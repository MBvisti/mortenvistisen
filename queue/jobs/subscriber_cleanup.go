package jobs

const subscriberCleanupJobKind string = "subscriber_cleanup_job"

type SubscriberCleanupJobArgs struct{}

func (SubscriberCleanupJobArgs) Kind() string { return subscriberCleanupJobKind }
