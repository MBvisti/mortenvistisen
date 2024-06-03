package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/MBvisti/mortenvistisen/models/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type listOptions struct {
	limit  sql.NullInt32
	offset sql.NullInt32
	// orderBy string
}

type listOpt func(*listOptions)

func WithOffset(val int32) listOpt {
	return func(lso *listOptions) {
		lso.offset = sql.NullInt32{Int32: val, Valid: true}
	}
}

func WithLimit(val int32) listOpt {
	return func(lso *listOptions) {
		lso.limit = sql.NullInt32{Int32: val, Valid: true}
	}
}

func WithPagination(limit, offset int32) listOpt {
	return func(lso *listOptions) {
		lso.offset = sql.NullInt32{Int32: offset, Valid: true}
		lso.limit = sql.NullInt32{Int32: limit, Valid: true}
	}
}

type Models struct {
	Subscriber SubscriberModel
}

func NewModels(pool *pgxpool.Pool) Models {
	db := database.New(pool)
	return Models{
		Subscriber: SubscriberModel{db},
	}
}

func SetupDatabasePool(ctx context.Context, databaseURL string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
