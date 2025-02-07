package seeds

import (
	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/jackc/pgx/v5"
)

type Seeder struct {
	dbtx db.DBTX
}

func NewSeeder(tx pgx.Tx) Seeder {
	return Seeder{
		tx,
	}
}
