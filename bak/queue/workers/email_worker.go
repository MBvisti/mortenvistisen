package workers

import (
	"context"

	"github.com/MBvisti/mortenvistisen/queue/jobs"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/riverqueue/river"
)

type EmailJobWorker struct {
	emailer services.Mail
	river.WorkerDefaults[jobs.EmailJobArgs]
}

func (w *EmailJobWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.EmailJobArgs],
) error {
	return w.emailer.Send(
		services.MailPayload{
			To:       job.Args.To,
			From:     job.Args.From,
			Subject:  job.Args.Subject,
			HtmlBody: job.Args.HtmlVersion,
			TextBody: job.Args.TextVersion,
		},
	)
}
