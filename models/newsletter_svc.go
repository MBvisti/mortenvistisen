package models

import (
	"context"
	"errors"
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
		filters NewsletterFilters,
		opts ...PaginationOption,
	) ([]domain.Newsletter, error)
	InsertNewsletter(
		ctx context.Context,
		newsletter domain.Newsletter,
	) ([]domain.Newsletter, error)
}

type newsletterEmailService interface {
	SendNewSubscriberEmail(
		ctx context.Context,
		subscriberEmail string,
		activationToken, unsubscribeToken string,
	) error
}

type NewsletterFilters map[string]any

type NewsletterService struct {
	storage      newsletterStorage
	emailService newsletterEmailService
}

func NewNewsletterSvc(
	storage newsletterStorage,
	emailService newsletterEmailService,
) NewsletterService {
	return NewsletterService{
		storage,
		emailService,
	}
}

func (svc NewsletterService) ByID(ctx context.Context, id uuid.UUID) (domain.Newsletter, error) {
	newsletter, err := svc.storage.QueryNewsletterByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Newsletter{}, ErrNewsletterNotFound
		}

		return domain.Newsletter{}, errors.Join(ErrUnrecoverableEvent, err)
	}

	return newsletter, nil
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

		article, err := svc.storage.QueryArticleByID(ctx, id)
		if err != nil {
			return domain.Newsletter{}, err
		}

		articleSlug = article.Slug
	}

	filters := NewsletterFilters{
		"IsReleased": true,
	}

	releasedArticles, err := svc.storage.ListNewsletters(ctx, filters)
	if err != nil {
		return domain.Newsletter{}, err
	}

	return domain.InitilizeNewsletter(
		title,
		int32(len(releasedArticles))+1,
		newParagraphsElements,
		articleSlug,
	), nil
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

		article, err := svc.storage.QueryArticleByID(ctx, id)
		if err != nil {
			return domain.Newsletter{}, err
		}

		articleSlug = article.Slug
	}

	newsletter, err := domain.NewNewsletter(title, edition, paragraphs, articleSlug)
	if err != nil {
		return domain.Newsletter{}, err
	}

	if _, err := svc.storage.InsertNewsletter(ctx, newsletter); err != nil {
		return domain.Newsletter{}, err
	}

	return newsletter, nil
}

func (svc NewsletterService) Release(
	ctx context.Context,
	title string,
	edition string,
	paragraphs []string,
	articleID string,
) (domain.Newsletter, error) {
	// var associatedArticleSlug string
	// var associatedArticleID uuid.UUID
	//
	// if articleID != "" {
	// 	id, err := uuid.Parse(articleID)
	// 	if err != nil {
	// 		return domain.Newsletter{}, err
	// 	}
	// 	associatedArticle, err := n.db.QueryPostByID(ctx, id)
	// 	if err != nil {
	// 		return domain.Newsletter{}, err
	// 	}
	//
	// 	associatedArticleSlug = associatedArticle.Slug
	// 	associatedArticleID = id
	// }
	//
	// parsedEdition, err := strconv.Atoi(edition)
	// if err != nil {
	// 	return domain.Newsletter{}, err
	// }
	//
	// newsletter := domain.InitilizeNewsletter(
	// 	title,
	// 	int32(parsedEdition),
	// 	paragraphs,
	// 	associatedArticleSlug,
	// )
	//
	// now := time.Now()
	//
	// var buf bytes.Buffer
	// if err := json.NewEncoder(&buf).Encode(newsletter.Paragraphs); err != nil {
	// 	return domain.Newsletter{}, err
	// }
	//
	// _, err = n.db.InsertNewsletter(ctx, database.InsertNewsletterParams{
	// 	ID:                  newsletter.ID,
	// 	CreatedAt:           database.ConvertToPGTimestamptz(now),
	// 	UpdatedAt:           database.ConvertToPGTimestamptz(now),
	// 	Title:               newsletter.Title,
	// 	Edition:             sql.NullInt32{Int32: newsletter.Edition, Valid: true},
	// 	Body:                buf.Bytes(),
	// 	AssociatedArticleID: associatedArticleID,
	// })
	// if err != nil {
	return domain.Newsletter{}, nil
	// }
	//
	// return newsletter, nil
}

// func (n *Newsletter) ReleaseNewsletter(
// 	ctx context.Context,
// 	newsletter domain.Newsletter,
// ) (domain.ValidationErrsMap, error) {
// 	now := time.Now()
//
// 	updatedNewsletter, validationErrs, err := newsletter.CanBeReleased(n.v)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if len(validationErrs) > 0 {
// 		return validationErrs, nil
// 	}
//
// 	var buf bytes.Buffer
// 	if err := json.NewEncoder(&buf).Encode(updatedNewsletter.Paragraphs); err != nil {
// 		return nil, err
// 	}
//
// 	associatedArticle, err := n.db.QueryPostBySlug(ctx, updatedNewsletter.ArticleSlug)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	_, err = n.db.UpdateNewsletter(ctx, database.UpdateNewsletterParams{
// 		UpdatedAt:           database.ConvertToPGTimestamptz(now),
// 		Title:               updatedNewsletter.Title,
// 		Edition:             sql.NullInt32{Int32: updatedNewsletter.Edition, Valid: true},
// 		Released:            pgtype.Bool{Bool: updatedNewsletter.Released, Valid: true},
// 		ReleasedAt:          database.ConvertToPGTimestamptz(updatedNewsletter.ReleasedAt),
// 		Body:                buf.Bytes(),
// 		AssociatedArticleID: associatedArticle.ID,
// 		ID:                  updatedNewsletter.ID,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	verifiedSubs, err := n.db.QueryVerifiedSubscribers(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	newsletterMail := templates.NewsletterMail{
// 		Title:       newsletter.Title,
// 		Edition:     strconv.Itoa(int(newsletter.Edition)),
// 		Paragraphs:  newsletter.Paragraphs,
// 		ArticleLink: BuildURLFromSlug(FormatArticleSlug(newsletter.ArticleSlug)),
// 	}
//
// 	htmlMail, err := newsletterMail.GenerateHtmlVersion()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	textMail, err := newsletterMail.GenerateTextVersion()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, verifiedSub := range verifiedSubs {
// 		if err := n.mail.Send(
// 			ctx,
// 			verifiedSub.Email.String,
// 			"newsletter@mortenvistisen.com",
// 			fmt.Sprintf("MBV newsletter edition: %v", newsletter.Edition),
// 			textMail,
// 			htmlMail,
// 		); err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return nil, nil
// }

func (svc *NewsletterService) Update(
	ctx context.Context,
	title string,
	edition string,
	paragraphs []string,
	articleID string,
	id uuid.UUID,
) (domain.Newsletter, error) {
	newsletter, err := svc.storage.QueryNewsletterByID(ctx, id)
	if err != nil {
		return domain.Newsletter{}, err
	}

	parsedArticleID, err := uuid.Parse(articleID)
	if err != nil {
		return domain.Newsletter{}, err
	}

	article, err := svc.storage.QueryArticleByID(ctx, parsedArticleID)
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

	if _, err := svc.storage.UpdateNewsletter(ctx, updatedNewsletter); err != nil {
		return domain.Newsletter{}, err
	}

	return updatedNewsletter, nil
}
