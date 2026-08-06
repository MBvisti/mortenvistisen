package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"mortenvistisen/config"
	"mortenvistisen/database"
	"mortenvistisen/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unexpected argument %q", args[0])
	}

	godotenv.Load()

	ctx := context.Background()

	db, err := storage.NewPostgres(ctx, buildDatabaseURL())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := storage.RunMigrations(ctx, db.Conn(), database.Migrations, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Migrations complete!")
	return nil
}

func buildDatabaseURL() string {
	return config.PostgresURL(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SSL_MODE"),
	)
}
