package models

import (
	"context"
	"errors"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
)

const tagTitleMaxLength = 255

type TagEntity struct {
	bun.BaseModel `bun:"table:tags,alias:tags"`
	ID            int32     `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	Title         string    `bun:"title"`
}

func (e *TagEntity) validationBuilder() *validation.ValidationBuilder {
	b := validation.NewBuilder()
	b.Required("title", e.Title)
	b.MaxLen("title", e.Title, tagTitleMaxLength)

	return b
}

func (e *TagEntity) Validate() error {
	return e.validationBuilder().Err()
}

func TagValidationRules() validation.Rules {
	return (&TagEntity{}).validationBuilder().Rules()
}

func (t tag) Find(ctx context.Context, db storage.Executor, id int32) (TagEntity, error) {
	var entity TagEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return TagEntity{}, err
	}

	return entity, nil
}

func (t tag) FindByTitle(
	ctx context.Context,
	db storage.Executor,
	title string,
) (TagEntity, error) {
	var entity TagEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("LOWER(BTRIM(title)) = LOWER(BTRIM(?))", title).
		Scan(ctx); err != nil {
		return TagEntity{}, err
	}

	return entity, nil
}

func (t tag) FindByIDs(ctx context.Context, db storage.Executor, ids []int32) ([]TagEntity, error) {
	if len(ids) == 0 {
		return []TagEntity{}, nil
	}

	entities := make([]TagEntity, 0, len(ids))
	if err := db.NewSelect().
		Model(&entities).
		Where("id IN (?)", bun.In(ids)).
		OrderExpr("LOWER(title) ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type CreateTagData struct {
	Title string
}

func (t tag) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateTagData,
) (TagEntity, error) {
	entity := TagEntity{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     strings.TrimSpace(data.Title),
	}

	if err := validation.Validate(&entity); err != nil {
		return TagEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return TagEntity{}, tagPersistenceError(err)
	}

	return entity, nil
}

type UpdateTagData struct {
	ID        int32
	UpdatedAt time.Time
	Title     string
}

func (t tag) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateTagData,
) (TagEntity, error) {
	entity := TagEntity{
		ID:        data.ID,
		UpdatedAt: time.Now(),
		Title:     strings.TrimSpace(data.Title),
	}

	if err := validation.Validate(&entity); err != nil {
		return TagEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("title").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		return TagEntity{}, tagPersistenceError(err)
	}

	return entity, nil
}

func tagPersistenceError(err error) error {
	var postgresErr *pgconn.PgError
	if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
		validationErr := validation.ValidationErrors{{
			Field:   "title",
			Code:    "unique",
			Message: "A tag with this title already exists.",
		}}
		return errors.Join(ErrDomainValidation, validationErr)
	}

	return err
}

func (t tag) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*TagEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (t tag) All(ctx context.Context, db storage.Executor) ([]TagEntity, error) {
	var entities []TagEntity
	if err := db.NewSelect().
		Model(&entities).
		OrderExpr("LOWER(title) ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type TagWithUsage struct {
	ID           int32     `bun:"id"`
	CreatedAt    time.Time `bun:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at"`
	Title        string    `bun:"title"`
	ArticleCount int64     `bun:"article_count"`
}

type PaginatedTagsWithUsage struct {
	Tags       []TagWithUsage
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func (t tag) PaginateWithUsage(
	ctx context.Context,
	db storage.Executor,
	page int64,
	pageSize int64,
) (PaginatedTagsWithUsage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	totalCount, err := db.NewSelect().Model((*TagEntity)(nil)).Count(ctx)
	if err != nil {
		return PaginatedTagsWithUsage{}, err
	}
	totalPages := (int64(totalCount) + pageSize - 1) / pageSize
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}

	items := make([]TagWithUsage, 0)
	if err := db.NewSelect().
		TableExpr("tags AS tags").
		Column("tags.id", "tags.created_at", "tags.updated_at", "tags.title").
		ColumnExpr("COUNT(article_tag_connections.id) AS article_count").
		Join("LEFT JOIN article_tag_connections ON article_tag_connections.tag_id = tags.id").
		Group("tags.id").
		OrderExpr("LOWER(tags.title) ASC").
		Limit(int(pageSize)).
		Offset(int((page-1)*pageSize)).
		Scan(ctx, &items); err != nil {
		return PaginatedTagsWithUsage{}, err
	}

	return PaginatedTagsWithUsage{
		Tags:       items,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (t tag) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateTagData,
) (TagEntity, error) {
	entity := TagEntity{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     strings.TrimSpace(data.Title),
	}

	if err := validation.Validate(&entity); err != nil {
		return TagEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT ((LOWER(BTRIM(title)))) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("updated_at = EXCLUDED.updated_at").
		Returning("*").
		Scan(ctx); err != nil {
		return TagEntity{}, err
	}

	return entity, nil
}
