package seeds

import (
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type Seeder struct {
	dbtx db.DBTX
}

func NewSeeder(dbtx db.DBTX) Seeder {
	return Seeder{dbtx}
}
