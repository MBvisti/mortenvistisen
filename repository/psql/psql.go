package psql

import (
	"context"
	"errors"

	"github.com/MBvisti/mortenvistisen/repository/psql/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoRowWithIdentifier = errors.New("could not find requested row in database")

type Postgres struct {
	db *database.Queries
}

func NewPostgres(ctx context.Context, uri string) Postgres {
	dbpool, err := pgxpool.New(ctx, uri)
	if err != nil {
		panic("could not establish connection to database")
	}

	db := database.New(dbpool)

	return Postgres{
		db,
	}
}
