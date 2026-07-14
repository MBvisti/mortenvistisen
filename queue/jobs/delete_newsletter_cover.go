package jobs

type DeleteNewsletterCoverArgs struct {
	NewsletterID int32  `json:"newsletter_id"`
	Key          string `json:"key"`
}

func (DeleteNewsletterCoverArgs) Kind() string { return "delete_newsletter_cover" }
