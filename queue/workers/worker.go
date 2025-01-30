package workers

import (
	"github.com/riverqueue/river"
)

type WorkerDependencies struct {
	// DB      *database.Queries
	// Emailer emails.EmailClient
	// Tracer  telemetry.Tracer
}

func SetupWorkers(deps WorkerDependencies) (*river.Workers, error) {
	workers := river.NewWorkers()

	// if err := river.AddWorkerSafely(workers, &EmailJobWorker{
	// 	emailer: deps.Emailer,
	// }); err != nil {
	// 	return nil, err
	// }

	return workers, nil
}
