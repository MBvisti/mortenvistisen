package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type Article struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PublishedAt     *time.Time
	Title           string
	Excerpt         string
	MetaTitle       string
	MetaDescription string
	Slug            string
	ImageLink       *string
	Content         *string
}

func (a Article) IsPublished() bool {
	return a.PublishedAt != nil && !a.PublishedAt.IsZero()
}

func (a Article) IsDraft() bool {
	return !a.IsPublished()
}

func GetArticleByID(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Article, error) {
	row, err := db.Stmts.QueryArticleByID(ctx, dbtx, id)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

func GetArticleByTitle(
	ctx context.Context,
	dbtx db.DBTX,
	title string,
) (Article, error) {
	row, err := db.Stmts.QueryArticleByTitle(ctx, dbtx, title)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

func GetArticleBySlug(
	ctx context.Context,
	dbtx db.DBTX,
	slug string,
) (Article, error) {
	row, err := db.Stmts.QueryArticleBySlug(ctx, dbtx, slug)
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

func GetArticles(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Article, error) {
	rows, err := db.Stmts.QueryArticles(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:              row.ID,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
			PublishedAt:     nullTimeToPointer(row.PublishedAt),
			Title:           row.Title,
			Excerpt:         row.Excerpt,
			MetaTitle:       row.MetaTitle,
			MetaDescription: row.MetaDescription,
			Slug:            row.Slug,
			ImageLink:       nullStringToPointer(row.ImageLink),
			Content:         nullStringToPointer(row.Content),
		}
	}

	return articles, nil
}

func GetPublishedArticles(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Article, error) {
	rows, err := db.Stmts.QueryPublishedArticles(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:              row.ID,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
			PublishedAt:     nullTimeToPointer(row.PublishedAt),
			Title:           row.Title,
			Excerpt:         row.Excerpt,
			MetaTitle:       row.MetaTitle,
			MetaDescription: row.MetaDescription,
			Slug:            row.Slug,
			ImageLink:       nullStringToPointer(row.ImageLink),
			Content:         nullStringToPointer(row.Content),
		}
	}

	return articles, nil
}

func GetDraftArticles(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Article, error) {
	rows, err := db.Stmts.QueryDraftArticles(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:              row.ID,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
			PublishedAt:     nullTimeToPointer(row.PublishedAt),
			Title:           row.Title,
			Excerpt:         row.Excerpt,
			MetaTitle:       row.MetaTitle,
			MetaDescription: row.MetaDescription,
			Slug:            row.Slug,
			ImageLink:       nullStringToPointer(row.ImageLink),
			Content:         nullStringToPointer(row.Content),
		}
	}

	return articles, nil
}

type PaginationResult struct {
	Articles    []Article
	TotalCount  int64
	Page        int
	PageSize    int
	TotalPages  int
	HasNext     bool
	HasPrevious bool
}

func GetArticlesPaginated(
	ctx context.Context,
	dbtx db.DBTX,
	page int,
	pageSize int,
) (PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Limit max page size
	}

	offset := (page - 1) * pageSize

	// Get total count
	totalCount, err := db.Stmts.CountArticles(ctx, dbtx)
	if err != nil {
		return PaginationResult{}, err
	}

	// Get paginated articles
	rows, err := db.Stmts.QueryArticlesPaginated(ctx, dbtx, db.QueryArticlesPaginatedParams{
		Limit:  int32(pageSize), //nolint:gosec // pageSize is bounded above
		Offset: int32(offset),   //nolint:gosec // offset is calculated from bounded values
	})
	if err != nil {
		return PaginationResult{}, err
	}

	articles := make([]Article, len(rows))
	for i, row := range rows {
		articles[i] = Article{
			ID:              row.ID,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
			PublishedAt:     nullTimeToPointer(row.PublishedAt),
			Title:           row.Title,
			Excerpt:         row.Excerpt,
			MetaTitle:       row.MetaTitle,
			MetaDescription: row.MetaDescription,
			Slug:            row.Slug,
			ImageLink:       nullStringToPointer(row.ImageLink),
			Content:         nullStringToPointer(row.Content),
		}
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return PaginationResult{
		Articles:    articles,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

type NewArticlePayload struct {
	Title           string  `validate:"required,max=100"`
	Excerpt         string  `validate:"required,max=255"`
	MetaTitle       string  `validate:"required,max=100"`
	MetaDescription string  `validate:"required,max=100"`
	Slug            string  `validate:"required,max=255"`
	ImageLink       *string `validate:"omitempty,max=255"`
	Content         *string
}

func NewArticle(
	ctx context.Context,
	dbtx db.DBTX,
	data NewArticlePayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	article := Article{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		ImageLink:       data.ImageLink,
		Content:         data.Content,
	}

	_, err := db.Stmts.InsertArticle(ctx, dbtx, db.InsertArticleParams{
		ID:              article.ID,
		CreatedAt:       pgtype.Timestamptz{Time: article.CreatedAt, Valid: true},
		UpdatedAt:       pgtype.Timestamptz{Time: article.UpdatedAt, Valid: true},
		PublishedAt:     pgtype.Timestamptz{Valid: false},
		Title:           article.Title,
		Excerpt:         article.Excerpt,
		MetaTitle:       article.MetaTitle,
		MetaDescription: article.MetaDescription,
		Slug:            article.Slug,
		ImageLink:       stringPointerToNullString(article.ImageLink),
		Content:         stringPointerToNullString(article.Content),
	})
	if err != nil {
		return Article{}, err
	}

	return article, nil
}

type UpdateArticlePayload struct {
	ID              uuid.UUID `validate:"required,uuid"`
	UpdatedAt       time.Time `validate:"required"`
	PublishedAt     *time.Time
	Title           string  `validate:"required,max=100"`
	Excerpt         string  `validate:"required,max=255"`
	MetaTitle       string  `validate:"required,max=100"`
	MetaDescription string  `validate:"required,max=100"`
	Slug            string  `validate:"required,max=255"`
	ImageLink       *string `validate:"omitempty,max=255"`
	Content         *string
}

func UpdateArticle(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateArticlePayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateArticle(ctx, dbtx, db.UpdateArticleParams{
		ID:              data.ID,
		UpdatedAt:       pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		PublishedAt:     timePointerToNullTime(data.PublishedAt),
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		ImageLink:       stringPointerToNullString(data.ImageLink),
		Content:         stringPointerToNullString(data.Content),
	})
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

type UpdateArticleContentPayload struct {
	ID        uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
	Content   *string
}

func UpdateArticleContent(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateArticleContentPayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateArticleContent(ctx, dbtx, db.UpdateArticleContentParams{
		ID:        data.ID,
		UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		Content:   stringPointerToNullString(data.Content),
	})
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

type UpdateArticleMetadataPayload struct {
	ID              uuid.UUID `validate:"required,uuid"`
	UpdatedAt       time.Time `validate:"required"`
	Title           string    `validate:"required,max=100"`
	Excerpt         string    `validate:"required,max=255"`
	MetaTitle       string    `validate:"required,max=100"`
	MetaDescription string    `validate:"required,max=100"`
	Slug            string    `validate:"required,max=255"`
	ImageLink       *string   `validate:"omitempty,max=255"`
}

func UpdateArticleMetadata(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateArticleMetadataPayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateArticleMetadata(ctx, dbtx, db.UpdateArticleMetadataParams{
		ID:              data.ID,
		UpdatedAt:       pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		Title:           data.Title,
		Excerpt:         data.Excerpt,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		ImageLink:       stringPointerToNullString(data.ImageLink),
	})
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

type PublishArticlePayload struct {
	ID          uuid.UUID `validate:"required,uuid"`
	UpdatedAt   time.Time `validate:"required"`
	PublishedAt time.Time `validate:"required"`
}

func PublishArticle(
	ctx context.Context,
	dbtx db.DBTX,
	data PublishArticlePayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.PublishArticle(ctx, dbtx, db.PublishArticleParams{
		ID:          data.ID,
		UpdatedAt:   pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		PublishedAt: pgtype.Timestamptz{Time: data.PublishedAt, Valid: true},
	})
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

type UnpublishArticlePayload struct {
	ID        uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
}

func UnpublishArticle(
	ctx context.Context,
	dbtx db.DBTX,
	data UnpublishArticlePayload,
) (Article, error) {
	if err := validate.Struct(data); err != nil {
		return Article{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UnpublishArticle(ctx, dbtx, db.UnpublishArticleParams{
		ID:        data.ID,
		UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
	})
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		PublishedAt:     nullTimeToPointer(row.PublishedAt),
		Title:           row.Title,
		Excerpt:         row.Excerpt,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		Slug:            row.Slug,
		ImageLink:       nullStringToPointer(row.ImageLink),
		Content:         nullStringToPointer(row.Content),
	}, nil
}

func DeleteArticle(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.Stmts.DeleteArticle(ctx, dbtx, id)
}

// Helper functions for handling nullable types
func nullTimeToPointer(nt pgtype.Timestamptz) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func timePointerToNullTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func nullStringToPointer(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func stringPointerToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}
