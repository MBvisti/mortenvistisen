package jobs

const marketingEmailJobKind string = "marketing_email_job"

type MarketingEmailJobArgs struct {
	To              string `json:"to"`
	From            string `json:"from"`
	Subject         string `json:"subject"`
	TextVersion     string `json:"text_version"`
	HtmlVersion     string `json:"html_version"`
	SubscriberID    string `json:"subscriber_id"`
	NewsletterID    string `json:"newsletter_id"`
	UnsubscribeLink string `json:"unsubscribe_link"`
}

func (MarketingEmailJobArgs) Kind() string { return marketingEmailJobKind }
