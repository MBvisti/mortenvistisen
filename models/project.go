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
	projectTitleMaxLength                      = 120
	projectMetaTitleRecommendedMinLength       = 50
	projectMetaTitleMaxLength                  = 60
	projectMetaDescriptionRecommendedMinLength = 120
	projectMetaDescriptionMaxLength            = 160
	projectStatusMaxLength                     = 80
	projectURLMaxLength                        = 255
)

type ProjectEntity struct {
	bun.BaseModel `bun:"table:projects,alias:projects"`

	ID              int32          `bun:"id,pk,autoincrement"`
	CreatedAt       time.Time      `bun:"created_at"`
	UpdatedAt       time.Time      `bun:"updated_at"`
	Published       bool           `bun:"published"`
	Title           string         `bun:"title"`
	Slug            string         `bun:"slug"`
	StartedAt       sql.NullTime   `bun:"started_at"`
	Status          string         `bun:"status"`
	Description     string         `bun:"description"`
	Content         string         `bun:"content"`
	ProjectUrl      sql.NullString `bun:"project_url"`
	MetaTitle       string         `bun:"meta_title"`
	MetaDescription string         `bun:"meta_description"`
	ImageLink       sql.NullString `bun:"image_link"`
	MetaImageLink   sql.NullString `bun:"meta_image_link"`
}

func (e *ProjectEntity) validationBuilder() *validation.ValidationBuilder {
	b := validation.NewBuilder()
	b.Required("title", e.Title)
	b.MaxLen("title", e.Title, projectTitleMaxLength)
	b.Required("slug", e.Slug)
	b.RequiredWhen(e.Published, "startedAt", e.StartedAt)
	b.RequiredWhen(e.Published, "status", e.Status)
	b.RequiredWhen(e.Published, "description", e.Description)
	b.RequiredWhen(e.Published, "metaTitle", e.MetaTitle)
	b.RequiredWhen(e.Published, "metaDescription", e.MetaDescription)
	b.RequiredWhen(e.Published, "imageLink", e.ImageLink)
	b.RequiredWhen(e.Published, "content", e.Content)
	b.URL("imageLink", e.ImageLink)
	b.URL("metaImageLink", e.MetaImageLink)
	b.MaxLen("status", e.Status, projectStatusMaxLength)
	b.MaxLen("metaTitle", e.MetaTitle, projectMetaTitleMaxLength)
	b.RecommendedLenBetween(
		"metaTitle",
		projectMetaTitleRecommendedMinLength,
		projectMetaTitleMaxLength,
	)
	b.MaxLen(
		"metaDescription",
		e.MetaDescription,
		projectMetaDescriptionMaxLength,
	)
	b.RecommendedLenBetween(
		"metaDescription",
		projectMetaDescriptionRecommendedMinLength,
		projectMetaDescriptionMaxLength,
	)
	b.MaxLen("projectUrl", e.ProjectUrl, projectURLMaxLength)
	b.URL("projectUrl", e.ProjectUrl)

	return b
}

func (e *ProjectEntity) Validate() error {
	return e.validationBuilder().Err()
}

func ProjectValidationRules(published bool) validation.Rules {
	return (&ProjectEntity{Published: published}).validationBuilder().Rules()
}

func (p project) Find(ctx context.Context, db storage.Executor, id int32) (ProjectEntity, error) {
	var entity ProjectEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return ProjectEntity{}, err
	}

	return entity, nil
}

func (p project) FindPublishedBySlug(
	ctx context.Context,
	db storage.Executor,
	slug string,
) (ProjectEntity, error) {
	var entity ProjectEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("slug = ?", slug).
		Where("published = TRUE").
		Where("started_at IS NOT NULL").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectEntity{}, ErrNotFound
		}
		return ProjectEntity{}, err
	}

	return entity, nil
}

func (p project) FindForUpdate(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (ProjectEntity, error) {
	var entity ProjectEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		For("UPDATE").
		Scan(ctx); err != nil {
		return ProjectEntity{}, err
	}

	return entity, nil
}

type CreateProjectData struct {
	Published       bool
	Title           string
	Slug            string
	StartedAt       sql.NullTime
	Status          string
	Description     string
	MetaTitle       string
	MetaDescription string
	ImageLink       sql.NullString
	MetaImageLink   sql.NullString
	Content         string
	ProjectUrl      sql.NullString
}

func (d CreateProjectData) Validate() error {
	entity := ProjectEntity{
		Published:       d.Published,
		Title:           d.Title,
		Slug:            d.Slug,
		StartedAt:       d.StartedAt,
		Status:          d.Status,
		Description:     d.Description,
		MetaTitle:       d.MetaTitle,
		MetaDescription: d.MetaDescription,
		ImageLink:       normalizeOptionalString(d.ImageLink),
		MetaImageLink:   normalizeOptionalString(d.MetaImageLink),
		Content:         d.Content,
		ProjectUrl:      normalizeOptionalString(d.ProjectUrl),
	}
	if err := validation.Validate(&entity); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}
	return nil
}

func (p project) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateProjectData,
) (ProjectEntity, error) {
	entity := ProjectEntity{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Published:       data.Published,
		Title:           data.Title,
		Slug:            data.Slug,
		StartedAt:       data.StartedAt,
		Status:          data.Status,
		Description:     data.Description,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		ImageLink:       normalizeOptionalString(data.ImageLink),
		MetaImageLink:   normalizeOptionalString(data.MetaImageLink),
		Content:         data.Content,
		ProjectUrl:      normalizeOptionalString(data.ProjectUrl),
	}

	if err := validation.Validate(&entity); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return ProjectEntity{}, err
	}

	return entity, nil
}

type UpdateProjectData struct {
	ID              int32
	UpdatedAt       time.Time
	Published       bool
	Title           string
	Slug            string
	StartedAt       sql.NullTime
	Status          string
	Description     string
	MetaTitle       string
	MetaDescription string
	ImageLink       sql.NullString
	MetaImageLink   sql.NullString
	Content         string
	ProjectUrl      sql.NullString
}

func (d UpdateProjectData) Validate() error {
	entity := ProjectEntity{
		ID:              d.ID,
		Published:       d.Published,
		Title:           d.Title,
		Slug:            d.Slug,
		StartedAt:       d.StartedAt,
		Status:          d.Status,
		Description:     d.Description,
		MetaTitle:       d.MetaTitle,
		MetaDescription: d.MetaDescription,
		ImageLink:       normalizeOptionalString(d.ImageLink),
		MetaImageLink:   normalizeOptionalString(d.MetaImageLink),
		Content:         d.Content,
		ProjectUrl:      normalizeOptionalString(d.ProjectUrl),
	}
	if err := validation.Validate(&entity); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}
	return nil
}

func (p project) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateProjectData,
) (ProjectEntity, error) {
	entity := ProjectEntity{
		ID:              data.ID,
		UpdatedAt:       time.Now(),
		Published:       data.Published,
		Title:           data.Title,
		Slug:            data.Slug,
		StartedAt:       data.StartedAt,
		Status:          data.Status,
		Description:     data.Description,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		ImageLink:       normalizeOptionalString(data.ImageLink),
		MetaImageLink:   normalizeOptionalString(data.MetaImageLink),
		Content:         data.Content,
		ProjectUrl:      normalizeOptionalString(data.ProjectUrl),
	}

	if err := validation.Validate(&entity); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("published").
		Column("title").
		Column("slug").
		Column("started_at").
		Column("status").
		Column("description").
		Column("meta_title").
		Column("meta_description").
		Column("image_link").
		Column("meta_image_link").
		Column("content").
		Column("project_url").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		return ProjectEntity{}, err
	}

	return entity, nil
}

func (p project) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*ProjectEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (p project) All(ctx context.Context, db storage.Executor) ([]ProjectEntity, error) {
	var entities []ProjectEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

func (p project) LatestPublished(
	ctx context.Context,
	db storage.Executor,
	limit int,
) ([]ProjectEntity, error) {
	var entities []ProjectEntity
	if err := db.NewSelect().
		Model(&entities).
		Where("published = TRUE").
		Where("started_at IS NOT NULL").
		OrderExpr("started_at DESC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type PaginatedProjects struct {
	Projects   []ProjectEntity
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func (p project) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedProjects, error) {
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

	totalCount, err := db.NewSelect().
		Model(&ProjectEntity{}).Count(ctx)
	if err != nil {
		return PaginatedProjects{}, err
	}

	entities := make([]ProjectEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedProjects{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedProjects{
		Projects:   entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (p project) PaginatePublished(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedProjects, error) {
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
		Model(&ProjectEntity{}).
		Where("published = TRUE").
		Where("started_at IS NOT NULL").
		Count(ctx)
	if err != nil {
		return PaginatedProjects{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}

	entities := make([]ProjectEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Where("published = TRUE").
		Where("started_at IS NOT NULL").
		OrderExpr("started_at DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Scan(ctx); err != nil {
		return PaginatedProjects{}, err
	}

	return PaginatedProjects{
		Projects:   entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (p project) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateProjectData,
) (ProjectEntity, error) {
	entity := ProjectEntity{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Published:       data.Published,
		Title:           data.Title,
		Slug:            data.Slug,
		StartedAt:       data.StartedAt,
		Status:          data.Status,
		Description:     data.Description,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		ImageLink:       normalizeOptionalString(data.ImageLink),
		MetaImageLink:   normalizeOptionalString(data.MetaImageLink),
		Content:         data.Content,
		ProjectUrl:      normalizeOptionalString(data.ProjectUrl),
	}

	if err := validation.Validate(&entity); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("published = excluded.published").
		Set("title = excluded.title").
		Set("slug = excluded.slug").
		Set("started_at = excluded.started_at").
		Set("status = excluded.status").
		Set("description = excluded.description").
		Set("meta_title = excluded.meta_title").
		Set("meta_description = excluded.meta_description").
		Set("image_link = excluded.image_link").
		Set("meta_image_link = excluded.meta_image_link").
		Set("content = excluded.content").
		Set("project_url = excluded.project_url").
		Returning("*").
		Scan(ctx); err != nil {
		return ProjectEntity{}, err
	}

	return entity, nil
}
