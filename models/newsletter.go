package models

import (
	"context"
	"database/sql"
	"errors"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"
	"time"

	"github.com/uptrace/bun"
)

const (
	newsletterTitleRecommendedMinLength           = 50
	newsletterTitleMaxLength                      = 60
	newsletterMetaTitleRecommendedMinLength       = 50
	newsletterMetaTitleMaxLength                  = 60
	newsletterMetaDescriptionRecommendedMinLength = 120
	newsletterMetaDescriptionMaxLength            = 160
)

type NewsletterEntity struct {
	bun.BaseModel `bun:"table:newsletters,alias:newsletters"`

	ID              int32          `bun:"id,pk,autoincrement"`
	CreatedAt       time.Time      `bun:"created_at"`
	UpdatedAt       time.Time      `bun:"updated_at"`
	Title           string         `bun:"title"`
	Slug            sql.NullString `bun:"slug"`
	MetaTitle       string         `bun:"meta_title"`
	MetaDescription string         `bun:"meta_description"`
	IsPublished     sql.NullBool   `bun:"is_published"`
	ReleasedAt      sql.NullTime   `bun:"released_at"`
	Content         sql.NullString `bun:"content"`
	ImageLink       sql.NullString `bun:"image_link"`
	MetaImageLink   sql.NullString `bun:"meta_image_link"`
}

func (e *NewsletterEntity) validationBuilder() *validation.ValidationBuilder {
	b := validation.NewBuilder()
	b.Required("title", e.Title)
	b.MaxLen("title", e.Title, newsletterTitleMaxLength)
	b.RecommendedLenBetween(
		"title",
		newsletterTitleRecommendedMinLength,
		newsletterTitleMaxLength,
	)
	b.Required("slug", e.Slug)
	b.RequiredWhen(e.IsPublished.Bool, "releasedAt", e.ReleasedAt)
	b.RequiredWhen(e.IsPublished.Bool, "metaTitle", e.MetaTitle)
	b.RequiredWhen(e.IsPublished.Bool, "metaDescription", e.MetaDescription)
	b.RequiredWhen(e.IsPublished.Bool, "imageLink", e.ImageLink)
	b.RequiredWhen(e.IsPublished.Bool, "content", e.Content)
	b.URL("imageLink", e.ImageLink)
	b.URL("metaImageLink", e.MetaImageLink)
	b.MaxLen("metaTitle", e.MetaTitle, newsletterMetaTitleMaxLength)
	b.RecommendedLenBetween(
		"metaTitle",
		newsletterMetaTitleRecommendedMinLength,
		newsletterMetaTitleMaxLength,
	)
	b.MaxLen(
		"metaDescription",
		e.MetaDescription,
		newsletterMetaDescriptionMaxLength,
	)
	b.RecommendedLenBetween(
		"metaDescription",
		newsletterMetaDescriptionRecommendedMinLength,
		newsletterMetaDescriptionMaxLength,
	)

	return b
}

func (e *NewsletterEntity) Validate() error {
	return e.validationBuilder().Err()
}

func NewsletterValidationRules(published bool) validation.Rules {
	return (&NewsletterEntity{
		IsPublished: sql.NullBool{Bool: published, Valid: true},
	}).validationBuilder().Rules()
}

func (n newsletter) Find(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (NewsletterEntity, error) {
	var entity NewsletterEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return NewsletterEntity{}, err
	}

	return entity, nil
}

func (n newsletter) FindPublishedBySlug(
	ctx context.Context,
	db storage.Executor,
	slug string,
) (NewsletterEntity, error) {
	var entity NewsletterEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("slug = ?", slug).
		Where("is_published = TRUE").
		Where("released_at IS NOT NULL").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewsletterEntity{}, ErrNotFound
		}
		return NewsletterEntity{}, err
	}

	return entity, nil
}

func (n newsletter) FindForUpdate(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (NewsletterEntity, error) {
	var entity NewsletterEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		For("UPDATE").
		Scan(ctx); err != nil {
		return NewsletterEntity{}, err
	}

	return entity, nil
}

type CreateNewsletterData struct {
	Title           string
	Slug            sql.NullString
	MetaTitle       string
	MetaDescription string
	IsPublished     sql.NullBool
	ReleasedAt      sql.NullTime
	ImageLink       sql.NullString
	MetaImageLink   sql.NullString
	Content         sql.NullString
}

func (d CreateNewsletterData) Validate() error {
	entity := NewsletterEntity{
		Title:           d.Title,
		Slug:            normalizeOptionalString(d.Slug),
		MetaTitle:       d.MetaTitle,
		MetaDescription: d.MetaDescription,
		IsPublished:     d.IsPublished,
		ReleasedAt:      d.ReleasedAt,
		ImageLink:       normalizeOptionalString(d.ImageLink),
		MetaImageLink:   normalizeOptionalString(d.MetaImageLink),
		Content:         normalizeOptionalString(d.Content),
	}
	if err := validation.Validate(&entity); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}
	return nil
}

func (n newsletter) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateNewsletterData,
) (NewsletterEntity, error) {
	if err := data.Validate(); err != nil {
		return NewsletterEntity{}, err
	}
	entity := NewsletterEntity{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Title:           data.Title,
		Slug:            normalizeOptionalString(data.Slug),
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		IsPublished:     data.IsPublished,
		ReleasedAt:      data.ReleasedAt,
		ImageLink:       normalizeOptionalString(data.ImageLink),
		MetaImageLink:   normalizeOptionalString(data.MetaImageLink),
		Content:         normalizeOptionalString(data.Content),
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return NewsletterEntity{}, err
	}

	return entity, nil
}

type UpdateNewsletterData struct {
	ID              int32
	UpdatedAt       time.Time
	Title           string
	Slug            sql.NullString
	MetaTitle       string
	MetaDescription string
	IsPublished     sql.NullBool
	ReleasedAt      sql.NullTime
	ImageLink       sql.NullString
	MetaImageLink   sql.NullString
	Content         sql.NullString
}

func (d UpdateNewsletterData) Validate() error {
	entity := NewsletterEntity{
		ID:              d.ID,
		Title:           d.Title,
		Slug:            normalizeOptionalString(d.Slug),
		MetaTitle:       d.MetaTitle,
		MetaDescription: d.MetaDescription,
		IsPublished:     d.IsPublished,
		ReleasedAt:      d.ReleasedAt,
		ImageLink:       normalizeOptionalString(d.ImageLink),
		MetaImageLink:   normalizeOptionalString(d.MetaImageLink),
		Content:         normalizeOptionalString(d.Content),
	}
	if err := validation.Validate(&entity); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}
	return nil
}

func (n newsletter) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateNewsletterData,
) (NewsletterEntity, error) {
	if err := data.Validate(); err != nil {
		return NewsletterEntity{}, err
	}
	entity := NewsletterEntity{
		ID:              data.ID,
		UpdatedAt:       time.Now(),
		Title:           data.Title,
		Slug:            normalizeOptionalString(data.Slug),
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		IsPublished:     data.IsPublished,
		ReleasedAt:      data.ReleasedAt,
		ImageLink:       normalizeOptionalString(data.ImageLink),
		MetaImageLink:   normalizeOptionalString(data.MetaImageLink),
		Content:         normalizeOptionalString(data.Content),
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("title").
		Column("slug").
		Column("meta_title").
		Column("meta_description").
		Column("is_published").
		Column("released_at").
		Column("image_link").
		Column("meta_image_link").
		Column("content").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		return NewsletterEntity{}, err
	}

	return entity, nil
}

func (n newsletter) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*NewsletterEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (n newsletter) All(ctx context.Context, db storage.Executor) ([]NewsletterEntity, error) {
	var entities []NewsletterEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

func (n newsletter) LatestPublished(
	ctx context.Context,
	db storage.Executor,
	limit int,
) ([]NewsletterEntity, error) {
	var entities []NewsletterEntity
	if err := db.NewSelect().
		Model(&entities).
		Where("is_published = TRUE").
		Where("released_at IS NOT NULL").
		OrderExpr("released_at DESC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type PaginatedNewsletters struct {
	Newsletters []NewsletterEntity
	TotalCount  int64
	Page        int64
	PageSize    int64
	TotalPages  int64
}

func (n newsletter) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedNewsletters, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	totalCount, err := db.NewSelect().
		Model(&NewsletterEntity{}).Count(ctx)
	if err != nil {
		return PaginatedNewsletters{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * pageSize

	entities := make([]NewsletterEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		OrderExpr("updated_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedNewsletters{}, err
	}

	return PaginatedNewsletters{
		Newsletters: entities,
		TotalCount:  int64(totalCount),
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

func (n newsletter) PaginatePublished(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedNewsletters, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	totalCount, err := db.NewSelect().
		Model(&NewsletterEntity{}).
		Where("is_published = TRUE").
		Where("released_at IS NOT NULL").
		Count(ctx)
	if err != nil {
		return PaginatedNewsletters{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}

	entities := make([]NewsletterEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Where("is_published = TRUE").
		Where("released_at IS NOT NULL").
		OrderExpr("released_at DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Scan(ctx); err != nil {
		return PaginatedNewsletters{}, err
	}

	return PaginatedNewsletters{
		Newsletters: entities,
		TotalCount:  int64(totalCount),
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

func (n newsletter) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateNewsletterData,
) (NewsletterEntity, error) {
	entity := NewsletterEntity{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Title:           data.Title,
		Slug:            normalizeOptionalString(data.Slug),
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		IsPublished:     data.IsPublished,
		ReleasedAt:      data.ReleasedAt,
		ImageLink:       normalizeOptionalString(data.ImageLink),
		MetaImageLink:   normalizeOptionalString(data.MetaImageLink),
		Content:         normalizeOptionalString(data.Content),
	}

	if err := validation.Validate(&entity); err != nil {
		return NewsletterEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("title = excluded.title").
		Set("slug = excluded.slug").
		Set("meta_title = excluded.meta_title").
		Set("meta_description = excluded.meta_description").
		Set("is_published = excluded.is_published").
		Set("released_at = excluded.released_at").
		Set("image_link = excluded.image_link").
		Set("meta_image_link = excluded.meta_image_link").
		Set("content = excluded.content").
		Returning("*").
		Scan(ctx); err != nil {
		return NewsletterEntity{}, err
	}

	return entity, nil
}
