package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/email"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/queue/jobs"
	"mortenvistisen/router/routes"

	"github.com/uptrace/bun"
)

const (
	subscriberEmailVerification = "subscriber_email_verification"
	subscriberCodeLifetime      = 24 * time.Hour
)

var (
	ErrSubscriberAlreadyVerified       = errors.New("subscriber email is already verified")
	ErrInvalidSubscriberConfirmation   = errors.New("invalid subscriber confirmation code")
	ErrExpiredSubscriberConfirmation   = errors.New("subscriber confirmation code has expired")
	ErrSubscriberConfirmationEmailDiff = errors.New(
		"subscriber email no longer matches confirmation",
	)
)

type subscriberConfirmationMetadata struct {
	SubscriberID int32  `json:"subscriber_id"`
	Email        string `json:"email"`
}

type SubscriberConfirmation struct {
	db              storage.Pool
	insertOnly      queue.InsertOnly
	tokenSigningKey string
}

func NewSubscriberConfirmation(
	db storage.Pool,
	insertOnly queue.InsertOnly,
	cfg config.Config,
) SubscriberConfirmation {
	return SubscriberConfirmation{
		db:              db,
		insertOnly:      insertOnly,
		tokenSigningKey: cfg.App.TokenSigningKey,
	}
}

func (s SubscriberConfirmation) Subscribe(
	ctx context.Context,
	data models.CreateSubscriberData,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin subscriber subscription transaction: %w", err)
	}

	subscriber, err := models.Subscriber.Upsert(ctx, tx, data)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if !subscriber.IsVerified.Bool {
		if err := s.send(ctx, tx, subscriber); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit subscriber subscription transaction: %w", err)
	}

	return nil
}

func (s SubscriberConfirmation) Send(ctx context.Context, subscriberID int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin subscriber confirmation transaction: %w", err)
	}

	subscriber, err := models.Subscriber.Find(ctx, tx, subscriberID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("find subscriber for confirmation: %w", err)
	}
	if subscriber.IsVerified.Bool {
		_ = tx.Rollback()
		return ErrSubscriberAlreadyVerified
	}
	if err := s.send(ctx, tx, subscriber); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit subscriber confirmation transaction: %w", err)
	}

	return nil
}

func (s SubscriberConfirmation) send(
	ctx context.Context,
	tx bun.Tx,
	subscriber models.SubscriberEntity,
) error {
	metadata, err := json.Marshal(subscriberConfirmationMetadata{
		SubscriberID: subscriber.ID,
		Email:        subscriber.Email.String,
	})
	if err != nil {
		return fmt.Errorf("marshal subscriber confirmation metadata: %w", err)
	}
	metadataFilter, err := json.Marshal(map[string]int32{
		"subscriber_id": subscriber.ID,
	})
	if err != nil {
		return fmt.Errorf("marshal subscriber confirmation filter: %w", err)
	}
	if err := models.Token.DestroyByScopeAndMetadata(
		ctx,
		tx,
		subscriberEmailVerification,
		metadataFilter,
	); err != nil {
		return fmt.Errorf("invalidate prior subscriber confirmation codes: %w", err)
	}

	code, err := models.Token.CreateNumericCode(
		ctx,
		tx,
		s.tokenSigningKey,
		subscriberEmailVerification,
		time.Now().Add(subscriberCodeLifetime),
		metadata,
	)
	if err != nil {
		return fmt.Errorf("create subscriber confirmation code: %w", err)
	}

	args, err := subscriberConfirmationJob(subscriber, code)
	if err != nil {
		return err
	}
	if _, err := s.insertOnly.InsertTx(ctx, tx.Tx, args, nil); err != nil {
		return fmt.Errorf("queue subscriber confirmation email: %w", err)
	}

	return nil
}

func subscriberConfirmationJob(
	subscriber models.SubscriberEntity,
	code string,
) (jobs.SendTransactionalEmailArgs, error) {
	confirmationURL := routes.SubscriberConfirmationNew.FullURL(config.BaseURL)
	template := email.SubscriberConfirmation{
		ConfirmationURL:  confirmationURL,
		VerificationCode: code,
	}

	html, err := template.ToHTML()
	if err != nil {
		return jobs.SendTransactionalEmailArgs{}, fmt.Errorf(
			"render subscriber confirmation email html: %w",
			err,
		)
	}
	text, err := template.ToText()
	if err != nil {
		return jobs.SendTransactionalEmailArgs{}, fmt.Errorf(
			"render subscriber confirmation email text: %w",
			err,
		)
	}

	return jobs.SendTransactionalEmailArgs{
		Data: email.TransactionalData{
			To:       subscriber.Email.String,
			From:     config.DefaultSenderSignature,
			Subject:  "Confirm your email address",
			HTMLBody: html,
			TextBody: text,
			Metadata: map[string]string{
				"subscriber_id": fmt.Sprintf("%d", subscriber.ID),
			},
		},
	}, nil
}

type ConfirmSubscriberEmailData struct {
	Code string
}

func (s SubscriberConfirmation) Confirm(
	ctx context.Context,
	data ConfirmSubscriberEmailData,
) (models.SubscriberEntity, error) {
	data.Code = strings.ToUpper(strings.TrimSpace(data.Code))
	b := validation.NewBuilder()
	b.Required("code", data.Code)
	b.LenBetween("code", data.Code, 6, 6, "must be a 6-digit code")
	if data.Code != "" && strings.IndexFunc(data.Code, func(character rune) bool {
		return character < '0' || character > '9'
	}) >= 0 {
		b.Add("code", "numeric", "must be a 6-digit code")
	}
	if !b.Errors().Empty() {
		return models.SubscriberEntity{}, b.Errors()
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.SubscriberEntity{}, fmt.Errorf(
			"begin subscriber confirmation transaction: %w",
			err,
		)
	}

	token, err := models.Token.FindByScopeAndHash(
		ctx,
		tx,
		s.tokenSigningKey,
		subscriberEmailVerification,
		data.Code,
	)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, models.ErrNotFound) {
			return models.SubscriberEntity{}, ErrInvalidSubscriberConfirmation
		}
		return models.SubscriberEntity{}, fmt.Errorf("find subscriber confirmation code: %w", err)
	}
	if !token.IsValid(data.Code, s.tokenSigningKey) {
		_ = tx.Rollback()
		return models.SubscriberEntity{}, ErrExpiredSubscriberConfirmation
	}

	var metadata subscriberConfirmationMetadata
	if err := json.Unmarshal(token.MetaData, &metadata); err != nil {
		_ = tx.Rollback()
		return models.SubscriberEntity{}, fmt.Errorf(
			"unmarshal subscriber confirmation metadata: %w",
			err,
		)
	}

	subscriber, err := models.Subscriber.Find(ctx, tx, metadata.SubscriberID)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, models.ErrNotFound) {
			return models.SubscriberEntity{}, ErrInvalidSubscriberConfirmation
		}
		return models.SubscriberEntity{}, fmt.Errorf("find confirmation subscriber: %w", err)
	}
	if !strings.EqualFold(subscriber.Email.String, metadata.Email) {
		_ = tx.Rollback()
		return models.SubscriberEntity{}, ErrSubscriberConfirmationEmailDiff
	}

	if !subscriber.IsVerified.Bool {
		subscriber, err = models.Subscriber.MarkVerified(ctx, tx, subscriber.ID)
		if err != nil {
			_ = tx.Rollback()
			return models.SubscriberEntity{}, fmt.Errorf("mark subscriber verified: %w", err)
		}
	}

	if err := models.Token.Destroy(ctx, tx, token.ID); err != nil {
		_ = tx.Rollback()
		return models.SubscriberEntity{}, fmt.Errorf(
			"consume subscriber confirmation code: %w",
			err,
		)
	}
	if err := tx.Commit(); err != nil {
		return models.SubscriberEntity{}, fmt.Errorf("commit subscriber confirmation: %w", err)
	}

	return subscriber, nil
}
