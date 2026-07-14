package queue

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"

	objectstorage "mortenvistisen/clients/storage"
	"mortenvistisen/queue/jobs"
)

type DeleteNewsletterCoverWorker struct {
	river.WorkerDefaults[jobs.DeleteNewsletterCoverArgs]
	store objectstorage.CoverStore
}

func NewDeleteNewsletterCoverWorker(store objectstorage.CoverStore) *DeleteNewsletterCoverWorker {
	return &DeleteNewsletterCoverWorker{store: store}
}

func (w *DeleteNewsletterCoverWorker) Register(workers *river.Workers) error {
	return river.AddWorkerSafely(workers, w)
}

func (w *DeleteNewsletterCoverWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.DeleteNewsletterCoverArgs],
) error {
	if err := w.store.Delete(ctx, job.Args.Key); err != nil {
		return fmt.Errorf(
			"delete cover %q for newsletter %d: %w",
			job.Args.Key,
			job.Args.NewsletterID,
			err,
		)
	}
	return nil
}
