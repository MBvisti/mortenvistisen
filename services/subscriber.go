package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
)

var ErrSubscriberExists = errors.New("subscriber already exists")

func SubscribeToNewsletter(
	ctx context.Context,
	db psql.Postgres,
	email string,
	referer string,
) (models.Subscriber, models.Token, error) {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return models.Subscriber{}, models.Token{}, err
	}
	defer tx.Rollback(ctx)

	_, err = models.GetSubscriberByEmail(ctx, tx, email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.Subscriber{}, models.Token{}, ErrSubscriberExists
		}

		return models.Subscriber{}, models.Token{}, err
	}

	subscriber, err := models.NewSubscriber(
		ctx,
		tx,
		models.NewSubscriberPayload{
			Email:        email,
			SubscribedAt: time.Now(),
			Referer:      referer,
		},
	)
	if err != nil {
		return models.Subscriber{}, models.Token{}, err
	}

	token, err := models.NewCodeToken(ctx, tx, models.NewTokenPayload{
		Expiration: time.Now().Add(24 * time.Hour),
		Meta: models.MetaInformation{
			Resource:   models.ResourceSubscriber,
			ResourceID: subscriber.ID,
			Scope:      models.ScopeEmailVerification,
		},
	})
	if err != nil {
		return models.Subscriber{}, models.Token{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Subscriber{}, models.Token{}, err
	}

	return subscriber, token, nil
}
