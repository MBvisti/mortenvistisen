package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/email"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/queue/jobs"
	"mortenvistisen/router/routes"
)

const subscriberEmailVerification = "subscriber_email_verification"

var (
	ErrSubscriberVerificationInvalidCode = errors.New("invalid subscriber verification code")
	ErrSubscriberVerificationExpiredCode = errors.New("subscriber verification code has expired")
	ErrSubscriberAlreadyVerified         = errors.New("subscriber already verified")
)

type RequestSubscriberVerificationData struct {
	Email   string
	Referer string
}

func RequestSubscriberVerification(
	ctx context.Context,
	db storage.Pool,
	insertOnly queue.InsertOnly,
	pepper string,
	data RequestSubscriberVerificationData,
) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	emailAddress := strings.ToLower(strings.TrimSpace(data.Email))
	if emailAddress == "" {
		return errors.New("email is required")
	}

	subscriber, err := models.FindSubscriberByEmail(ctx, tx, emailAddress)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("find subscriber by email: %w", err)
		}

		subscriber, err = models.CreateSubscriber(ctx, tx, models.CreateSubscriberData{
			Email:        emailAddress,
			SubscribedAt: time.Now().UTC(),
			Referer:      data.Referer,
			IsVerified:   false,
		})
		if err != nil {
			return fmt.Errorf("create subscriber: %w", err)
		}
	}

	if subscriber.IsVerified {
		return ErrSubscriberAlreadyVerified
	}

	meta, err := json.Marshal(map[string]string{
		"subscriber_id": strconv.Itoa(int(subscriber.ID)),
	})
	if err != nil {
		return fmt.Errorf("marshal subscriber token metadata: %w", err)
	}

	code, err := models.CreateCodeToken(
		ctx,
		tx,
		pepper,
		subscriberEmailVerification,
		time.Now().Add(30*time.Minute),
		meta,
	)
	if err != nil {
		return fmt.Errorf("create subscriber verification token: %w", err)
	}

	verifyURL := fmt.Sprintf("%s%s", config.BaseURL, routes.SubscriberVerificationNew.URL())
	verifyEmail := email.VerifyNewsletterSubscription{
		VerificationCode: code,
		VerificationURL:  verifyURL,
	}

	html, err := verifyEmail.ToHTML()
	if err != nil {
		return fmt.Errorf("render subscriber verification html email: %w", err)
	}

	text, err := verifyEmail.ToText()
	if err != nil {
		return fmt.Errorf("render subscriber verification text email: %w", err)
	}

	_, err = insertOnly.InsertTx(ctx, tx, jobs.SendTransactionalEmailArgs{
		Data: email.TransactionalData{
			To:       emailAddress,
			From:     "newsletter@mortenvistisen.com",
			Subject:  "Verify your newsletter subscription",
			HTMLBody: html,
			TextBody: text,
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("enqueue subscriber verification email: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit subscriber verification transaction: %w", err)
	}

	return nil
}

type VerifySubscriberData struct {
	Code string
}

func VerifySubscriber(
	ctx context.Context,
	db storage.Pool,
	pepper string,
	data VerifySubscriberData,
) (models.Subscriber, error) {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return models.Subscriber{}, err
	}
	defer tx.Rollback(ctx)

	token, err := models.FindTokenByScopeAndHash(
		ctx,
		tx,
		pepper,
		subscriberEmailVerification,
		data.Code,
	)
	if err != nil {
		return models.Subscriber{}, ErrSubscriberVerificationInvalidCode
	}

	if !token.IsValid(data.Code, pepper) {
		return models.Subscriber{}, ErrSubscriberVerificationExpiredCode
	}

	var meta map[string]string
	if err := json.Unmarshal(token.MetaData, &meta); err != nil {
		return models.Subscriber{}, err
	}

	rawID, ok := meta["subscriber_id"]
	if !ok {
		return models.Subscriber{}, errors.New("token metadata missing subscriber_id")
	}

	id, err := strconv.ParseInt(rawID, 10, 32)
	if err != nil {
		return models.Subscriber{}, err
	}

	subscriber, err := models.FindSubscriber(ctx, tx, int32(id))
	if err != nil {
		return models.Subscriber{}, err
	}

	if !subscriber.IsVerified {
		subscriber, err = models.UpdateSubscriber(ctx, tx, models.UpdateSubscriberData{
			ID:           subscriber.ID,
			Email:        subscriber.Email,
			SubscribedAt: subscriber.SubscribedAt,
			Referer:      subscriber.Referer,
			IsVerified:   true,
		})
		if err != nil {
			return models.Subscriber{}, err
		}
	}

	if err := models.DestroyToken(ctx, tx, token.ID); err != nil {
		return models.Subscriber{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Subscriber{}, err
	}

	return subscriber, nil
}
