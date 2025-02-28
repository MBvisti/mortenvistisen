package jobs

import "github.com/google/uuid"

const scheduleNewsletterReleaseJobKind string = "ScheduleNewsletterRelease"

type ScheduleNewsletterRelease struct {
	NewsletterID uuid.UUID `json:"newsletter_id"`
}

func (ScheduleNewsletterRelease) Kind() string { return scheduleNewsletterReleaseJobKind }
