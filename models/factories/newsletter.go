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

// NewsletterFactory wraps models.Newsletter for testing
type NewsletterFactory struct {
	models.Newsletter // Embedded
}

type NewsletterOption func(*NewsletterFactory)

// BuildNewsletter creates an in-memory Newsletter with default test values
func BuildNewsletter(opts ...NewsletterOption) models.Newsletter {
	f := &NewsletterFactory{
		Newsletter: models.Newsletter{
			ID:              uuid.New(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Title:           faker.Word(),
			MetaTitle:       faker.Word(),
			MetaDescription: faker.Word(),
			IsPublished:     false,
			ReleasedAt:      time.Time{}, // Optional timestamp - zero by default
			Slug:            faker.Word(),
			Content:         faker.Word(),
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.Newsletter
}

// CreateNewsletter creates and persists a Newsletter to the database
func CreateNewsletter(
	ctx context.Context,
	exec storage.Executor,
	opts ...NewsletterOption,
) (models.Newsletter, error) {
	// Build with defaults and required FKs
	built := BuildNewsletter(opts...)

	// Prepare creation data
	data := models.CreateNewsletterData{
		Title:           built.Title,
		MetaTitle:       built.MetaTitle,
		MetaDescription: built.MetaDescription,
		IsPublished:     built.IsPublished,
		ReleasedAt:      built.ReleasedAt,
		Slug:            built.Slug,
		Content:         built.Content,
	}

	// Use model's Create function
	newsletter, err := models.CreateNewsletter(ctx, exec, data)
	if err != nil {
		return models.Newsletter{}, err
	}

	return newsletter, nil
}

// CreateNewsletters creates multiple Newsletter records at once
func CreateNewsletters(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...NewsletterOption,
) ([]models.Newsletter, error) {
	newsletters := make([]models.Newsletter, 0, count)

	for i := 0; i < count; i++ {
		newsletter, err := CreateNewsletter(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create newsletter %d: %w", i+1, err)
		}
		newsletters = append(newsletters, newsletter)
	}

	return newsletters, nil
}

// Option functions

// WithNewslettersTitle sets the Title field
func WithNewslettersTitle(value string) NewsletterOption {
	return func(f *NewsletterFactory) {
		f.Newsletter.Title = value
	}
}

// WithNewslettersMetaTitle sets the MetaTitle field
func WithNewslettersMetaTitle(value string) NewsletterOption {
	return func(f *NewsletterFactory) {
		f.Newsletter.MetaTitle = value
	}
}

// WithNewslettersMetaDescription sets the MetaDescription field
func WithNewslettersMetaDescription(value string) NewsletterOption {
	return func(f *NewsletterFactory) {
		f.Newsletter.MetaDescription = value
	}
}

// WithNewslettersIsPublished sets the IsPublished field
func WithNewslettersIsPublished(value bool) NewsletterOption {
	return func(f *NewsletterFactory) {
		f.Newsletter.IsPublished = value
	}
}

// WithNewslettersReleasedAt sets the ReleasedAt field
func WithNewslettersReleasedAt(value time.Time) NewsletterOption {
	return func(f *NewsletterFactory) {
		f.Newsletter.ReleasedAt = value
	}
}

// WithNewslettersSlug sets the Slug field
func WithNewslettersSlug(value string) NewsletterOption {
	return func(f *NewsletterFactory) {
		f.Newsletter.Slug = value
	}
}

// WithNewslettersContent sets the Content field
func WithNewslettersContent(value string) NewsletterOption {
	return func(f *NewsletterFactory) {
		f.Newsletter.Content = value
	}
}
