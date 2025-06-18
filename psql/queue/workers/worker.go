package workers

import (
	"github.com/mbvisti/mortenvistisen/clients"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/riverqueue/river"
)

type WorkerDependencies struct {
	DB          psql.Postgres
	EmailClient clients.Email
}

func SetupWorkers(deps WorkerDependencies) (*river.Workers, error) {
	workers := river.NewWorkers()

	if err := river.AddWorkerSafely(workers, &EmailJobWorker{
		emailClient: deps.EmailClient,
		db:          deps.DB,
	}); err != nil {
		return nil, err
	}

	return workers, nil
}
