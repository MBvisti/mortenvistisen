package models

import (
	"context"
	"errors"
	"time"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models/internal/db"
)

type Article struct {
	ID               int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstPublishedAt time.Time
	Published        bool
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
	id int32,
) (Article, error) {
	row, err := queries.QueryArticleByID(ctx, exec, id)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

func FindArticleBySlug(
	ctx context.Context,
	exec storage.Executor,
	slug string,
) (Article, error) {
	row, err := queries.QueryArticleBySlug(ctx, exec, slug)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

type CreateArticleData struct {
	FirstPublishedAt time.Time
	Published        bool
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
	if err := Validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.InsertArticleParams{
		FirstPublishedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: data.Published,
		},
		Published:       data.Published,
		Title:           data.Title,
		Excerpt:         pgtype.Text{String: data.Excerpt, Valid: true},
		MetaTitle:       pgtype.Text{String: data.MetaTitle, Valid: true},
		MetaDescription: pgtype.Text{String: data.MetaDescription, Valid: true},
		Slug:            slug.Make(data.Title),
		ImageLink:       pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:        pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:         pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.InsertArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

type UpdateArticleData struct {
	ID              int32
	UpdatedAt       time.Time
	Published       bool
	Title           string
	Excerpt         string
	MetaTitle       string
	MetaDescription string
	Slug            string
	ImageLink       string
	ReadTime        int32
	Content         string
}

func UpdateArticle(
	ctx context.Context,
	exec storage.Executor,
	data UpdateArticleData,
) (Article, error) {
	if err := Validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	current, err := FindArticle(ctx, exec, data.ID)
	if err != nil {
		return Article{}, err
	}

	firstPublishedAt := current.FirstPublishedAt
	if data.Published && firstPublishedAt.IsZero() {
		firstPublishedAt = time.Now().UTC()
	}

	params := db.UpdateArticleParams{
		ID: data.ID,
		FirstPublishedAt: pgtype.Timestamptz{
			Time:  firstPublishedAt,
			Valid: !firstPublishedAt.IsZero(),
		},
		Published:       data.Published,
		Title:           data.Title,
		Excerpt:         pgtype.Text{String: data.Excerpt, Valid: true},
		MetaTitle:       pgtype.Text{String: data.MetaTitle, Valid: true},
		MetaDescription: pgtype.Text{String: data.MetaDescription, Valid: true},
		Slug:            data.Slug,
		ImageLink:       pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:        pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:         pgtype.Text{String: data.Content, Valid: true},
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
	id int32,
) error {
	if err := clearArticleTagConnections(ctx, exec, id); err != nil {
		return err
	}

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

func AllPublishedArticles(
	ctx context.Context,
	exec storage.Executor,
) ([]Article, error) {
	rows, err := queries.QueryPublishedArticles(ctx, exec)
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
	if err := Validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	firstPublishedAt := data.FirstPublishedAt
	if data.Published && firstPublishedAt.IsZero() {
		firstPublishedAt = time.Now().UTC()
	}

	params := db.UpsertArticleParams{
		FirstPublishedAt: pgtype.Timestamptz{
			Time:  firstPublishedAt,
			Valid: !firstPublishedAt.IsZero(),
		},
		Published:       data.Published,
		Title:           data.Title,
		Excerpt:         pgtype.Text{String: data.Excerpt, Valid: true},
		MetaTitle:       pgtype.Text{String: data.MetaTitle, Valid: true},
		MetaDescription: pgtype.Text{String: data.MetaDescription, Valid: true},
		Slug:            data.Slug,
		ImageLink:       pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:        pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:         pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.UpsertArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row), nil
}

func CountArticles(
	ctx context.Context,
	exec storage.Executor,
) (int64, error) {
	return queries.CountArticles(ctx, exec)
}

func AttachTagsToArticle(
	ctx context.Context,
	exec storage.Executor,
	articleID int32,
	tagIDs []int32,
) error {
	seen := make(map[int32]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		if tagID <= 0 {
			continue
		}
		if _, exists := seen[tagID]; exists {
			continue
		}
		seen[tagID] = struct{}{}

		if _, err := queries.InsertArticleTagConnection(
			ctx,
			exec,
			db.InsertArticleTagConnectionParams{
				ArticleID: articleID,
				TagID:     tagID,
			},
		); err != nil {
			return err
		}
	}

	return nil
}

func TagIDsForArticle(
	ctx context.Context,
	exec storage.Executor,
	articleID int32,
) ([]int32, error) {
	connections, err := queries.QueryArticleTagConnection(ctx, exec)
	if err != nil {
		return nil, err
	}

	tagIDs := make([]int32, 0)
	for _, connection := range connections {
		if connection.ArticleID != articleID {
			continue
		}
		tagIDs = append(tagIDs, connection.TagID)
	}

	return tagIDs, nil
}

func ReplaceTagsForArticle(
	ctx context.Context,
	exec storage.Executor,
	articleID int32,
	tagIDs []int32,
) error {
	if err := clearArticleTagConnections(ctx, exec, articleID); err != nil {
		return err
	}

	return AttachTagsToArticle(ctx, exec, articleID, tagIDs)
}

func clearArticleTagConnections(
	ctx context.Context,
	exec storage.Executor,
	articleID int32,
) error {
	connections, err := queries.QueryArticleTagConnection(ctx, exec)
	if err != nil {
		return err
	}

	for _, connection := range connections {
		if connection.ArticleID != articleID {
			continue
		}
		if err := queries.DeleteArticleTagConnection(ctx, exec, connection.ID); err != nil {
			return err
		}
	}

	return nil
}

func rowToArticle(row db.Article) Article {
	return Article{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		FirstPublishedAt: row.FirstPublishedAt.Time,
		Published:        row.Published,
		Title:            row.Title,
		Excerpt:          row.Excerpt.String,
		MetaTitle:        row.MetaTitle.String,
		MetaDescription:  row.MetaDescription.String,
		Slug:             row.Slug,
		ImageLink:        row.ImageLink.String,
		ReadTime:         row.ReadTime.Int32,
		Content:          row.Content.String,
	}
}
