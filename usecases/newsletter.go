package usecases

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Newsletter struct {
	db   database.Queries
	v    *validator.Validate
	mail mail.Mail
}

func NewNewsletter(
	db database.Queries,
	v *validator.Validate,
	mail mail.Mail,
) Newsletter {
	return Newsletter{
		db:   db,
		v:    v,
		mail: mail,
	}
}

func (n *Newsletter) Get(ctx context.Context, id uuid.UUID) (domain.Newsletter, error) {
	newsletterModel, err := n.db.QueryNewsletterByID(ctx, id)
	if err != nil {
		return domain.Newsletter{}, err
	}

	associatedArticle, err := n.db.QueryPostByID(
		ctx,
		newsletterModel.AssociatedArticleID,
	)
	if err != nil {
		return domain.Newsletter{}, err
	}

	return newsletterModel.ConvertNewsletterToDomain(associatedArticle.Slug)
}

func (n *Newsletter) Preview(
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

		article, err := n.db.QueryPostByID(ctx, id)
		if err != nil {
			return domain.Newsletter{}, err
		}

		articleSlug = article.Slug
	}

	releasedNewslettersCount, err := n.db.QueryReleasedNewslettersCount(
		ctx,
	)
	if err != nil {
		return domain.Newsletter{}, err
	}

	return domain.InitilizeNewsletter(
		title,
		int32(releasedNewslettersCount)+1,
		newParagraphsElements,
		articleSlug,
	), nil
}

func (n *Newsletter) Create(
	ctx context.Context,
	title string,
	edition string,
	paragraphs []string,
	articleID string,
) (domain.Newsletter, error) {
	var associatedArticleSlug string
	var associatedArticleID uuid.UUID

	if articleID != "" {
		id, err := uuid.Parse(articleID)
		if err != nil {
			return domain.Newsletter{}, err
		}
		associatedArticle, err := n.db.QueryPostByID(ctx, id)
		if err != nil {
			return domain.Newsletter{}, err
		}

		associatedArticleSlug = associatedArticle.Slug
		associatedArticleID = id
	}

	parsedEdition, err := strconv.Atoi(edition)
	if err != nil {
		return domain.Newsletter{}, err
	}

	newsletter := domain.InitilizeNewsletter(
		title,
		int32(parsedEdition),
		paragraphs,
		associatedArticleSlug,
	)

	now := time.Now()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(newsletter.Paragraphs); err != nil {
		return domain.Newsletter{}, err
	}

	_, err = n.db.InsertNewsletter(ctx, database.InsertNewsletterParams{
		ID:                  newsletter.ID,
		CreatedAt:           database.ConvertToPGTimestamptz(now),
		UpdatedAt:           database.ConvertToPGTimestamptz(now),
		Title:               newsletter.Title,
		Edition:             sql.NullInt32{Int32: newsletter.Edition, Valid: true},
		Body:                buf.Bytes(),
		AssociatedArticleID: associatedArticleID,
	})
	if err != nil {
		return domain.Newsletter{}, err
	}

	return newsletter, nil
}

func (n *Newsletter) ReleaseNewsletter(
	ctx context.Context,
	newsletter domain.Newsletter,
) (domain.ValidationErrsMap, error) {
	now := time.Now()

	updatedNewsletter, validationErrs, err := newsletter.CanBeReleased(n.v)
	if err != nil {
		return nil, err
	}

	if len(validationErrs) > 0 {
		return validationErrs, nil
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updatedNewsletter.Paragraphs); err != nil {
		return nil, err
	}

	associatedArticle, err := n.db.QueryPostBySlug(ctx, updatedNewsletter.ArticleSlug)
	if err != nil {
		return nil, err
	}

	_, err = n.db.UpdateNewsletter(ctx, database.UpdateNewsletterParams{
		UpdatedAt:           database.ConvertToPGTimestamptz(now),
		Title:               updatedNewsletter.Title,
		Edition:             sql.NullInt32{Int32: updatedNewsletter.Edition, Valid: true},
		Released:            pgtype.Bool{Bool: updatedNewsletter.Released, Valid: true},
		ReleasedAt:          database.ConvertToPGTimestamptz(updatedNewsletter.ReleasedAt),
		Body:                buf.Bytes(),
		AssociatedArticleID: associatedArticle.ID,
		ID:                  updatedNewsletter.ID,
	})
	if err != nil {
		return nil, err
	}

	verifiedSubs, err := n.db.QueryVerifiedSubscribers(ctx)
	if err != nil {
		return nil, err
	}

	newsletterMail := templates.NewsletterMail{
		Title:       newsletter.Title,
		Edition:     strconv.Itoa(int(newsletter.Edition)),
		Paragraphs:  newsletter.Paragraphs,
		ArticleLink: BuildURLFromSlug(FormatArticleSlug(newsletter.ArticleSlug)),
	}

	htmlMail, err := newsletterMail.GenerateHtmlVersion()
	if err != nil {
		return nil, err
	}

	textMail, err := newsletterMail.GenerateTextVersion()
	if err != nil {
		return nil, err
	}

	for _, verifiedSub := range verifiedSubs {
		if err := n.mail.Send(
			ctx,
			verifiedSub.Email.String,
			"newsletter@mortenvistisen.com",
			fmt.Sprintf("MBV newsletter edition: %v", newsletter.Edition),
			textMail,
			htmlMail,
		); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (n *Newsletter) Update(
	ctx context.Context,
	title string,
	edition string,
	paragraphs []string,
	articleID string,
	id uuid.UUID,
) (domain.Newsletter, domain.ValidationErrsMap, error) {
	newsletterModel, err := n.db.QueryNewsletterByID(ctx, id)
	if err != nil {
		return domain.Newsletter{}, nil, err
	}

	associatedArticle, err := n.db.QueryPostByID(ctx, newsletterModel.AssociatedArticleID)
	if err != nil {
		return domain.Newsletter{}, nil, err
	}

	newsletter, err := newsletterModel.ConvertNewsletterToDomain(associatedArticle.Slug)
	if err != nil {
		return domain.Newsletter{}, nil, err
	}

	parsedEdition, err := strconv.Atoi(edition)
	if err != nil {
		return domain.Newsletter{}, nil, err
	}

	updatedNewsletter, err := newsletter.Update(domain.UpdateNewsletterPayload{
		ID:          id,
		Title:       title,
		Edition:     int32(parsedEdition),
		Paragraphs:  paragraphs,
		ArticleSlug: associatedArticle.Slug,
	}, n.v)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return domain.Newsletter{}, nil, err
		}

		errors := make(domain.ValidationErrsMap, len(validationErrors))
		for _, valiErr := range validationErrors {
			switch valiErr.Field() {
			case "ID":
				errors[valiErr.Field()] = "a valid uuid v4 must be provided"
			case "Title":
				errors[valiErr.Field()] = "title cannot be empty"
			case "Paragraphs":
				errors[valiErr.Field()] = "atleast one paragraph is needed"
			case "ArticleSlug":
				errors[valiErr.Field()] = "an article slug must be provided"
			case "Edition":
				errors[valiErr.Field()] = "edition is required and must be > 0"
			}
		}
		return domain.Newsletter{}, errors, nil
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updatedNewsletter.Paragraphs); err != nil {
		return domain.Newsletter{}, nil, err
	}

	_, err = n.db.UpdateNewsletter(ctx, database.UpdateNewsletterParams{
		UpdatedAt:           database.ConvertToPGTimestamptz(time.Now()),
		Title:               updatedNewsletter.Title,
		Edition:             sql.NullInt32{Int32: updatedNewsletter.Edition, Valid: true},
		Body:                buf.Bytes(),
		AssociatedArticleID: newsletterModel.AssociatedArticleID,
		ID:                  updatedNewsletter.ID,
	})
	if err != nil {
		return domain.Newsletter{}, nil, err
	}

	return updatedNewsletter, nil, nil
}
