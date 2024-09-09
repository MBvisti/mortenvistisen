package psql

import (
	"context"
	"database/sql"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql/database"
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

	return models.Newsletter{
		ID:         newsletter.ID,
		CreatedAt:  newsletter.CreatedAt.Time,
		UpdatedAt:  newsletter.UpdatedAt.Time,
		Title:      newsletter.Title,
		Content:    newsletter.Content,
		ReleasedAt: newsletter.ReleasedAt.Time,
		Released:   newsletter.Released.Bool,
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
		newsletters[i] = models.Newsletter{
			ID:         row.NewsletterID,
			CreatedAt:  row.NewsletterCreatedAt.Time,
			UpdatedAt:  row.NewsletterUpdatedAt.Time,
			Title:      row.NewsletterTitle,
			Content:    row.NewsletterContent,
			ReleasedAt: row.NewsletterReleasedAt.Time,
			Released:   row.NewsletterReleased.Bool,
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

	newNewsletter, err := p.Queries.InsertNewsletter(
		ctx,
		database.InsertNewsletterParams{
			ID:        data.ID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Title:     data.Title,
			Content:   data.Content,
		},
	)
	if err != nil {
		return models.Newsletter{}, err
	}

	return models.Newsletter{
		ID:        newNewsletter.ID,
		CreatedAt: newNewsletter.CreatedAt.Time,
		UpdatedAt: newNewsletter.UpdatedAt.Time,
		Content:   newNewsletter.Content,
		Title:     newNewsletter.Title,
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

	updatedNewsletter, err := p.Queries.UpdateNewsletter(
		ctx,
		database.UpdateNewsletterParams{
			UpdatedAt: updatedAt,
			Title:     newsletter.Title,
			Released: pgtype.Bool{
				Bool:  newsletter.Released,
				Valid: true,
			},
			ReleasedAt: releasedAt,
			ID:         newsletter.ID,
		},
	)
	if err != nil {
		return models.Newsletter{}, nil
	}

	return models.Newsletter{
		ID:         updatedNewsletter.ID,
		CreatedAt:  updatedNewsletter.CreatedAt.Time,
		UpdatedAt:  updatedNewsletter.UpdatedAt.Time,
		Title:      updatedNewsletter.Title,
		ReleasedAt: updatedNewsletter.ReleasedAt.Time,
		Released:   updatedNewsletter.Released.Bool,
	}, nil
}

func (p Postgres) Count(ctx context.Context) (int64, error) {
	return p.Queries.CountNewsletters(ctx)
}

func (p Postgres) CountReleased(ctx context.Context) (int64, error) {
	return p.Queries.CountReleasedNewsletters(ctx)
}
