package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Subscriber struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	SubscribedAt time.Time
	Referer      string
	IsVerified   bool
	db           database.Queries
}

func NewSubscriber(db database.Queries) Subscriber {
	return Subscriber{
		db: db,
	}
}

func (s *Subscriber) ByID(ctx context.Context, id uuid.UUID) (Subscriber, error) {
	subscriber, err := s.db.QuerySubscriberByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Subscriber{}, ErrNoRowWithIdentifier
		}

		return Subscriber{}, err
	}

	return Subscriber{
		ID:           subscriber.ID,
		CreatedAt:    subscriber.CreatedAt.Time,
		UpdatedAt:    subscriber.UpdatedAt.Time,
		Email:        subscriber.Email.String,
		SubscribedAt: subscriber.SubscribedAt.Time,
		Referer:      subscriber.Referer.String,
		IsVerified:   subscriber.IsVerified.Bool,
	}, nil
}

func (s *Subscriber) ByEmail(ctx context.Context, email string) (Subscriber, error) {
	subscriber, err := s.db.QuerySubscriberByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Subscriber{}, ErrNoRowWithIdentifier
		}

		return Subscriber{}, err
	}

	return Subscriber{
		ID:           subscriber.ID,
		CreatedAt:    subscriber.CreatedAt.Time,
		UpdatedAt:    subscriber.UpdatedAt.Time,
		Email:        subscriber.Email.String,
		SubscribedAt: subscriber.SubscribedAt.Time,
		Referer:      subscriber.Referer.String,
		IsVerified:   subscriber.IsVerified.Bool,
	}, nil
}

func (s *Subscriber) List(ctx context.Context, opts ...listOpt) ([]Subscriber, error) {
	options := &listOptions{}

	for _, opt := range opts {
		opt(options)
	}

	subs, err := s.db.QuerySubscribers(ctx, database.QuerySubscribersParams{
		Offset: options.offset,
		Limit:  options.limit,
	})
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(subs))
	for i, sub := range subs {
		subscribers[i] = Subscriber{
			ID:           sub.ID,
			CreatedAt:    sub.CreatedAt.Time,
			UpdatedAt:    sub.UpdatedAt.Time,
			Email:        sub.Email.String,
			SubscribedAt: sub.SubscribedAt.Time,
			Referer:      sub.Referer.String,
			IsVerified:   sub.IsVerified.Bool,
		}
	}

	return subscribers, nil
}

func (s *Subscriber) Count(ctx context.Context) (int64, error) {
	count, err := s.db.QuerySubscriberCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Subscriber) NewForCurrentMonth(ctx context.Context) (int64, error) {
	count, err := s.db.QueryNewSubscribersForCurrentMonth(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Subscriber) UnverifiedCount(ctx context.Context) (int64, error) {
	count, err := s.db.QueryUnverifiedSubCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
