package jobs

import "mortenvistisen/email"

type SendTransactionalEmailArgs struct {
	Data email.TransactionalData
}

func (SendTransactionalEmailArgs) Kind() string { return "send_transactional_email" }
