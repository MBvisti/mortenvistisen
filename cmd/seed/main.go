package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/models/seeds"
	"github.com/mbvisti/mortenvistisen/psql"
)

func main() {
	pool, err := psql.CreatePooledConnection(
		context.Background(),
		config.Cfg.GetDatabaseURL(),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		panic(err)
	}

	///nolint:errcheck
	defer tx.Rollback(ctx)

	slog.Info("Starting seed script...")

	seeder := seeds.NewSeeder(pool)
	_, err = seeder.PlantUser(
		ctx,
		seeds.WithUserEmailVerifiedAt(time.Now()),
		seeds.WithUserEmail("aryastark@gmail.com"),
	)
	if err != nil {
		panic(err)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
	slog.Info("Seed script finished")
}
