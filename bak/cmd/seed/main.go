package main

import (
	"context"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/models/seeds"
	"github.com/MBvisti/mortenvistisen/psql"
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

	seeder := seeds.NewSeeder(tx)
	if _, err := seeder.PlantUser(
		ctx,
		seeds.WithUserEmail("admin@mbvlabs.com"),
		seeds.WithUserIsAdmin(true),
	); err != nil {
		panic(err)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
}
