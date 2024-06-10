package psql

import (
	"context"
	"errors"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/repository/psql/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInternalDBErr       = errors.New("an error occurred that was not possible to recover from")
	ErrNoRowWithIdentifier = errors.New("could not find requested row in database")
)

type Postgres struct {
	db *database.Queries
}

func NewPostgres(dbPool *pgxpool.Pool) Postgres {
	return Postgres{
		db: database.New(dbPool),
	}
}

func CreatePooledConnection(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(ctx, uri)
	if err != nil {
		slog.Error("could not establish connection to database", "error", err)
		return nil, err
	}

	return dbpool, nil
}
