package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/MBvisti/mortenvistisen/migrations"
	"github.com/MBvisti/mortenvistisen/pkg/config"
)

func main() {
	cfg := config.New()

	conn, err := sql.Open("postgres", cfg.Db.GetUrlString())
	if err != nil {
		panic(err)
	}

	if err := migrations.Up(conn); err != nil {
		panic(err)
	}
}
