package psql

import (
	"context"
	"errors"
	"log/slog"
	"time"

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
}

func (p *Postgres) GetSubscriberByID(ctx context.Context, id uuid.UUID) (Subscriber, error) {
	subscriber, err := p.db.QuerySubscriberByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Subscriber{}, ErrNoRowWithIdentifier
		}

		slog.Error("could not get subscriber by the provided id", "id", id, "error", err)
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
