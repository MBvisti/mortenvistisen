package models

import (
	"context"
	"database/sql"
	"errors"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

const (
	articleTitleRecommendedMinLength           = 50
	articleTitleMaxLength                      = 60
	articleMetaTitleRecommendedMinLength       = 50
	articleMetaTitleMaxLength                  = 60
	articleMetaDescriptionRecommendedMinLength = 120
	articleMetaDescriptionMaxLength            = 160
)

type ArticleEntity struct {
	bun.BaseModel `bun:"table:articles,alias:articles"`

	ID               int32          `bun:"id,pk,autoincrement"`
	CreatedAt        time.Time      `bun:"created_at"`
	UpdatedAt        time.Time      `bun:"updated_at"`
	FirstPublishedAt sql.NullTime   `bun:"first_published_at"`
	Published        bool           `bun:"published"`
	Title            string         `bun:"title"`
	Excerpt          sql.NullString `bun:"excerpt"`
	MetaTitle        sql.NullString `bun:"meta_title"`
	MetaDescription  sql.NullString `bun:"meta_description"`
	Slug             string         `bun:"slug"`
	ImageLink        sql.NullString `bun:"image_link"`
	ReadTime         sql.NullInt32  `bun:"read_time"`
	Content          sql.NullString `bun:"content"`
	MetaImageLink    sql.NullString `bun:"meta_image_link"`
}

func (e *ArticleEntity) validationBuilder() *validation.ValidationBuilder {
	b := validation.NewBuilder()
	b.Required("title", e.Title)
	b.MaxLen("title", e.Title, articleTitleMaxLength)
	b.RecommendedLenBetween(
		"title",
		articleTitleRecommendedMinLength,
		articleTitleMaxLength,
	)
	b.Required("slug", e.Slug)
	b.RequiredWhen(e.Published, "firstPublishedAt", e.FirstPublishedAt)
	b.RequiredWhen(e.Published, "excerpt", e.Excerpt)
	b.RequiredWhen(e.Published, "metaTitle", e.MetaTitle)
	b.RequiredWhen(e.Published, "metaDescription", e.MetaDescription)
	b.RequiredWhen(e.Published, "imageLink", e.ImageLink)
	b.RequiredWhen(e.Published, "readTime", e.ReadTime)
	b.RequiredWhen(e.Published, "content", e.Content)
	b.MaxLen("metaTitle", e.MetaTitle, articleMetaTitleMaxLength)
	b.RecommendedLenBetween(
		"metaTitle",
		articleMetaTitleRecommendedMinLength,
		articleMetaTitleMaxLength,
	)
	b.MaxLen(
		"metaDescription",
		e.MetaDescription,
		articleMetaDescriptionMaxLength,
	)
	b.RecommendedLenBetween(
		"metaDescription",
		articleMetaDescriptionRecommendedMinLength,
		articleMetaDescriptionMaxLength,
	)
	b.URL("imageLink", e.ImageLink)
	b.URL("metaImageLink", e.MetaImageLink)
	if e.Published {
		b.MinInt("readTime", e.ReadTime, 1)
	} else {
		b.MinInt("readTime", e.ReadTime, 0)
	}

	return b
}

func (e *ArticleEntity) Validate() error {
	return e.validationBuilder().Err()
}

func ArticleValidationRules(published bool) validation.Rules {
	return (&ArticleEntity{Published: published}).validationBuilder().Rules()
}

func (a article) Find(ctx context.Context, db storage.Executor, id int32) (ArticleEntity, error) {
	var entity ArticleEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return ArticleEntity{}, err
	}

	return entity, nil
}

func (a article) FindPublishedBySlug(
	ctx context.Context,
	db storage.Executor,
	slug string,
) (ArticleEntity, error) {
	var entity ArticleEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("slug = ?", slug).
		Where("published = TRUE").
		Where("first_published_at IS NOT NULL").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ArticleEntity{}, ErrNotFound
		}
		return ArticleEntity{}, err
	}

	return entity, nil
}

func (a article) FindForUpdate(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (ArticleEntity, error) {
	var entity ArticleEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		For("UPDATE").
		Scan(ctx); err != nil {
		return ArticleEntity{}, err
	}

	return entity, nil
}

func (a article) FindForUpdateBySlug(
	ctx context.Context,
	db storage.Executor,
	slug string,
) (ArticleEntity, error) {
	var entity ArticleEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("slug = ?", slug).
		For("UPDATE").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ArticleEntity{}, ErrNotFound
		}
		return ArticleEntity{}, err
	}

	return entity, nil
}

type CreateArticleData struct {
	FirstPublishedAt sql.NullTime
	Published        bool
	Title            string
	Excerpt          sql.NullString
	MetaTitle        sql.NullString
	MetaDescription  sql.NullString
	Slug             string
	ImageLink        sql.NullString
	MetaImageLink    sql.NullString
	ReadTime         sql.NullInt32
	Content          sql.NullString
}

func (d CreateArticleData) Validate() error {
	entity := ArticleEntity{
		FirstPublishedAt: d.FirstPublishedAt,
		Published:        d.Published,
		Title:            d.Title,
		Excerpt:          normalizeOptionalString(d.Excerpt),
		MetaTitle:        normalizeOptionalString(d.MetaTitle),
		MetaDescription:  normalizeOptionalString(d.MetaDescription),
		Slug:             d.Slug,
		ImageLink:        normalizeOptionalString(d.ImageLink),
		MetaImageLink:    normalizeOptionalString(d.MetaImageLink),
		ReadTime:         d.ReadTime,
		Content:          normalizeOptionalString(d.Content),
	}

	if err := validation.Validate(&entity); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}

	return nil
}

func normalizeOptionalString(value sql.NullString) sql.NullString {
	if !value.Valid || strings.TrimSpace(value.String) == "" {
		return sql.NullString{}
	}

	return value
}

func (a article) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateArticleData,
) (ArticleEntity, error) {
	if err := data.Validate(); err != nil {
		return ArticleEntity{}, err
	}

	entity := ArticleEntity{
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		FirstPublishedAt: data.FirstPublishedAt,
		Published:        data.Published,
		Title:            data.Title,
		Excerpt:          normalizeOptionalString(data.Excerpt),
		MetaTitle:        normalizeOptionalString(data.MetaTitle),
		MetaDescription:  normalizeOptionalString(data.MetaDescription),
		Slug:             data.Slug,
		ImageLink:        normalizeOptionalString(data.ImageLink),
		MetaImageLink:    normalizeOptionalString(data.MetaImageLink),
		ReadTime:         data.ReadTime,
		Content:          normalizeOptionalString(data.Content),
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return ArticleEntity{}, err
	}

	return entity, nil
}

type UpdateArticleData struct {
	ID               int32
	UpdatedAt        time.Time
	FirstPublishedAt sql.NullTime
	Published        bool
	Title            string
	Excerpt          sql.NullString
	MetaTitle        sql.NullString
	MetaDescription  sql.NullString
	Slug             string
	ImageLink        sql.NullString
	MetaImageLink    sql.NullString
	ReadTime         sql.NullInt32
	Content          sql.NullString
}

func (d UpdateArticleData) Validate() error {
	entity := ArticleEntity{
		ID:               d.ID,
		FirstPublishedAt: d.FirstPublishedAt,
		Published:        d.Published,
		Title:            d.Title,
		Excerpt:          normalizeOptionalString(d.Excerpt),
		MetaTitle:        normalizeOptionalString(d.MetaTitle),
		MetaDescription:  normalizeOptionalString(d.MetaDescription),
		Slug:             d.Slug,
		ImageLink:        normalizeOptionalString(d.ImageLink),
		MetaImageLink:    normalizeOptionalString(d.MetaImageLink),
		ReadTime:         d.ReadTime,
		Content:          normalizeOptionalString(d.Content),
	}

	if err := validation.Validate(&entity); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}

	return nil
}

func (a article) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateArticleData,
) (ArticleEntity, error) {
	if err := data.Validate(); err != nil {
		return ArticleEntity{}, err
	}

	entity := ArticleEntity{
		ID:               data.ID,
		UpdatedAt:        time.Now(),
		FirstPublishedAt: data.FirstPublishedAt,
		Published:        data.Published,
		Title:            data.Title,
		Excerpt:          normalizeOptionalString(data.Excerpt),
		MetaTitle:        normalizeOptionalString(data.MetaTitle),
		MetaDescription:  normalizeOptionalString(data.MetaDescription),
		Slug:             data.Slug,
		ImageLink:        normalizeOptionalString(data.ImageLink),
		MetaImageLink:    normalizeOptionalString(data.MetaImageLink),
		ReadTime:         data.ReadTime,
		Content:          normalizeOptionalString(data.Content),
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("first_published_at").
		Column("published").
		Column("title").
		Column("excerpt").
		Column("meta_title").
		Column("meta_description").
		Column("slug").
		Column("image_link").
		Column("meta_image_link").
		Column("read_time").
		Column("content").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		return ArticleEntity{}, err
	}

	return entity, nil
}

func (a article) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*ArticleEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (a article) All(ctx context.Context, db storage.Executor) ([]ArticleEntity, error) {
	var entities []ArticleEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

func (a article) LatestPublished(
	ctx context.Context,
	db storage.Executor,
	limit int,
) ([]ArticleEntity, error) {
	var entities []ArticleEntity
	if err := db.NewSelect().
		Model(&entities).
		Where("published = TRUE").
		Where("first_published_at IS NOT NULL").
		OrderExpr("first_published_at DESC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type ArticleFilter struct {
	Published *bool
}

type PaginatedArticles struct {
	Articles   []ArticleEntity
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func applyArticleFilter(query *bun.SelectQuery, filter ArticleFilter) *bun.SelectQuery {
	if filter.Published != nil {
		query = query.Where("published = ?", *filter.Published)
	}
	return query
}

func (a article) Paginate(
	ctx context.Context,
	db storage.Executor,
	filter ArticleFilter,
	page, pageSize int64,
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

	totalCount, err := applyArticleFilter(
		db.NewSelect().Model((*ArticleEntity)(nil)),
		filter,
	).Count(ctx)
	if err != nil {
		return PaginatedArticles{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * pageSize

	entities := make([]ArticleEntity, 0, int(pageSize))
	if err := applyArticleFilter(
		db.NewSelect().Model(&entities),
		filter,
	).
		OrderExpr("updated_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedArticles{}, err
	}

	return PaginatedArticles{
		Articles:   entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (a article) PaginatePublished(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
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

	totalCount, err := db.NewSelect().
		Model(&ArticleEntity{}).
		Where("published = TRUE").
		Where("first_published_at IS NOT NULL").
		Count(ctx)
	if err != nil {
		return PaginatedArticles{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}

	entities := make([]ArticleEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Where("published = TRUE").
		Where("first_published_at IS NOT NULL").
		OrderExpr("first_published_at DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Scan(ctx); err != nil {
		return PaginatedArticles{}, err
	}

	return PaginatedArticles{
		Articles:   entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (a article) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateArticleData,
) (ArticleEntity, error) {
	entity := ArticleEntity{
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		FirstPublishedAt: data.FirstPublishedAt,
		Published:        data.Published,
		Title:            data.Title,
		Excerpt:          normalizeOptionalString(data.Excerpt),
		MetaTitle:        normalizeOptionalString(data.MetaTitle),
		MetaDescription:  normalizeOptionalString(data.MetaDescription),
		Slug:             data.Slug,
		ImageLink:        normalizeOptionalString(data.ImageLink),
		MetaImageLink:    normalizeOptionalString(data.MetaImageLink),
		ReadTime:         data.ReadTime,
		Content:          normalizeOptionalString(data.Content),
	}

	if err := validation.Validate(&entity); err != nil {
		return ArticleEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("first_published_at = excluded.first_published_at").
		Set("published = excluded.published").
		Set("title = excluded.title").
		Set("excerpt = excluded.excerpt").
		Set("meta_title = excluded.meta_title").
		Set("meta_description = excluded.meta_description").
		Set("slug = excluded.slug").
		Set("image_link = excluded.image_link").
		Set("meta_image_link = excluded.meta_image_link").
		Set("read_time = excluded.read_time").
		Set("content = excluded.content").
		Returning("*").
		Scan(ctx); err != nil {
		return ArticleEntity{}, err
	}

	return entity, nil
}
