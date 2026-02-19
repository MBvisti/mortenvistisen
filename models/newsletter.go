package models

import (
	"context"
	"errors"
	"time"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models/internal/db"
)

var ErrNewsletterCannotBeUnpublished = errors.New("published newsletters cannot be unpublished")

type Newsletter struct {
	ID              int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Title           string
	Slug            string
	MetaTitle       string
	MetaDescription string
	IsPublished     bool
	ReleasedAt      time.Time
	Content         string
}

func FindNewsletter(
	ctx context.Context,
	exec storage.Executor,
	id int32,
) (Newsletter, error) {
	row, err := queries.QueryNewsletterByID(ctx, exec, id)
	if err != nil {
		return Newsletter{}, err
	}

	return rowToNewsletter(row), nil
}

func FindNewsletterBySlug(
	ctx context.Context,
	exec storage.Executor,
	slug string,
) (Newsletter, error) {
	rows, err := queries.QueryNewsletters(ctx, exec)
	if err != nil {
		return Newsletter{}, err
	}

	for _, row := range rows {
		if row.Slug.Valid && row.Slug.String == slug {
			return rowToNewsletter(row), nil
		}
	}

	return Newsletter{}, pgx.ErrNoRows
}

type CreateNewsletterData struct {
	Title           string
	Slug            string
	MetaTitle       string
	MetaDescription string
	IsPublished     bool
	ReleasedAt      time.Time
	Content         string
}

func CreateNewsletter(
	ctx context.Context,
	exec storage.Executor,
	data CreateNewsletterData,
) (Newsletter, error) {
	if err := Validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}
	params := db.InsertNewsletterParams{
		Title:           data.Title,
		Slug:            pgtype.Text{String: slug.Make(data.Title), Valid: data.Title != ""},
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		IsPublished:     pgtype.Bool{Bool: data.IsPublished, Valid: true},
		ReleasedAt:      pgtype.Timestamptz{Time: data.ReleasedAt, Valid: true},
		Content:         pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.InsertNewsletter(ctx, exec, params)
	if err != nil {
		return Newsletter{}, err
	}

	return rowToNewsletter(row), nil
}

type UpdateNewsletterData struct {
	ID              int32
	UpdatedAt       time.Time
	Title           string
	Slug            string
	MetaTitle       string
	MetaDescription string
	IsPublished     bool
	ReleasedAt      time.Time
	Content         string
}

func UpdateNewsletter(
	ctx context.Context,
	exec storage.Executor,
	data UpdateNewsletterData,
) (Newsletter, error) {
	if err := Validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}

	current, err := queries.QueryNewsletterByID(ctx, exec, data.ID)
	if err != nil {
		return Newsletter{}, err
	}

	if current.IsPublished.Bool && !data.IsPublished {
		return Newsletter{}, errors.Join(
			ErrDomainValidation,
			ErrNewsletterCannotBeUnpublished,
		)
	}

	params := db.UpdateNewsletterParams{
		ID:              data.ID,
		Title:           data.Title,
		Slug:            pgtype.Text{String: slug.Make(data.Title), Valid: data.Title != ""},
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		IsPublished:     pgtype.Bool{Bool: data.IsPublished, Valid: true},
		ReleasedAt:      pgtype.Timestamptz{Time: data.ReleasedAt, Valid: true},
		Content:         pgtype.Text{String: data.Content, Valid: true},
	}

	row, err := queries.UpdateNewsletter(ctx, exec, params)
	if err != nil {
		return Newsletter{}, err
	}

	return rowToNewsletter(row), nil
}

func DestroyNewsletter(
	ctx context.Context,
	exec storage.Executor,
	id int32,
) error {
	return queries.DeleteNewsletter(ctx, exec, id)
}

func AllNewsletters(
	ctx context.Context,
	exec storage.Executor,
) ([]Newsletter, error) {
	rows, err := queries.QueryNewsletters(ctx, exec)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = rowToNewsletter(row)
	}

	return newsletters, nil
}

func AllPublishedNewsletters(
	ctx context.Context,
	exec storage.Executor,
) ([]Newsletter, error) {
	rows, err := queries.QueryPublishedNewsletters(ctx, exec)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = rowToNewsletter(row)
	}

	return newsletters, nil
}

type PaginatedNewsletters struct {
	Newsletters []Newsletter
	TotalCount  int64
	Page        int64
	PageSize    int64
	TotalPages  int64
}

func PaginateNewsletters(
	ctx context.Context,
	exec storage.Executor,
	page int64,
	pageSize int64,
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

	offset := (page - 1) * pageSize

	totalCount, err := queries.CountNewsletters(ctx, exec)
	if err != nil {
		return PaginatedNewsletters{}, err
	}

	rows, err := queries.QueryPaginatedNewsletters(
		ctx,
		exec,
		db.QueryPaginatedNewslettersParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedNewsletters{}, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = rowToNewsletter(row)
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedNewsletters{
		Newsletters: newsletters,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

func UpsertNewsletter(
	ctx context.Context,
	exec storage.Executor,
	data CreateNewsletterData,
) (Newsletter, error) {
	if err := Validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}
	params := db.UpsertNewsletterParams{
		Title:           data.Title,
		Slug:            pgtype.Text{String: slug.Make(data.Title), Valid: data.Title != ""},
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		IsPublished:     pgtype.Bool{Bool: data.IsPublished, Valid: true},
		ReleasedAt:      pgtype.Timestamptz{Time: data.ReleasedAt, Valid: true},
		Content:         pgtype.Text{String: data.Content, Valid: true},
	}
	row, err := queries.UpsertNewsletter(ctx, exec, params)
	if err != nil {
		return Newsletter{}, err
	}

	return rowToNewsletter(row), nil
}

func CountNewsletters(
	ctx context.Context,
	exec storage.Executor,
) (int64, error) {
	return queries.CountNewsletters(ctx, exec)
}

func rowToNewsletter(row db.Newsletter) Newsletter {
	return Newsletter{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
		Title:           row.Title,
		Slug:            row.Slug.String,
		MetaTitle:       row.MetaTitle,
		MetaDescription: row.MetaDescription,
		IsPublished:     row.IsPublished.Bool,
		ReleasedAt:      row.ReleasedAt.Time,
		Content:         row.Content.String,
	}
}
