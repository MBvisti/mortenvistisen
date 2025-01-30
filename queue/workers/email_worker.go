package workers

// type EmailJobWorker struct {
// 	emailer services.Mail
// 	river.WorkerDefaults[jobs.EmailJobArgs]
// }
//
// func (w *EmailJobWorker) Work(
// 	ctx context.Context,
// 	job *river.Job[jobs.EmailJobArgs],
// ) error {
// 	return w.emailer.SendNewSubscriber(
// 		ctx,
// 		job.Args.To,
// 		job.Args.From,
// 		job.Args.Subject,
// 	)
// }
