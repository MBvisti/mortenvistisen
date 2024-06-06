package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Newsletter struct {
	ID         uuid.UUID `validate:"required"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string    `validate:"required,gte=3"`
	Edition    int32     `validate:"required,gte=1"`
	Released   bool      `validate:"required,eq=true"`
	ReleasedAt time.Time `validate:"required"`
	Body       []string  `validate:"required,gte=1"`
	ArticleID  uuid.UUID `validate:"required"`
}

type NewsletterModel struct {
	db *database.Queries
	v  *validator.Validate
}

func NewNewsletter(db *database.Queries, v *validator.Validate) NewsletterModel {
	return NewsletterModel{
		db,
		v,
	}
}

func (n *NewsletterModel) ByID(ctx context.Context, id uuid.UUID) (Newsletter, error) {
	newsletter, err := n.db.QueryNewsletterByID(ctx, id)
	if err != nil {
		return Newsletter{}, err
	}

	var newsletterBody []string
	if err := json.Unmarshal(newsletter.Body, &newsletterBody); err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:         newsletter.ID,
		CreatedAt:  newsletter.CreatedAt.Time,
		UpdatedAt:  newsletter.UpdatedAt.Time,
		Title:      newsletter.Title,
		Edition:    newsletter.Edition.Int32,
		Released:   newsletter.Released.Bool,
		ReleasedAt: newsletter.ReleasedAt.Time,
		Body:       newsletterBody,
		ArticleID:  newsletter.AssociatedArticleID,
	}, nil
}

func (n *NewsletterModel) CreateDraft(
	ctx context.Context,
	title string,
	edition int32,
	body []string,
	articeID uuid.UUID,
) (Newsletter, error) {
	newsletter := Newsletter{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     title,
		Edition:   edition,
		Body:      body,
		ArticleID: articeID,
	}

	if err := n.v.StructPartial(newsletter, "ArticleID"); err != nil {
		if errors.As(err, validator.ValidationErrors) {
		}
		return newsletter, nil
	}

	return Newsletter{
		ID:         newsletter.ID,
		CreatedAt:  newsletter.CreatedAt.Time,
		UpdatedAt:  newsletter.UpdatedAt.Time,
		Title:      newsletter.Title,
		Edition:    newsletter.Edition.Int32,
		Released:   newsletter.Released.Bool,
		ReleasedAt: newsletter.ReleasedAt.Time,
		Body:       newsletterBody,
		ArticleID:  newsletter.AssociatedArticleID,
	}, nil, nil
}

func (n *NewsletterModel) ReadyForRelease(
	ctx context.Context,
	id uuid.UUID,
) (ValidationErrsMap, error) {
	newsletter, err := n.ByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := n.v.Struct(newsletter); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil, errors.New("could not convert errors to ValidationErrors")
		}

		errors := make(ValidationErrsMap, len(validationErrors))
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
			case "Released":
				errors[valiErr.Field()] = "released must be set to true"
			case "ReleasedAt":
				errors[valiErr.Field()] = "released date must be set"
			}
		}

		return errors, nil
	}

	return nil, nil
}

func (n *NewsletterModel) Update(
	ctx context.Context,
	id uuid.UUID,
	title string,
	edition int32,
	released bool,
	releasedAt time.Time,
	body []string,
	articleID uuid.UUID,
) (Newsletter, error) {
	marshaledBody, err := json.Marshal(body)
	if err != nil {
		return Newsletter{}, err
	}

	updatedNewsletter, err := n.db.UpdateNewsletter(ctx, database.UpdateNewsletterParams{
		UpdatedAt:           ConvertToPGTimestamptz(time.Now()),
		Title:               title,
		Edition:             sql.NullInt32{Int32: edition, Valid: true},
		Released:            pgtype.Bool{Bool: released, Valid: true},
		ReleasedAt:          ConvertToPGTimestamptz(releasedAt),
		Body:                marshaledBody,
		AssociatedArticleID: articleID,
		ID:                  id,
	})
	if err != nil {
		return Newsletter{}, err
	}

	var updatedNewsletterBody []string
	if err := json.Unmarshal(updatedNewsletter.Body, &updatedNewsletterBody); err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:         updatedNewsletter.ID,
		CreatedAt:  updatedNewsletter.CreatedAt.Time,
		UpdatedAt:  updatedNewsletter.UpdatedAt.Time,
		Title:      updatedNewsletter.Title,
		Edition:    updatedNewsletter.Edition.Int32,
		Released:   updatedNewsletter.Released.Bool,
		ReleasedAt: updatedNewsletter.ReleasedAt.Time,
		Body:       updatedNewsletterBody,
		ArticleID:  updatedNewsletter.AssociatedArticleID,
	}, nil
}

func (n *NewsletterModel) List(
	ctx context.Context,
	opts ...listOpt,
) ([]Newsletter, error) {
	options := &listOptions{}

	for _, opt := range opts {
		opt(options)
	}

	newslettersModels, err := n.db.QueryNewsletters(ctx, database.QueryNewslettersParams{
		Offset: options.offset,
		Limit:  options.limit,
	})
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(newslettersModels))
	for i, newsletter := range newslettersModels {
		var newsletterBody []string
		if err := json.Unmarshal(newsletter.Body, &newsletterBody); err != nil {
			return nil, err
		}

		newsletters[i] = Newsletter{
			ID:         newsletter.ID,
			CreatedAt:  newsletter.CreatedAt.Time,
			UpdatedAt:  newsletter.UpdatedAt.Time,
			Title:      newsletter.Title,
			Edition:    newsletter.Edition.Int32,
			Released:   newsletter.Released.Bool,
			ReleasedAt: newsletter.ReleasedAt.Time,
			Body:       newsletterBody,
			ArticleID:  newsletter.AssociatedArticleID,
		}
	}

	return newsletters, nil
}

func (n *NewsletterModel) GetCount(ctx context.Context) (int64, error) {
	count, err := n.db.QueryNewslettersCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (n *NewsletterModel) GetReleasedCount(ctx context.Context) (int64, error) {
	count, err := n.db.QueryReleasedNewslettersCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
