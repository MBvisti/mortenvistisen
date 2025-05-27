package main

import (
	"context"
	"log/slog"

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

	// // Create admin user
	// _, err = seeder.PlantUser(
	// 	ctx,
	// 	seeds.WithUserEmailVerifiedAt(time.Now()),
	// 	seeds.WithUserEmail("admin@mortenvistisen.com"),
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// slog.Info("Created admin user")
	//
	// Create sample articles for pagination testing
	articles, err := seeder.PlantArticles(ctx, 25)
	if err != nil {
		panic(err)
	}
	slog.Info("Created sample articles", "count", len(articles))

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
	slog.Info("Seed script finished")
}
