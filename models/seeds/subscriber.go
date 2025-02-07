package seeds

import (
	"math/rand"
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	return func(sd *subscriberSeedData) {
		sd.ID = id
	}
}

func WithSubscriberCreatedAt(createdAt time.Time) subscriberSeedOption {
	return func(sd *subscriberSeedData) {
		sd.CreatedAt = createdAt
	}
}

func WithSubscriberUpdatedAt(updatedAt time.Time) subscriberSeedOption {
	return func(sd *subscriberSeedData) {
		sd.UpdatedAt = updatedAt
	}
}

func WithSubscriberEmail(email string) subscriberSeedOption {
	return func(sd *subscriberSeedData) {
		sd.Email = email
	}
}

func WithSubscriberSubscribedAt(subscribedAt time.Time) subscriberSeedOption {
	return func(sd *subscriberSeedData) {
		sd.SubscribedAt = subscribedAt
	}
}

func WithSubscriberReferer(referer string) subscriberSeedOption {
	return func(sd *subscriberSeedData) {
		sd.Referer = referer
	}
}

func WithSubscriberIsVerified(isVerified bool) subscriberSeedOption {
	return func(sd *subscriberSeedData) {
		sd.IsVerified = isVerified
	}
}

func (s Seeder) PlantSubscriber(
	ctx context.Context,
	opts ...subscriberSeedOption,
) (models.Subscriber, error) {
	trueOrFalse := rand.Float32() < 0.5

	data := &subscriberSeedData{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        faker.Email(),
		SubscribedAt: time.Now(),
		Referer:      "https://mortenvistisen.com",
		IsVerified:   trueOrFalse,
	}

	for _, opt := range opts {
		opt(data)
	}

	subscriber, err := models.NewSubscriber(
		ctx,
		s.dbtx,
		models.NewSubscriberPayload{
			Email:        data.Email,
			SubscribedAt: data.SubscribedAt,
			Referer:      data.Referer,
		},
	)
	if err != nil {
		return models.Subscriber{}, err
	}

	if data.IsVerified {
		if _, err := db.Stmts.UpdateSubscriberVerification(ctx, s.dbtx, db.UpdateSubscriberVerificationParams{
			ID: subscriber.ID,
			UpdatedAt: pgtype.Timestamptz{
				Time:  data.UpdatedAt,
				Valid: true,
			},
			IsVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		}); err != nil {
			return models.Subscriber{}, err
		}
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
