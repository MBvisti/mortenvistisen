package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models/internal/db"
)

type Tag struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
}

func FindTag(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) (Tag, error) {
	row, err := queries.QueryTagByID(ctx, exec, id)
	if err != nil {
		return Tag{}, err
	}

	return rowToTag(row), nil
}

type CreateTagData struct {
	Title string
}

func CreateTag(
	ctx context.Context,
	exec storage.Executor,
	data CreateTagData,
) (Tag, error) {
	if err := Validate.Struct(data); err != nil {
		return Tag{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.InsertTagParams{
		ID:    uuid.New(),
		Title: data.Title,
	}
	row, err := queries.InsertTag(ctx, exec, params)
	if err != nil {
		return Tag{}, err
	}

	return rowToTag(row), nil
}

type UpdateTagData struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	Title     string
}

func UpdateTag(
	ctx context.Context,
	exec storage.Executor,
	data UpdateTagData,
) (Tag, error) {
	if err := Validate.Struct(data); err != nil {
		return Tag{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpdateTagParams{
		ID:    data.ID,
		Title: data.Title,
	}

	row, err := queries.UpdateTag(ctx, exec, params)
	if err != nil {
		return Tag{}, err
	}

	return rowToTag(row), nil
}

func DestroyTag(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) error {
	return queries.DeleteTag(ctx, exec, id)
}

func AllTags(
	ctx context.Context,
	exec storage.Executor,
) ([]Tag, error) {
	rows, err := queries.QueryTags(ctx, exec)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, len(rows))
	for i, row := range rows {
		tags[i] = rowToTag(row)
	}

	return tags, nil
}

type PaginatedTags struct {
	Tags       []Tag
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func PaginateTags(
	ctx context.Context,
	exec storage.Executor,
	page int64,
	pageSize int64,
) (PaginatedTags, error) {
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

	totalCount, err := queries.CountTags(ctx, exec)
	if err != nil {
		return PaginatedTags{}, err
	}

	rows, err := queries.QueryPaginatedTags(
		ctx,
		exec,
		db.QueryPaginatedTagsParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedTags{}, err
	}

	tags := make([]Tag, len(rows))
	for i, row := range rows {
		tags[i] = rowToTag(row)
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedTags{
		Tags:       tags,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func UpsertTag(
	ctx context.Context,
	exec storage.Executor,
	data CreateTagData,
) (Tag, error) {
	if err := Validate.Struct(data); err != nil {
		return Tag{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpsertTagParams{
		ID:    uuid.New(),
		Title: data.Title,
	}
	row, err := queries.UpsertTag(ctx, exec, params)
	if err != nil {
		return Tag{}, err
	}

	return rowToTag(row), nil
}

func rowToTag(row db.Tag) Tag {
	return Tag{
		ID:        row.ID,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
		Title:     row.Title,
	}
}
