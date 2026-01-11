package factories

import (
	"context"
	"fmt"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

// TagFactory wraps models.Tag for testing
type TagFactory struct {
	models.Tag // Embedded
}

type TagOption func(*TagFactory)

// BuildTag creates an in-memory Tag with default test values
func BuildTag(opts ...TagOption) models.Tag {
	f := &TagFactory{
		Tag: models.Tag{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     faker.Word(),
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.Tag
}

// CreateTag creates and persists a Tag to the database
func CreateTag(ctx context.Context, exec storage.Executor, opts ...TagOption) (models.Tag, error) {
	// Build with defaults and required FKs
	built := BuildTag(opts...)

	// Prepare creation data
	data := models.CreateTagData{
		Title: built.Title,
	}

	// Use model's Create function
	tag, err := models.CreateTag(ctx, exec, data)
	if err != nil {
		return models.Tag{}, err
	}

	return tag, nil
}

// CreateTags creates multiple Tag records at once
func CreateTags(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...TagOption,
) ([]models.Tag, error) {
	tags := make([]models.Tag, 0, count)

	for i := 0; i < count; i++ {
		tag, err := CreateTag(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create tag %d: %w", i+1, err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// Option functions

// WithTagsTitle sets the Title field
func WithTagsTitle(value string) TagOption {
	return func(f *TagFactory) {
		f.Tag.Title = value
	}
}
