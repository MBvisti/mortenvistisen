package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models/internal/db"
)

type Project struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Published   bool
	Title       string
	Slug        string
	StartedAt   time.Time
	Status      string
	Description string
	Content     string
	ProjectURL  string
}

func FindProject(
	ctx context.Context,
	exec storage.Executor,
	id int32,
) (Project, error) {
	row, err := queries.QueryProjectByID(ctx, exec, id)
	if err != nil {
		return Project{}, err
	}

	return rowToProject(row), nil
}

func FindProjectBySlug(
	ctx context.Context,
	exec storage.Executor,
	slug string,
) (Project, error) {
	row, err := queries.QueryProjectBySlug(ctx, exec, slug)
	if err != nil {
		return Project{}, err
	}

	return rowToProject(row), nil
}

type CreateProjectData struct {
	Published   bool
	Title       string
	Slug        string
	StartedAt   time.Time
	Status      string
	Description string
	Content     string
	ProjectURL  string
}

func CreateProject(
	ctx context.Context,
	exec storage.Executor,
	data CreateProjectData,
) (Project, error) {
	if err := Validate.Struct(data); err != nil {
		return Project{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.InsertProjectParams{
		Published:   data.Published,
		Title:       data.Title,
		Slug:        data.Slug,
		StartedAt:   pgtype.Timestamptz{Time: data.StartedAt, Valid: !data.StartedAt.IsZero()},
		Status:      data.Status,
		Description: data.Description,
		Content:     data.Content,
		ProjectUrl:  pgtype.Text{String: data.ProjectURL, Valid: data.ProjectURL != ""},
	}

	row, err := queries.InsertProject(ctx, exec, params)
	if err != nil {
		return Project{}, err
	}

	return rowToProject(row), nil
}

type UpdateProjectData struct {
	ID          int32
	Published   bool
	Title       string
	Slug        string
	StartedAt   time.Time
	Status      string
	Description string
	Content     string
	ProjectURL  string
}

func UpdateProject(
	ctx context.Context,
	exec storage.Executor,
	data UpdateProjectData,
) (Project, error) {
	if err := Validate.Struct(data); err != nil {
		return Project{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpdateProjectParams{
		ID:          data.ID,
		Published:   data.Published,
		Title:       data.Title,
		Slug:        data.Slug,
		StartedAt:   pgtype.Timestamptz{Time: data.StartedAt, Valid: !data.StartedAt.IsZero()},
		Status:      data.Status,
		Description: data.Description,
		Content:     data.Content,
		ProjectUrl:  pgtype.Text{String: data.ProjectURL, Valid: data.ProjectURL != ""},
	}

	row, err := queries.UpdateProject(ctx, exec, params)
	if err != nil {
		return Project{}, err
	}

	return rowToProject(row), nil
}

func DestroyProject(
	ctx context.Context,
	exec storage.Executor,
	id int32,
) error {
	return queries.DeleteProject(ctx, exec, id)
}

func AllProjects(
	ctx context.Context,
	exec storage.Executor,
) ([]Project, error) {
	rows, err := queries.QueryProjects(ctx, exec)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, len(rows))
	for i, row := range rows {
		projects[i] = rowToProject(row)
	}

	return projects, nil
}

func AllPublishedProjects(
	ctx context.Context,
	exec storage.Executor,
) ([]Project, error) {
	rows, err := queries.QueryPublishedProjects(ctx, exec)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, len(rows))
	for i, row := range rows {
		projects[i] = rowToProject(row)
	}

	return projects, nil
}

type PaginatedProjects struct {
	Projects   []Project
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func PaginateProjects(
	ctx context.Context,
	exec storage.Executor,
	page int64,
	pageSize int64,
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

	totalCount, err := queries.CountProjects(ctx, exec)
	if err != nil {
		return PaginatedProjects{}, err
	}

	rows, err := queries.QueryPaginatedProjects(
		ctx,
		exec,
		db.QueryPaginatedProjectsParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedProjects{}, err
	}

	projects := make([]Project, len(rows))
	for i, row := range rows {
		projects[i] = rowToProject(row)
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedProjects{
		Projects:   projects,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func UpsertProject(
	ctx context.Context,
	exec storage.Executor,
	data CreateProjectData,
) (Project, error) {
	if err := Validate.Struct(data); err != nil {
		return Project{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpsertProjectParams{
		Published:   data.Published,
		Title:       data.Title,
		Slug:        data.Slug,
		StartedAt:   pgtype.Timestamptz{Time: data.StartedAt, Valid: !data.StartedAt.IsZero()},
		Status:      data.Status,
		Description: data.Description,
		Content:     data.Content,
		ProjectUrl:  pgtype.Text{String: data.ProjectURL, Valid: data.ProjectURL != ""},
	}

	row, err := queries.UpsertProject(ctx, exec, params)
	if err != nil {
		return Project{}, err
	}

	return rowToProject(row), nil
}

func CountProjects(
	ctx context.Context,
	exec storage.Executor,
) (int64, error) {
	return queries.CountProjects(ctx, exec)
}

func rowToProject(row db.Project) Project {
	return Project{
		ID:          row.ID,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
		Published:   row.Published,
		Title:       row.Title,
		Slug:        row.Slug,
		StartedAt:   row.StartedAt.Time,
		Status:      row.Status,
		Description: row.Description,
		Content:     row.Content,
		ProjectURL:  row.ProjectUrl.String,
	}
}
