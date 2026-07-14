package queue

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"

	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/queue/jobs"
)

type DeleteArticleCoverWorker struct {
	river.WorkerDefaults[jobs.DeleteArticleCoverArgs]
	store objectstorage.CoverStore
}

func NewDeleteArticleCoverWorker(store objectstorage.CoverStore) *DeleteArticleCoverWorker {
	return &DeleteArticleCoverWorker{store: store}
}

func (w *DeleteArticleCoverWorker) Register(workers *river.Workers) error {
	return river.AddWorkerSafely(workers, w)
}

func (w *DeleteArticleCoverWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.DeleteArticleCoverArgs],
) error {
	if err := w.store.Delete(ctx, job.Args.Key); err != nil {
		return fmt.Errorf(
			"delete cover %q for article %d: %w",
			job.Args.Key,
			job.Args.ArticleID,
			err,
		)
	}

	return nil
}
