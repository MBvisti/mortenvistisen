package psql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p Postgres) QueryNewsletterByID(
	ctx context.Context,
	id uuid.UUID,
) (domain.Newsletter, error) {
	newsletter, err := p.Queries.QueryNewsletterByID(ctx, id)
	if err != nil {
		return domain.Newsletter{}, err
	}

	var paragraphs []string
	if err := json.Unmarshal(newsletter.Body, &paragraphs); err != nil {
		return domain.Newsletter{}, err
	}

	article, err := p.Queries.QueryPostByID(ctx, newsletter.AssociatedArticleID)
	if err != nil {
		return domain.Newsletter{}, err
	}

	return domain.Newsletter{
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
) ([]domain.Newsletter, error) {
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

	newsletters := make([]domain.Newsletter, len(newsL))
	for i, row := range newsL {
		var paragraphs []string
		if err := json.Unmarshal(row.NewsletterBody, &paragraphs); err != nil {
			return nil, err
		}

		newsletters[i] = domain.Newsletter{
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
	data domain.Newsletter,
) (domain.Newsletter, error) {
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
		return domain.Newsletter{}, err
	}

	newNewsletter, err := p.Queries.InsertNewsletter(ctx, database.InsertNewsletterParams{
		ID:                  data.ID,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		Title:               data.Title,
		Edition:             sql.NullInt32{Int32: data.Edition, Valid: true},
		Body:                []byte{},
		AssociatedArticleID: article.ID,
	})
	if err != nil {
		return domain.Newsletter{}, err
	}

	var paragraphs []string
	if err := json.Unmarshal(newNewsletter.Body, &paragraphs); err != nil {
		return domain.Newsletter{}, err
	}

	return domain.Newsletter{
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
	newsletter domain.Newsletter,
) (domain.Newsletter, error) {
	return domain.Newsletter{}, nil
}

func (p Postgres) Count(ctx context.Context) (int64, error) {
	return p.Queries.CountNewsletters(ctx, pgtype.Bool{
		Bool:  false,
		Valid: false,
	})
}

func (p Postgres) CountReleased(ctx context.Context) (int64, error) {
	return p.Queries.CountNewsletters(ctx, pgtype.Bool{
		Bool:  true,
		Valid: true,
	})
}
