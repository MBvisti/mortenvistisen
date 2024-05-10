package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func ConvertToPGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func ConvertFromPGTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func ConvertFromPGTimestampToTime(t pgtype.Timestamp) time.Time {
	return t.Time
}
