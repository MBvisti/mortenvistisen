package jobs

const subscriberCleanupJobKind string = "SubscriberCleanupJob"

type SubscriberCleanupJobArgs struct{}

func (SubscriberCleanupJobArgs) Kind() string { return subscriberCleanupJobKind }
