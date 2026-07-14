package factories

import (
	"context"
	"fmt"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"

	"github.com/go-faker/faker/v4"
)

// TagFactory wraps models.TagEntity for testing
type TagFactory struct {
	models.TagEntity
}

type TagOption func(*TagFactory)

// BuildTag creates an in-memory Tag with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateTag.
func BuildTag(opts ...TagOption) models.TagEntity {
	f := &TagFactory{
		TagEntity: models.TagEntity{
			Title: faker.Word(),
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.TagEntity
}

// CreateTag creates and persists a Tag to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateTag(
	ctx context.Context,
	exec storage.Executor,
	opts ...TagOption,
) (models.TagEntity, error) {
	built := BuildTag(opts...)

	entity := models.TagEntity{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     built.Title,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.TagEntity{}, err
	}

	return entity, nil
}

// CreateTags creates multiple Tag records at once
func CreateTags(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...TagOption,
) ([]models.TagEntity, error) {
	tags := make([]models.TagEntity, 0, count)

	for i := range count {
		entity, err := CreateTag(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create tag %d: %w", i+1, err)
		}
		tags = append(tags, entity)
	}

	return tags, nil
}

// Option functions

// WithTagsTitle sets the Title field
func WithTagsTitle(value string) TagOption {
	return func(f *TagFactory) {
		f.TagEntity.Title = value
	}
}
