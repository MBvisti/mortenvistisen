package models

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type newsletterStorage interface {
	QueryNewsletterByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Newsletter, error)
	UpdateNewsletter(
		ctx context.Context,
		newsletter domain.Newsletter,
	) (domain.Newsletter, error)
	QueryArticleBySlug(
		ctx context.Context,
		slug string,
	) (domain.Article, error)
	QueryArticleByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Article, error)
	ListNewsletters(
		ctx context.Context,
		filters QueryFilters,
		opts ...PaginationOption,
	) ([]domain.Newsletter, error)
	InsertNewsletter(
		ctx context.Context,
		newsletter domain.Newsletter,
	) (domain.Newsletter, error)
	Count(ctx context.Context) (int64, error)
	CountReleased(ctx context.Context) (int64, error)
}

type newsletterEmailService interface {
	SendNewSubscriberEmail(
		ctx context.Context,
		subscriberEmail string,
		activationToken, unsubscribeToken string,
	) error
}

type QueryFilters map[string]any

type NewsletterService struct {
	newsletterStorage newsletterStorage
	subscriberStorage subscriberStorage
	tknService        subscriberTokenService
	emailService      newsletterEmailService
}

func NewNewsletterSvc(
	newsletterStorage newsletterStorage,
	subscriberStorage subscriberStorage,
	tknService subscriberTokenService,
	emailService newsletterEmailService,
) NewsletterService {
	return NewsletterService{
		newsletterStorage,
		subscriberStorage,
		tknService,
		emailService,
	}
}

func (svc NewsletterService) ByID(ctx context.Context, id uuid.UUID) (domain.Newsletter, error) {
	newsletter, err := svc.newsletterStorage.QueryNewsletterByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Newsletter{}, ErrNewsletterNotFound
		}

		return domain.Newsletter{}, errors.Join(ErrUnrecoverableEvent, err)
	}

	return newsletter, nil
}

func (svc NewsletterService) List(
	ctx context.Context,
	limit, offset int32,
) ([]domain.Newsletter, error) {
	return svc.newsletterStorage.ListNewsletters(ctx, nil, WithPagination(limit, offset))
}

func (svc NewsletterService) Count(
	ctx context.Context,
	releasedOnly bool,
) (int64, error) {
	if releasedOnly {
		return svc.newsletterStorage.CountReleased(ctx)
	}

	return svc.newsletterStorage.Count(ctx)
}

func (svc NewsletterService) Preview(
	ctx context.Context,
	paragraphIndex string,
	action string,
	title string,
	paragraphElements []string,
	newParagraphElement string,
	articleID string,
) (domain.Newsletter, error) {
	newParagraphsElements := paragraphElements
	if paragraphIndex != "" && action == "del" {
		index, err := strconv.Atoi(paragraphIndex)
		if err != nil {
			return domain.Newsletter{}, err
		}

		if action == "del" {
			newParagraphsElements = append(
				newParagraphsElements[:index],
				newParagraphsElements[index+1:]...)
		}
	}

	if newParagraphElement != "" && action != "del" {
		newParagraphsElements = append(
			newParagraphsElements,
			newParagraphElement,
		)
	}

	var articleSlug string
	if articleID != "" {
		id, err := uuid.Parse(articleID)
		if err != nil {
			return domain.Newsletter{}, err
		}

		article, err := svc.newsletterStorage.QueryArticleByID(ctx, id)
		if err != nil {
			return domain.Newsletter{}, err
		}

		articleSlug = article.Slug
	}

	filters := QueryFilters{
		"IsReleased": true,
	}

	releasedArticles, err := svc.newsletterStorage.ListNewsletters(ctx, filters)
	if err != nil {
		return domain.Newsletter{}, err
	}

	return domain.CreateNewsletter(
		title,
		int32(len(releasedArticles))+1,
		newParagraphsElements,
		articleSlug,
	)
}

func (svc NewsletterService) CreateDraft(
	ctx context.Context,
	title string,
	edition int32,
	paragraphs []string,
	articleID string,
) (domain.Newsletter, error) {
	var articleSlug string
	if articleID != "" {
		id, err := uuid.Parse(articleID)
		if err != nil {
			return domain.Newsletter{}, err
		}

		article, err := svc.newsletterStorage.QueryArticleByID(ctx, id)
		if err != nil {
			return domain.Newsletter{}, err
		}

		articleSlug = article.Slug
	}

	newsletter, err := domain.CreateNewsletter(title, edition, paragraphs, articleSlug)
	if err != nil {
		return domain.Newsletter{}, err
	}

	if _, err := svc.newsletterStorage.InsertNewsletter(ctx, newsletter); err != nil {
		return domain.Newsletter{}, err
	}

	return newsletter, nil
}

func (svc NewsletterService) Release(
	ctx context.Context,
	instance domain.Newsletter,
) (domain.Newsletter, error) {
	newsletter, err := instance.Release()
	if err != nil {
		return domain.Newsletter{}, err
	}

	updatedNewsletter, err := svc.newsletterStorage.UpdateNewsletter(ctx, newsletter)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Newsletter{}, errors.Join(ErrNewsletterNotFound, err)
		}

		slog.Error("could not update newsletter", "error", err)
		return domain.Newsletter{}, err
	}

	verifiedSubscribers, err := svc.subscriberStorage.ListSubscribers(
		ctx,
		QueryFilters{"IsVerified": true},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("could not find any verified subscribers")
			return domain.Newsletter{}, nil
		}

		return domain.Newsletter{}, errors.Join(ErrUnrecoverableEvent, err)
	}

	for _, verifiedSubscriber := range verifiedSubscribers {
		subTkn, err := svc.tknService.CreateSubscriptionToken(ctx, verifiedSubscriber.ID)
		if err != nil {
			return domain.Newsletter{}, errors.Join(ErrUnrecoverableEvent, err)
		}

		unsubTkn, err := svc.tknService.CreateUnsubscribeToken(ctx, verifiedSubscriber.ID)
		if err != nil {
			return domain.Newsletter{}, errors.Join(ErrUnrecoverableEvent, err)
		}

		if err := svc.emailService.SendNewSubscriberEmail(ctx, verifiedSubscriber.Email, subTkn, unsubTkn); err != nil {
			return domain.Newsletter{}, errors.Join(ErrUnrecoverableEvent, err)
		}
	}

	return updatedNewsletter, nil
}

func (svc *NewsletterService) Update(
	ctx context.Context,
	title string,
	edition string,
	paragraphs []string,
	articleID string,
	id uuid.UUID,
) (domain.Newsletter, error) {
	newsletter, err := svc.newsletterStorage.QueryNewsletterByID(ctx, id)
	if err != nil {
		return domain.Newsletter{}, err
	}

	parsedArticleID, err := uuid.Parse(articleID)
	if err != nil {
		return domain.Newsletter{}, err
	}

	article, err := svc.newsletterStorage.QueryArticleByID(ctx, parsedArticleID)
	if err != nil {
		return domain.Newsletter{}, err
	}

	parsedEdition, err := strconv.Atoi(edition)
	if err != nil {
		return domain.Newsletter{}, err
	}

	updatedNewsletter, err := newsletter.Update(
		title,
		int32(parsedEdition),
		paragraphs,
		article.Slug,
	)
	if err != nil {
		return domain.Newsletter{}, err
	}

	if _, err := svc.newsletterStorage.UpdateNewsletter(ctx, updatedNewsletter); err != nil {
		return domain.Newsletter{}, err
	}

	return updatedNewsletter, nil
}
