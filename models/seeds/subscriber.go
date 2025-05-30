package seeds

import (
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
	"golang.org/x/net/context"
)

type subscriberSeedData struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	SubscribedAt time.Time
	Referer      string
	IsVerified   bool
}

type subscriberSeedOption func(*subscriberSeedData)

func WithSubscriberID(id uuid.UUID) subscriberSeedOption {
	return func(ssd *subscriberSeedData) {
		ssd.ID = id
	}
}

func WithSubscriberCreatedAt(createdAt time.Time) subscriberSeedOption {
	return func(ssd *subscriberSeedData) {
		ssd.CreatedAt = createdAt
	}
}

func WithSubscriberUpdatedAt(updatedAt time.Time) subscriberSeedOption {
	return func(ssd *subscriberSeedData) {
		ssd.UpdatedAt = updatedAt
	}
}

func WithSubscriberEmail(email string) subscriberSeedOption {
	return func(ssd *subscriberSeedData) {
		ssd.Email = email
	}
}

func WithSubscriberSubscribedAt(subscribedAt time.Time) subscriberSeedOption {
	return func(ssd *subscriberSeedData) {
		ssd.SubscribedAt = subscribedAt
	}
}

func WithSubscriberReferer(referer string) subscriberSeedOption {
	return func(ssd *subscriberSeedData) {
		ssd.Referer = referer
	}
}

func WithSubscriberIsVerified(isVerified bool) subscriberSeedOption {
	return func(ssd *subscriberSeedData) {
		ssd.IsVerified = isVerified
	}
}

func (s Seeder) PlantSubscriber(
	ctx context.Context,
	opts ...subscriberSeedOption,
) (models.Subscriber, error) {
	now := time.Now()
	data := &subscriberSeedData{
		ID:           uuid.New(),
		CreatedAt:    now,
		UpdatedAt:    now,
		Email:        faker.Email(),
		SubscribedAt: now,
		Referer:      faker.URL(),
		IsVerified:   false,
	}

	for _, opt := range opts {
		opt(data)
	}

	subscriber, err := models.NewSubscriber(ctx, s.dbtx, models.NewSubscriberPayload{
		Email:        data.Email,
		SubscribedAt: data.SubscribedAt,
		Referer:      data.Referer,
	})
	if err != nil {
		return models.Subscriber{}, err
	}

	if data.IsVerified {
		if err := db.Stmts.VerifySubscriber(ctx, s.dbtx, db.VerifySubscriberParams{
			ID: subscriber.ID,
			UpdatedAt: pgtype.Timestamptz{
				Time:  data.UpdatedAt,
				Valid: true,
			},
			IsVerified: pgtype.Bool{Bool: data.IsVerified, Valid: true},
		}); err != nil {
			return models.Subscriber{}, err
		}
		subscriber.IsVerified = data.IsVerified
	}

	return subscriber, nil
}

func (s Seeder) PlantSubscribers(
	ctx context.Context,
	amount int,
) ([]models.Subscriber, error) {
	subscribers := make([]models.Subscriber, amount)

	for i := range amount {
		subscriber, err := s.PlantSubscriber(ctx)
		if err != nil {
			return nil, err
		}

		subscribers[i] = subscriber
	}

	return subscribers, nil
}

func (s Seeder) PlantVerifiedSubscribers(
	ctx context.Context,
	amount int,
) ([]models.Subscriber, error) {
	subscribers := make([]models.Subscriber, amount)

	for i := range amount {
		subscriber, err := s.PlantSubscriber(ctx, WithSubscriberIsVerified(true))
		if err != nil {
			return nil, err
		}

		subscribers[i] = subscriber
	}

	return subscribers, nil
}
