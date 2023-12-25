package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func SetupDatabaseConnection(ctx context.Context, databaseURL string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
