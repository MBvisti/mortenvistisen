package psql

import (
	"context"
	"time"

	"github.com/MBvisti/mortenvistisen/repository/psql/database"
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
	err := p.Queries.InsertSubscriberToken(ctx, database.InsertSubscriberTokenParams{
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
		return err
	}

	return nil
}

func (p Postgres) InsertToken(
	ctx context.Context,
	hash, scope string, expiresAt time.Time,
	userID uuid.UUID,
) error {
	err := p.Queries.InsertToken(ctx, database.InsertTokenParams{
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
		Scope:  scope,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p Postgres) DeleteTokenByHash(ctx context.Context, hash string) error {
	return p.Queries.DeleteTokenByHash(ctx, hash)
}

func (p Postgres) DeleteTokenBySubID(ctx context.Context, id uuid.UUID) error {
	return p.Queries.DeleteSubscriberTokenBySubscriberID(ctx, id)
}

func (p Postgres) QueryTokenByHash(ctx context.Context, hash string) (database.Token, error) {
	return p.Queries.QueryTokenByHash(ctx, hash)
}

func (p Postgres) QuerySubscriberTokenByHash(
	ctx context.Context,
	hash string,
) (database.SubscriberToken, error) {
	return p.Queries.QuerySubscriberTokenByHash(ctx, hash)
}
