package psql

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/repository/psql/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SubscriberToken struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Hash         string
	ExpiresAt    time.Time
	Scope        string
	SubscriberID uuid.UUID
}

type Token struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Hash      string
	ExpiresAt time.Time
	Scope     string
	UserID    uuid.UUID
}

func (p Postgres) InsertSubscriberToken(
	ctx context.Context,
	hash, scope string, expiresAt time.Time,
	subscriberID uuid.UUID,
) error {
	err := p.db.InsertSubscriberToken(ctx, database.InsertSubscriberTokenParams{
		ID: uuid.New(),
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		Hash: hash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
		Scope:        scope,
		SubscriberID: subscriberID,
	})
	if err != nil {
		slog.Error("could not insert subscriber token", "error", err)
		return errors.Join(ErrInternalDBErr, err)
	}

	return nil
}
