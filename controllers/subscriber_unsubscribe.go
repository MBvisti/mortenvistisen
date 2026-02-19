package controllers

import (
	"context"
	"encoding/json"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
)

const subscriberUnsubscribeScope = "subscriber_newsletter_unsubscribe"

type subscriberUnsubscribeMeta struct {
	SubscriberID int32 `json:"subscriber_id"`
}

func createSubscriberUnsubscribeToken(
	ctx context.Context,
	exec storage.Executor,
	pepper string,
	subscriberID int32,
) (string, error) {
	meta, err := json.Marshal(subscriberUnsubscribeMeta{
		SubscriberID: subscriberID,
	})
	if err != nil {
		return "", err
	}

	// Long-lived unsubscribe links are expected to keep working from old emails.
	expiresAt := time.Now().AddDate(2, 0, 0)

	return models.CreateToken(
		ctx,
		exec,
		pepper,
		subscriberUnsubscribeScope,
		expiresAt,
		meta,
	)
}
