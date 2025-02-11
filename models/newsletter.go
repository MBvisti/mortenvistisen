package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/google/uuid"
)

type Newsletter struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string
	Content    string
	ReleasedAt time.Time
	Released   bool
	Slug       string
}

type NewNewsletterPayload struct {
	Title      string
	Content    string
	ReleasedAt time.Time
	Released   bool
}

func NewNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	payload NewNewsletterPayload,
) (Newsletter, error) {
	now := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}
	releasedAt := pgtype.Timestamptz{
		Time:  payload.ReleasedAt,
		Valid: true,
	}
	id := uuid.New()

	dbID, err := db.Stmts.InsertNewsletter(ctx, dbtx, db.InsertNewsletterParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Title:     payload.Title,
		Slug: sql.NullString{
			String: slug.Make(payload.Title),
			Valid:  true,
		},
		Content:    payload.Content,
		ReleasedAt: releasedAt,
		Released: pgtype.Bool{
			Bool:  payload.Released,
			Valid: true,
		},
	})
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:         dbID,
		CreatedAt:  now.Time,
		UpdatedAt:  now.Time,
		Title:      payload.Title,
		Content:    payload.Content,
		ReleasedAt: releasedAt.Time,
		Released:   payload.Released,
	}, nil
}

func GetNewsletterByID(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Newsletter, error) {
	dbNewsletter, err := db.Stmts.QueryNewsletterByID(ctx, dbtx, id)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:         dbNewsletter.ID,
		CreatedAt:  dbNewsletter.CreatedAt.Time,
		UpdatedAt:  dbNewsletter.UpdatedAt.Time,
		Title:      dbNewsletter.Title,
		Content:    dbNewsletter.Content,
		ReleasedAt: dbNewsletter.ReleasedAt.Time,
		Released:   dbNewsletter.Released.Bool,
	}, nil
}

type QueryNewslettersParams struct {
	Limit  int32
	Offset int32
}

func GetNewslettersPage(
	ctx context.Context,
	dbtx db.DBTX,
	params QueryNewslettersParams,
) ([]Newsletter, error) {
	dbNewsletters, err := db.Stmts.QueryNewsletters(
		ctx,
		dbtx,
		db.QueryNewslettersParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(dbNewsletters))
	for i, dbNewsletter := range dbNewsletters {
		newsletters[i] = Newsletter{
			ID:         dbNewsletter.ID,
			CreatedAt:  dbNewsletter.CreatedAt.Time,
			UpdatedAt:  dbNewsletter.UpdatedAt.Time,
			Title:      dbNewsletter.Title,
			Content:    dbNewsletter.Content,
			ReleasedAt: dbNewsletter.ReleasedAt.Time,
			Released:   dbNewsletter.Released.Bool,
		}
	}

	return newsletters, nil
}

func GetNewslettersCount(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	return db.Stmts.QueryNewslettersCount(ctx, dbtx)
}

func UpdateNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	newsletter Newsletter,
) (Newsletter, error) {
	updatedAt := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}
	releasedAt := pgtype.Timestamptz{
		Time:  newsletter.ReleasedAt,
		Valid: true,
	}
	err := db.Stmts.UpdateNewsletter(ctx, dbtx, db.UpdateNewsletterParams{
		ID:         newsletter.ID,
		UpdatedAt:  updatedAt,
		Title:      newsletter.Title,
		Content:    newsletter.Content,
		ReleasedAt: releasedAt,
		Released: pgtype.Bool{
			Bool:  newsletter.Released,
			Valid: true,
		},
	})
	if err != nil {
		return Newsletter{}, err
	}

	return GetNewsletterByID(ctx, dbtx, newsletter.ID)
}

func GetAllNewsletters(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Newsletter, error) {
	newsletterRows, err := db.Stmts.QueryAllNewsletters(
		ctx,
		dbtx,
	)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(newsletterRows))
	for i, row := range newsletterRows {
		newsletters[i] = Newsletter{
			ID:         row.ID,
			CreatedAt:  row.CreatedAt.Time,
			UpdatedAt:  row.UpdatedAt.Time,
			Title:      row.Title,
			Content:    row.Content,
			ReleasedAt: row.ReleasedAt.Time,
			Released:   row.Released.Bool,
		}
	}

	return newsletters, nil
}

func DeleteNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.Stmts.DeleteNewsletter(ctx, dbtx, id)
}

func GetNewsletterBySlug(
	ctx context.Context,
	dbtx db.DBTX,
	slug string,
) (Newsletter, error) {
	newsletterRow, err := db.Stmts.QueryNewsletterBySlug(
		ctx,
		dbtx,
		sql.NullString{
			String: slug,
			Valid:  true,
		},
	)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:         newsletterRow.ID,
		CreatedAt:  newsletterRow.CreatedAt.Time,
		UpdatedAt:  newsletterRow.UpdatedAt.Time,
		Title:      newsletterRow.Title,
		Content:    newsletterRow.Content,
		ReleasedAt: newsletterRow.ReleasedAt.Time,
		Released:   newsletterRow.Released.Bool,
		Slug:       newsletterRow.Slug.String,
	}, nil
}
