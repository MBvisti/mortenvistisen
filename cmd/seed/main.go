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

	// Create admin user
	_, err = seeder.PlantUser(
		ctx,
		seeds.WithUserEmailVerifiedAt(time.Now()),
		seeds.WithUserEmail("admin@mortenvistisen.com"),
	)
	if err != nil {
		panic(err)
	}
	slog.Info("Created admin user")

	// Create predefined article tags
	tags, err := seeder.PlantPredefinedArticleTags(ctx)
	if err != nil {
		panic(err)
	}
	slog.Info("Created article tags", "count", len(tags))

	// Create sample articles with random tags assigned (1-5 tags per article)
	articles, err := seeder.PlantArticlesWithRandomTags(ctx, 25, tags, 1, 5)
	if err != nil {
		panic(err)
	}
	slog.Info("Created sample articles with tags", "count", len(articles))

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
	slog.Info("Seed script finished")
}
