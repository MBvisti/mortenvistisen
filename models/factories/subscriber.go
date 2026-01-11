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

// SubscriberFactory wraps models.Subscriber for testing
type SubscriberFactory struct {
	models.Subscriber // Embedded
}

type SubscriberOption func(*SubscriberFactory)

// BuildSubscriber creates an in-memory Subscriber with default test values
func BuildSubscriber(opts ...SubscriberOption) models.Subscriber {
	f := &SubscriberFactory{
		Subscriber: models.Subscriber{
			ID:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Email:        faker.Word(),
			SubscribedAt: time.Time{}, // Optional timestamp - zero by default
			Referer:      faker.Word(),
			IsVerified:   false,
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.Subscriber
}

// CreateSubscriber creates and persists a Subscriber to the database
func CreateSubscriber(
	ctx context.Context,
	exec storage.Executor,
	opts ...SubscriberOption,
) (models.Subscriber, error) {
	// Build with defaults and required FKs
	built := BuildSubscriber(opts...)

	// Prepare creation data
	data := models.CreateSubscriberData{
		Email:        built.Email,
		SubscribedAt: built.SubscribedAt,
		Referer:      built.Referer,
		IsVerified:   built.IsVerified,
	}

	// Use model's Create function
	subscriber, err := models.CreateSubscriber(ctx, exec, data)
	if err != nil {
		return models.Subscriber{}, err
	}

	return subscriber, nil
}

// CreateSubscribers creates multiple Subscriber records at once
func CreateSubscribers(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...SubscriberOption,
) ([]models.Subscriber, error) {
	subscribers := make([]models.Subscriber, 0, count)

	for i := 0; i < count; i++ {
		subscriber, err := CreateSubscriber(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create subscriber %d: %w", i+1, err)
		}
		subscribers = append(subscribers, subscriber)
	}

	return subscribers, nil
}

// Option functions

// WithSubscribersEmail sets the Email field
func WithSubscribersEmail(value string) SubscriberOption {
	return func(f *SubscriberFactory) {
		f.Subscriber.Email = value
	}
}

// WithSubscribersSubscribedAt sets the SubscribedAt field
func WithSubscribersSubscribedAt(value time.Time) SubscriberOption {
	return func(f *SubscriberFactory) {
		f.Subscriber.SubscribedAt = value
	}
}

// WithSubscribersReferer sets the Referer field
func WithSubscribersReferer(value string) SubscriberOption {
	return func(f *SubscriberFactory) {
		f.Subscriber.Referer = value
	}
}

// WithSubscribersIsVerified sets the IsVerified field
func WithSubscribersIsVerified(value bool) SubscriberOption {
	return func(f *SubscriberFactory) {
		f.Subscriber.IsVerified = value
	}
}
