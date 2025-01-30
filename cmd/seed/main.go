package main

import (
	"context"

	"github.com/MBvisti/mortenvistisen/config"
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
	_, err = pool.Begin(ctx)
	if err != nil {
		panic(err)
	}
}
