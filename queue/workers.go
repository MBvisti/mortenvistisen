// Package queue provides queue-specific resources.
package queue

import (
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

var wrksConstructors = fx.Provide(
	river.NewWorkers,
	NewSendTransactionalEmailWorker,
	NewSendMarketingEmailWorker,
	NewDeleteArticleCoverWorker,
	NewDeleteNewsletterCoverWorker,
	NewDeleteProjectCoverWorker,
)

var WorkersModule = fx.Module(
	"queue-workers",
	wrksConstructors,
	fx.Invoke(func(workers *river.Workers, worker *SendTransactionalEmailWorker) error {
		return worker.Register(workers)
	}),
	fx.Invoke(func(workers *river.Workers, worker *SendMarketingEmailWorker) error {
		return worker.Register(workers)
	}),
	fx.Invoke(func(workers *river.Workers, worker *DeleteArticleCoverWorker) error {
		return worker.Register(workers)
	}),
	fx.Invoke(func(workers *river.Workers, worker *DeleteNewsletterCoverWorker) error {
		return worker.Register(workers)
	}),
	fx.Invoke(func(workers *river.Workers, worker *DeleteProjectCoverWorker) error {
		return worker.Register(workers)
	}),
)
