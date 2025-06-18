package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/emails"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/psql/queue/jobs"
	"github.com/riverqueue/river"
)

var ErrSubscriberExists = errors.New("subscriber already exists")

func SubscribeToNewsletter(
	ctx context.Context,
	db psql.Postgres,
	queue *river.Client[pgx.Tx],
	email string,
	referer string,
) (models.Subscriber, models.Token, error) {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return models.Subscriber{}, models.Token{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx) // Rollback error is ignored as it's likely due to successful commit
	}()

	_, err = models.GetSubscriberByEmail(ctx, tx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

	html, text, err := emails.SubscriberWelcome{
		Email: subscriber.Email,
		Code:  token.Value,
	}.Generate(ctx)
	if err != nil {
		return models.Subscriber{}, models.Token{}, err
	}

	_, err = queue.InsertTx(ctx, tx, jobs.EmailJobArgs{
		Type:        "transaction",
		To:          subscriber.Email,
		From:        config.Cfg.DefaultSenderSignature,
		Subject:     "Welcome! Please verify your email",
		HtmlVersion: html.String(),
		TextVersion: text.String(),
	}, nil)
	if err != nil {
		return models.Subscriber{}, models.Token{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Subscriber{}, models.Token{}, err
	}

	return subscriber, token, nil
}

func VerifySubscriberEmail(
	ctx context.Context,
	db psql.Postgres,
	email string,
	code string,
) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx) // Rollback error is ignored as it's likely due to successful commit
	}()

	subscriber, err := models.GetSubscriberByEmail(ctx, tx, email)
	if err != nil {
		return err
	}

	if subscriber.IsVerified {
		return errors.New("subscriber already verified")
	}

	token, err := models.GetHashedToken(ctx, tx, code)
	if err != nil {
		return err
	}

	if !token.IsValid() ||
		token.Meta.Scope != models.ScopeEmailVerification ||
		token.Meta.Resource != models.ResourceSubscriber ||
		token.Meta.ResourceID != subscriber.ID {
		return errors.New("invalid verification code")
	}

	if err := models.VerifySubscriber(ctx, tx, models.VerifySubscriberPayload{
		ID:         subscriber.ID,
		UpdatedAt:  time.Now(),
		IsVerified: true,
	}); err != nil {
		return err
	}

	if err := models.DeleteToken(ctx, tx, token.ID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func GenerateUnsubscribeLink(
	ctx context.Context,
	db psql.Postgres,
	subscriberID uuid.UUID,
) (string, error) {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	expiration := time.Now().Add(30 * 24 * time.Hour)

	token, err := models.NewHashedToken(ctx, tx, models.NewTokenPayload{
		Expiration: expiration,
		Meta: models.MetaInformation{
			Resource:   models.ResourceSubscriber,
			ResourceID: subscriberID,
			Scope:      models.ScopeUnsubscribe,
		},
	})
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	unsubscribeURL := fmt.Sprintf("%s/unsubscribe/%s",
		config.Cfg.App.GetFullDomain(),
		token.Value)

	return unsubscribeURL, nil
}
