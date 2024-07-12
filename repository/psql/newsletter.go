package psql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p Postgres) QueryNewsletterByID(
	ctx context.Context,
	id uuid.UUID,
) (models.Newsletter, error) {
	newsletter, err := p.Queries.QueryNewsletterByID(ctx, id)
	if err != nil {
		return models.Newsletter{}, err
	}

	var paragraphs []string
	if err := json.Unmarshal(newsletter.Body, &paragraphs); err != nil {
		return models.Newsletter{}, err
	}

	article, err := p.Queries.QueryPostByID(ctx, newsletter.AssociatedArticleID)
	if err != nil {
		return models.Newsletter{}, err
	}

	return models.Newsletter{
		ID:          newsletter.ID,
		CreatedAt:   newsletter.CreatedAt.Time,
		UpdatedAt:   newsletter.UpdatedAt.Time,
		Title:       newsletter.Title,
		Edition:     newsletter.Edition.Int32,
		ReleasedAt:  newsletter.ReleasedAt.Time,
		Released:    newsletter.Released.Bool,
		Paragraphs:  paragraphs,
		ArticleSlug: article.Slug,
	}, nil
}

func (p Postgres) ListNewsletters(
	ctx context.Context,
	filters models.QueryFilters,
	opts ...models.PaginationOption,
) ([]models.Newsletter, error) {
	options := &models.PaginationOptions{}

	for _, opt := range opts {
		opt(options)
	}

	params := database.QueryNewslettersParams{
		Offset: sql.NullInt32{Int32: options.Offset, Valid: true},
		Limit:  sql.NullInt32{Int32: options.Limit, Valid: true},
	}

	for k, v := range filters {
		if k == "IsReleased" {
			val, ok := v.(bool)
			if ok {
				params.IsReleased = pgtype.Bool{Bool: val, Valid: true}
			}
		}
	}

	newsL, err := p.Queries.QueryNewsletters(ctx, params)
	if err != nil {
		return nil, err
	}

	newsletters := make([]models.Newsletter, len(newsL))
	for i, row := range newsL {
		var paragraphs []string
		if err := json.Unmarshal(row.NewsletterBody, &paragraphs); err != nil {
			return nil, err
		}

		newsletters[i] = models.Newsletter{
			ID:          row.NewsletterID,
			CreatedAt:   row.NewsletterCreatedAt.Time,
			UpdatedAt:   row.NewsletterUpdatedAt.Time,
			Title:       row.NewsletterTitle,
			Edition:     row.NewsletterEdition.Int32,
			ReleasedAt:  row.NewsletterReleasedAt.Time,
			Released:    row.NewsletterReleased.Bool,
			Paragraphs:  paragraphs,
			ArticleSlug: row.PostSlug,
		}
	}

	return newsletters, nil
}

func (p Postgres) InsertNewsletter(
	ctx context.Context,
	data models.Newsletter,
) (models.Newsletter, error) {
	createdAt := pgtype.Timestamptz{
		Time:  data.CreatedAt,
		Valid: true,
	}
	updatedAt := pgtype.Timestamptz{
		Time:  data.UpdatedAt,
		Valid: true,
	}

	article, err := p.Queries.QueryPostBySlug(ctx, data.ArticleSlug)
	if err != nil {
		return models.Newsletter{}, err
	}

	body, err := json.Marshal(data.Paragraphs)
	if err != nil {
		return models.Newsletter{}, err
	}

	newNewsletter, err := p.Queries.InsertNewsletter(ctx, database.InsertNewsletterParams{
		ID:                  data.ID,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		Title:               data.Title,
		Edition:             sql.NullInt32{Int32: data.Edition, Valid: true},
		Body:                body,
		AssociatedArticleID: article.ID,
	})
	if err != nil {
		return models.Newsletter{}, err
	}

	var paragraphs []string
	if err := json.Unmarshal(newNewsletter.Body, &paragraphs); err != nil {
		return models.Newsletter{}, err
	}

	return models.Newsletter{
		ID:          newNewsletter.ID,
		CreatedAt:   newNewsletter.CreatedAt.Time,
		UpdatedAt:   newNewsletter.UpdatedAt.Time,
		Title:       newNewsletter.Title,
		Edition:     newNewsletter.Edition.Int32,
		Paragraphs:  paragraphs,
		ArticleSlug: article.Slug,
	}, nil
}

// TODO
func (p Postgres) UpdateNewsletter(
	ctx context.Context,
	newsletter models.Newsletter,
) (models.Newsletter, error) {
	updatedAt := pgtype.Timestamptz{
		Time:  newsletter.UpdatedAt,
		Valid: true,
	}
	releasedAt := pgtype.Timestamptz{
		Time:  newsletter.ReleasedAt,
		Valid: true,
	}

	body, err := json.Marshal(newsletter.Paragraphs)
	if err != nil {
		return models.Newsletter{}, err
	}

	article, err := p.QueryArticleBySlug(ctx, newsletter.ArticleSlug)
	if err != nil {
		return models.Newsletter{}, err
	}

	updatedNewsletter, err := p.Queries.UpdateNewsletter(ctx, database.UpdateNewsletterParams{
		UpdatedAt:           updatedAt,
		Title:               newsletter.Title,
		Edition:             sql.NullInt32{Int32: newsletter.Edition, Valid: true},
		Released:            pgtype.Bool{Bool: newsletter.Released, Valid: true},
		ReleasedAt:          releasedAt,
		Body:                body,
		AssociatedArticleID: article.ID,
		ID:                  newsletter.ID,
	})
	if err != nil {
		return models.Newsletter{}, nil
	}

	var paragraphs []string
	if err := json.Unmarshal(updatedNewsletter.Body, &paragraphs); err != nil {
		return models.Newsletter{}, nil
	}

	return models.Newsletter{
		ID:          updatedNewsletter.ID,
		CreatedAt:   updatedNewsletter.CreatedAt.Time,
		UpdatedAt:   updatedNewsletter.UpdatedAt.Time,
		Title:       updatedNewsletter.Title,
		Edition:     updatedNewsletter.Edition.Int32,
		ReleasedAt:  updatedNewsletter.ReleasedAt.Time,
		Released:    updatedNewsletter.Released.Bool,
		Paragraphs:  paragraphs,
		ArticleSlug: article.Slug,
	}, nil
}

func (p Postgres) Count(ctx context.Context) (int64, error) {
	return p.Queries.CountNewsletters(ctx)
}

func (p Postgres) CountReleased(ctx context.Context) (int64, error) {
	return p.Queries.CountReleasedNewsletters(ctx)
}
