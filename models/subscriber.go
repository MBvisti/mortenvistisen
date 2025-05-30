package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type Subscriber struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	SubscribedAt time.Time
	Referer      string
	IsVerified   bool
}

func (s Subscriber) IsActive() bool {
	return s.IsVerified && !s.SubscribedAt.IsZero()
}

func GetSubscriberByEmail(
	ctx context.Context,
	dbtx db.DBTX,
	email string,
) (Subscriber, error) {
	row, err := db.Stmts.QuerySubscriberByEmail(
		ctx,
		dbtx,
		sql.NullString{String: email, Valid: email != ""},
	)
	if err != nil {
		return Subscriber{}, err
	}

	return Subscriber{
		ID:           row.ID,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		Email:        row.Email.String,
		SubscribedAt: row.SubscribedAt.Time,
		Referer:      row.Referer.String,
		IsVerified:   row.IsVerified.Bool,
	}, nil
}

func GetSubscriber(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Subscriber, error) {
	row, err := db.Stmts.QuerySubscriberByID(ctx, dbtx, id)
	if err != nil {
		return Subscriber{}, err
	}

	return Subscriber{
		ID:           row.ID,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		Email:        row.Email.String,
		SubscribedAt: row.SubscribedAt.Time,
		Referer:      row.Referer.String,
		IsVerified:   row.IsVerified.Bool,
	}, nil
}

func GetAllSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Subscriber, error) {
	rows, err := db.Stmts.QuerySubscribers(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(rows))
	for i, row := range rows {
		subscribers[i] = Subscriber{
			ID:           row.ID,
			CreatedAt:    row.CreatedAt.Time,
			UpdatedAt:    row.UpdatedAt.Time,
			Email:        row.Email.String,
			SubscribedAt: row.SubscribedAt.Time,
			Referer:      row.Referer.String,
			IsVerified:   row.IsVerified.Bool,
		}
	}

	return subscribers, nil
}

func GetVerifiedSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Subscriber, error) {
	rows, err := db.Stmts.QueryVerifiedSubscribers(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(rows))
	for i, row := range rows {
		subscribers[i] = Subscriber{
			ID:           row.ID,
			CreatedAt:    row.CreatedAt.Time,
			UpdatedAt:    row.UpdatedAt.Time,
			Email:        row.Email.String,
			SubscribedAt: row.SubscribedAt.Time,
			Referer:      row.Referer.String,
			IsVerified:   row.IsVerified.Bool,
		}
	}

	return subscribers, nil
}

type NewSubscriberPayload struct {
	Email        string `validate:"required,email"`
	SubscribedAt time.Time
	Referer      string
}

func NewSubscriber(
	ctx context.Context,
	dbtx db.DBTX,
	data NewSubscriberPayload,
) (Subscriber, error) {
	if err := validate.Struct(data); err != nil {
		return Subscriber{}, errors.Join(ErrDomainValidation, err)
	}

	subscriber := Subscriber{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        data.Email,
		SubscribedAt: data.SubscribedAt,
		Referer:      data.Referer,
		IsVerified:   false,
	}

	_, err := db.Stmts.InsertSubscriber(ctx, dbtx, db.InsertSubscriberParams{
		ID:        subscriber.ID,
		CreatedAt: pgtype.Timestamptz{Time: subscriber.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: subscriber.UpdatedAt, Valid: true},
		Email: sql.NullString{
			String: subscriber.Email,
			Valid:  subscriber.Email != "",
		},
		SubscribedAt: pgtype.Timestamptz{
			Time:  subscriber.SubscribedAt,
			Valid: !subscriber.SubscribedAt.IsZero(),
		},
		Referer: sql.NullString{
			String: subscriber.Referer,
			Valid:  subscriber.Referer != "",
		},
		IsVerified: pgtype.Bool{Bool: subscriber.IsVerified, Valid: true},
	})
	if err != nil {
		return Subscriber{}, err
	}

	return subscriber, nil
}

type UpdateSubscriberPayload struct {
	ID        uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
	Email     string    `validate:"required,email"`
	Referer   string
}

func UpdateSubscriber(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateSubscriberPayload,
) (Subscriber, error) {
	if err := validate.Struct(data); err != nil {
		return Subscriber{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateSubscriber(ctx, dbtx, db.UpdateSubscriberParams{
		ID: data.ID,
		UpdatedAt: pgtype.Timestamptz{
			Time:  data.UpdatedAt,
			Valid: true,
		},
		Email: sql.NullString{String: data.Email, Valid: data.Email != ""},
		Referer: sql.NullString{
			String: data.Referer,
			Valid:  data.Referer != "",
		},
	})
	if err != nil {
		return Subscriber{}, err
	}

	return Subscriber{
		ID:           row.ID,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		Email:        row.Email.String,
		SubscribedAt: row.SubscribedAt.Time,
		Referer:      row.Referer.String,
		IsVerified:   row.IsVerified.Bool,
	}, nil
}

type VerifySubscriberPayload struct {
	ID         uuid.UUID `validate:"required,uuid"`
	UpdatedAt  time.Time `validate:"required"`
	IsVerified bool      `validate:"required"`
}

func VerifySubscriber(
	ctx context.Context,
	dbtx db.DBTX,
	data VerifySubscriberPayload,
) error {
	if err := validate.Struct(data); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}

	return db.Stmts.VerifySubscriber(ctx, dbtx, db.VerifySubscriberParams{
		ID: data.ID,
		UpdatedAt: pgtype.Timestamptz{
			Time:  data.UpdatedAt,
			Valid: true,
		},
		IsVerified: pgtype.Bool{Bool: data.IsVerified, Valid: true},
	})
}

func DeleteSubscriber(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.Stmts.DeleteSubscriber(ctx, dbtx, id)
}

func CountSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	return db.Stmts.CountSubscribers(ctx, dbtx)
}

func CountVerifiedSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	return db.Stmts.CountVerifiedSubscribers(ctx, dbtx)
}
