package seeds

import (
	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/jackc/pgx/v5"
)

//	type SeedBuilder[T, V any] interface {
//		WithRandoms(n int) *T
//		WithSpecific(data map[string]any) *T
//		Build() *V
//	}
//
//	type Seed interface {
//		Generate(ctx context.Context, dbtx db.DBTX) error
//	}
type Seeder struct {
	dbtx db.DBTX
}

func NewSeeder(tx pgx.Tx) Seeder {
	return Seeder{
		tx,
	}
}
