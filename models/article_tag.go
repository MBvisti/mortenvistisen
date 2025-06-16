package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type ArticleTag struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
}

func GetArticleTagByID(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (ArticleTag, error) {
	row, err := db.Stmts.QueryArticleTagByID(ctx, dbtx, id)
	if err != nil {
		return ArticleTag{}, err
	}

	return ArticleTag{
		ID:        row.ID,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
		Title:     row.Title,
	}, nil
}

func GetArticleTagByTitle(
	ctx context.Context,
	dbtx db.DBTX,
	title string,
) (ArticleTag, error) {
	row, err := db.Stmts.QueryArticleTagByTitle(ctx, dbtx, title)
	if err != nil {
		return ArticleTag{}, err
	}

	return ArticleTag{
		ID:        row.ID,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
		Title:     row.Title,
	}, nil
}

func GetArticleTags(
	ctx context.Context,
	dbtx db.DBTX,
) ([]ArticleTag, error) {
	rows, err := db.Stmts.QueryArticleTags(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	tags := make([]ArticleTag, len(rows))
	for i, row := range rows {
		tags[i] = ArticleTag{
			ID:        row.ID,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			Title:     row.Title,
		}
	}

	return tags, nil
}

func GetArticleTagsByArticleID(
	ctx context.Context,
	dbtx db.DBTX,
	articleID uuid.UUID,
) ([]ArticleTag, error) {
	rows, err := db.Stmts.QueryArticleTagsByArticleID(ctx, dbtx, articleID)
	if err != nil {
		return nil, err
	}

	tags := make([]ArticleTag, len(rows))
	for i, row := range rows {
		tags[i] = ArticleTag{
			ID:        row.ID,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			Title:     row.Title,
		}
	}

	return tags, nil
}

type NewArticleTagPayload struct {
	Title string `validate:"required,max=255"`
}

func NewArticleTag(
	ctx context.Context,
	dbtx db.DBTX,
	data NewArticleTagPayload,
) (ArticleTag, error) {
	if err := validate.Struct(data); err != nil {
		return ArticleTag{}, errors.Join(ErrDomainValidation, err)
	}

	tag := ArticleTag{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     strings.ToLower(slug.Make(data.Title)),
	}

	_, err := db.Stmts.InsertArticleTag(ctx, dbtx, db.InsertArticleTagParams{
		ID: tag.ID,
		CreatedAt: pgtype.Timestamptz{
			Time:  tag.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  tag.UpdatedAt,
			Valid: true,
		},
		Title: tag.Title,
	})
	if err != nil {
		return ArticleTag{}, err
	}

	return tag, nil
}

type UpdateArticleTagPayload struct {
	ID        uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
	Title     string    `validate:"required,max=255"`
}

func UpdateArticleTag(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateArticleTagPayload,
) (ArticleTag, error) {
	if err := validate.Struct(data); err != nil {
		return ArticleTag{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateArticleTag(ctx, dbtx, db.UpdateArticleTagParams{
		ID:        data.ID,
		UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		Title:     data.Title,
	})
	if err != nil {
		return ArticleTag{}, err
	}

	return ArticleTag{
		ID:        row.ID,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
		Title:     row.Title,
	}, nil
}

func DeleteArticleTag(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.Stmts.DeleteArticleTag(ctx, dbtx, id)
}
