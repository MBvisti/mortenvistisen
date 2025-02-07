package seeds

import (
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type newsletterSeedData struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string
	Content    string
	ReleasedAt time.Time
	Released   bool
}

type newsletterSeedOption func(*newsletterSeedData)

func WithNewsletterID(id uuid.UUID) newsletterSeedOption {
	return func(nsd *newsletterSeedData) {
		nsd.ID = id
	}
}

func WithNewsletterCreatedAt(createdAt time.Time) newsletterSeedOption {
	return func(nsd *newsletterSeedData) {
		nsd.CreatedAt = createdAt
	}
}

func WithNewsletterUpdatedAt(updatedAt time.Time) newsletterSeedOption {
	return func(nsd *newsletterSeedData) {
		nsd.UpdatedAt = updatedAt
	}
}

func WithNewsletterTitle(title string) newsletterSeedOption {
	return func(nsd *newsletterSeedData) {
		nsd.Title = title
	}
}

func WithNewsletterContent(content string) newsletterSeedOption {
	return func(nsd *newsletterSeedData) {
		nsd.Content = content
	}
}

func WithNewsletterReleasedAt(releasedAt time.Time) newsletterSeedOption {
	return func(nsd *newsletterSeedData) {
		nsd.ReleasedAt = releasedAt
	}
}

func WithNewsletterReleased(released bool) newsletterSeedOption {
	return func(nsd *newsletterSeedData) {
		nsd.Released = released
	}
}

func (s Seeder) PlantNewsletter(
	ctx context.Context,
	opts ...newsletterSeedOption,
) (models.Newsletter, error) {
	data := &newsletterSeedData{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Title:      faker.Sentence(),
		Content:    faker.Paragraph(),
		ReleasedAt: time.Now().Add(24 * time.Hour), // Default to tomorrow
		Released:   false,
	}

	for _, opt := range opts {
		opt(data)
	}

	newsletter, err := models.NewNewsletter(
		ctx,
		s.dbtx,
		models.NewNewsletterPayload{
			Title:      data.Title,
			Content:    data.Content,
			ReleasedAt: data.ReleasedAt,
			Released:   data.Released,
		},
	)
	if err != nil {
		return models.Newsletter{}, err
	}

	return newsletter, nil
}

func (s Seeder) PlantNewsletters(
	ctx context.Context,
	amount int,
) ([]models.Newsletter, error) {
	newsletters := make([]models.Newsletter, amount)

	for i := range amount {
		newsletter, err := s.PlantNewsletter(ctx)
		if err != nil {
			return nil, err
		}

		newsletters[i] = newsletter
	}

	return newsletters, nil
}
