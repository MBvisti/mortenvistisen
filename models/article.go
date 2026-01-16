package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
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
	Published        bool
	Tags             []string
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

	tags, err := FindTagsForArticle(ctx, exec, row.ID)
	if err != nil {
		return Article{}, err
	}

	tagNames := make([]string, len(tags))
	for j, tag := range tags {
		tagNames[j] = tag.Title
	}

	return rowToArticle(row, tagNames), nil
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

	tags, err := FindTagsForArticle(ctx, exec, row.ID)
	if err != nil {
		return Article{}, err
	}

	tagNames := make([]string, len(tags))
	for j, tag := range tags {
		tagNames[j] = tag.Title
	}

	return rowToArticle(row, tagNames), nil
}

type CreateArticleData struct {
	FirstPublishedAt time.Time
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	ImageLink        string
	Published        bool
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

	// If published but no first_published_at, auto-set to now
	firstPublishedAt := data.FirstPublishedAt
	if data.Published && firstPublishedAt.IsZero() {
		firstPublishedAt = time.Now()
	}

	params := db.InsertArticleParams{
		ID: uuid.New(),
		FirstPublishedAt: pgtype.Timestamptz{
			Time:  firstPublishedAt,
			Valid: !firstPublishedAt.IsZero(),
		},
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Published:       data.Published,
		Slug:            slug.Make(data.Title),
		ImageLink:       pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:        pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:         pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.InsertArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	return rowToArticle(row, nil), nil
}

type UpdateArticleData struct {
	ID               uuid.UUID
	UpdatedAt        time.Time
	FirstPublishedAt time.Time
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	// Slug             string
	ImageLink string
	Published bool
	ReadTime  int32
	Content   string
}

func UpdateArticle(
	ctx context.Context,
	exec storage.Executor,
	data UpdateArticleData,
) (Article, error) {
	if err := Validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	// If published but no first_published_at, auto-set to now
	firstPublishedAt := data.FirstPublishedAt
	if data.Published && firstPublishedAt.IsZero() {
		firstPublishedAt = time.Now()
	}

	params := db.UpdateArticleParams{
		ID: data.ID,
		FirstPublishedAt: pgtype.Timestamptz{
			Time:  firstPublishedAt,
			Valid: !firstPublishedAt.IsZero(),
		},
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            slug.Make(data.Title),
		Published:       data.Published,
		ImageLink:       pgtype.Text{String: data.ImageLink, Valid: data.ImageLink != ""},
		ReadTime:        pgtype.Int4{Int32: data.ReadTime, Valid: data.ReadTime != 0},
		Content:         pgtype.Text{String: data.Content, Valid: data.Content != ""},
	}

	row, err := queries.UpdateArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	tags, err := FindTagsForArticle(ctx, exec, row.ID)
	if err != nil {
		return Article{}, err
	}

	tagNames := make([]string, len(tags))
	for j, tag := range tags {
		tagNames[j] = tag.Title
	}

	return rowToArticle(row, tagNames), nil
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
		tags, err := FindTagsForArticle(ctx, exec, row.ID)
		if err != nil {
			return nil, err
		}

		tagNames := make([]string, len(tags))
		for j, tag := range tags {
			tagNames[j] = tag.Title
		}

		articles[i] = rowToArticle(row, tagNames)
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

		tags, err := FindTagsForArticle(ctx, exec, row.ID)
		if err != nil {
			return PaginatedArticles{}, err
		}

		tagNames := make([]string, len(tags))
		for j, tag := range tags {
			tagNames[j] = tag.Title
		}

		articles[i] = rowToArticle(row, tagNames)
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

	params := db.UpsertArticleParams{
		ID:               uuid.New(),
		FirstPublishedAt: pgtype.Timestamptz{Time: data.FirstPublishedAt, Valid: true},
		Title:            data.Title,
		Excerpt:          data.Excerpt,
		MetaTitle:        data.MetaTitle,
		MetaDescription:  data.MetaDescription,
		Published:        data.Published,
		Slug:             slug.Make(data.Title),
		ImageLink:        pgtype.Text{String: data.ImageLink, Valid: true},
		ReadTime:         pgtype.Int4{Int32: data.ReadTime, Valid: true},
		Content:          pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.UpsertArticle(ctx, exec, params)
	if err != nil {
		return Article{}, err
	}

	tags, err := FindTagsForArticle(ctx, exec, row.ID)
	if err != nil {
		return Article{}, err
	}

	tagNames := make([]string, len(tags))
	for j, tag := range tags {
		tagNames[j] = tag.Title
	}

	return rowToArticle(row, tagNames), nil
}

func AssociateTagsWithArticle(
	ctx context.Context,
	exec storage.Executor,
	articleID uuid.UUID,
	tagIDs []uuid.UUID,
) error {
	for _, tagID := range tagIDs {
		params := db.InsertArticleTagConnectionParams{
			ID:        uuid.New(),
			ArticleID: articleID,
			TagID:     tagID,
		}
		_, err := queries.InsertArticleTagConnection(ctx, exec, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindPublishedArticles(
	ctx context.Context,
	exec storage.Executor,
) ([]Article, error) {
	rows, err := queries.QueryPublishedArticles(ctx, exec)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		tags, err := FindTagsForArticle(ctx, exec, row.ID)
		if err != nil {
			return nil, err
		}

		tagNames := make([]string, len(tags))
		for j, tag := range tags {
			tagNames[j] = tag.Title
		}

		articles[i] = rowToArticle(row, tagNames)
	}

	return articles, nil
}

func FindTagsForArticle(
	ctx context.Context,
	exec storage.Executor,
	articleID uuid.UUID,
) ([]Tag, error) {
	rows, err := queries.QueryTagsByArticleID(ctx, exec, articleID)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, len(rows))
	for i, row := range rows {
		tags[i] = rowToTag(row)
	}

	return tags, nil
}

func rowToArticle(row db.Article, tags []string) Article {
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
		Published:        row.Published,
		Content:          row.Content.String,
		Tags:             tags,
	}
}
