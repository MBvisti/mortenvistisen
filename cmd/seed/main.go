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

	_, err = seeder.PlantSubscriber(
		ctx,
		seeds.WithSubscriberEmail("hello@mbvlabs.com"),
		seeds.WithSubscriberIsVerified(true),
	)
	if err != nil {
		panic(err)
	}
	_, err = seeder.PlantSubscriber(
		ctx,
		seeds.WithSubscriberEmail("hello1@mbvlabs.com"),
		seeds.WithSubscriberIsVerified(true),
	)
	if err != nil {
		panic(err)
	}
	_, err = seeder.PlantSubscriber(
		ctx,
		seeds.WithSubscriberEmail("hello2@mbvlabs.com"),
		seeds.WithSubscriberIsVerified(true),
	)
	if err != nil {
		panic(err)
	}
	_, err = seeder.PlantSubscriber(
		ctx,
		seeds.WithSubscriberEmail("hello3@mbvlabs.com"),
		seeds.WithSubscriberIsVerified(true),
	)
	if err != nil {
		panic(err)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(err)
	}
}
