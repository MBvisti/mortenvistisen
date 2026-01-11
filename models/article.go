package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models/internal/db"
)

type Article struct {
	ID               uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstPublishedAt time.Time
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	Slug             string
	ImageLink        string
	ReadTime         int32
	Content          string
}

func FindArticle(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) (Article, error) {
	row, err := queries.QueryArticleByID(ctx, exec, id)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

type CreateArticleData struct {
	FirstPublishedAt time.Time
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	Slug             string
	ImageLink        string
	ReadTime         int32
	Content          string
}

func CreateArticle(
	ctx context.Context,
	exec storage.Executor,
	data CreateArticleData,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.InsertArticleParams{
		ID:               uuid.New(),
		FirstPublishedAt: pgtype.Timestamptz{Time: data.FirstPublishedAt, Valid: true},
		Title:            data.Title,
		Excerpt:          data.Excerpt,
		MetaTitle:        data.MetaTitle,
		MetaDescription:  data.MetaDescription,
		Slug:             data.Slug,
		ImageLink:        pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:         pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:          pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.InsertArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

type UpdateArticleData struct {
	ID               uuid.UUID
	UpdatedAt        time.Time
	FirstPublishedAt time.Time
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	Slug             string
	ImageLink        string
	ReadTime         int32
	Content          string
}

func UpdateArticle(
	ctx context.Context,
	exec storage.Executor,
	data UpdateArticleData,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpdateArticleParams{
		ID:               data.ID,
		FirstPublishedAt: pgtype.Timestamptz{Time: data.FirstPublishedAt, Valid: true},
		Title:            data.Title,
		Excerpt:          data.Excerpt,
		MetaTitle:        data.MetaTitle,
		MetaDescription:  data.MetaDescription,
		Slug:             data.Slug,
		ImageLink:        pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:         pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:          pgtype.Text{String: data.Content, Valid: true},
	}

	row, err := queries.UpdateArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

func DestroyArticle(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) error {
	return queries.DeleteArticle(ctx, exec, id)
}

func AllArticles(
	ctx context.Context,
	exec storage.Executor,
) ([]Article, error) {
	rows, err := queries.QueryArticles(ctx, exec)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = rowToArticle(row)
	}

	return articles, nil
}

type PaginatedArticles struct {
	Articles   []Article
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func PaginateArticles(
	ctx context.Context,
	exec storage.Executor,
	page int64,
	pageSize int64,
) (PaginatedArticles, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := queries.CountArticles(ctx, exec)
	if err != nil {
		return PaginatedArticles{}, err
	}

	rows, err := queries.QueryPaginatedArticles(
		ctx,
		exec,
		db.QueryPaginatedArticlesParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedArticles{}, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = rowToArticle(row)
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedArticles{
		Articles:   articles,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func UpsertArticle(
	ctx context.Context,
	exec storage.Executor,
	data CreateArticleData,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpsertArticleParams{
		ID:               uuid.New(),
		FirstPublishedAt: pgtype.Timestamptz{Time: data.FirstPublishedAt, Valid: true},
		Title:            data.Title,
		Excerpt:          data.Excerpt,
		MetaTitle:        data.MetaTitle,
		MetaDescription:  data.MetaDescription,
		Slug:             data.Slug,
		ImageLink:        pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:         pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:          pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.UpsertArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

func rowToArticle(row db.Article) Article {
	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		Title:            row.Title,
		Excerpt:          row.Excerpt,
		MetaTitle:        row.MetaTitle,
		MetaDescription:  row.MetaDescription,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		ReadTime:         row.ReadTime.Int32,
		Content:          row.Content.String,
	}
}
