package queue

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"

	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/queue/jobs"
)

type DeleteProjectCoverWorker struct {
	river.WorkerDefaults[jobs.DeleteProjectCoverArgs]
	store objectstorage.CoverStore
}

func NewDeleteProjectCoverWorker(store objectstorage.CoverStore) *DeleteProjectCoverWorker {
	return &DeleteProjectCoverWorker{store: store}
}

func (w *DeleteProjectCoverWorker) Register(workers *river.Workers) error {
	return river.AddWorkerSafely(workers, w)
}

func (w *DeleteProjectCoverWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.DeleteProjectCoverArgs],
) error {
	if err := w.store.Delete(ctx, job.Args.Key); err != nil {
		return fmt.Errorf(
			"delete image %q for project %d: %w",
			job.Args.Key,
			job.Args.ProjectID,
			err,
		)
	}
	return nil
}
