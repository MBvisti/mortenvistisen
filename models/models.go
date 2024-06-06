package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
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
	Newsletter NewsletterModel
	Article    ArticleModel
}

func NewModels(pool *pgxpool.Pool, v *validator.Validate) Models {
	db := database.New(pool)
	return Models{
		SubscriberModel{db},
		NewsletterModel{db, v},
		ArticleModel{db},
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

func ConvertToPGTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
