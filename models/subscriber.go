package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

type PaginatedSubscribers struct {
	Subscribers []Subscriber
	Total       int64
	Page        int32
	PageSize    int32
}

type NewSubscriberPayload struct {
	Email        string    `validate:"required,email"`
	SubscribedAt time.Time `validate:"required"`
	Referer      string    `validate:"required"`
}

func NewSubscriber(
	ctx context.Context,
	dbtx db.DBTX,
	payload NewSubscriberPayload,
) (Subscriber, error) {
	if err := validate.Struct(payload); err != nil {
		return Subscriber{}, errors.Join(ErrDomainValidation, err)
	}

	now := time.Now()
	params := db.InsertSubscriberParams{
		ID:         uuid.New(),
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		Email:      sql.NullString{String: payload.Email, Valid: true},
		IsVerified: pgtype.Bool{Bool: false, Valid: true},
		SubscribedAt: pgtype.Timestamptz{
			Time:  payload.SubscribedAt,
			Valid: true,
		},
		Referer: sql.NullString{String: payload.Referer, Valid: true},
	}

	newSubscriber, err := db.Stmts.InsertSubscriber(ctx, dbtx, params)
	if err != nil {
		return Subscriber{}, err
	}

	return convertDBSubscriber(newSubscriber), nil
}

func GetSubscriberByID(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Subscriber, error) {
	dbSub, err := db.Stmts.QuerySubscriberByID(ctx, dbtx, id)
	if err != nil {
		return Subscriber{}, err
	}

	return convertDBSubscriber(dbSub), nil
}

func GetSubscriberByEmail(
	ctx context.Context,
	dbtx db.DBTX,
	email string,
) (Subscriber, error) {
	dbSub, err := db.Stmts.QuerySubscriberByEmail(
		ctx,
		dbtx,
		sql.NullString{String: email, Valid: true},
	)
	if err != nil {
		return Subscriber{}, err
	}

	return convertDBSubscriber(dbSub), nil
}

func GetAllSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Subscriber, error) {
	dbSubs, err := db.Stmts.QuerySubscribers(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(dbSubs))
	for i, dbSub := range dbSubs {
		subscribers[i] = convertDBSubscriber(dbSub)
	}

	return subscribers, nil
}

func GetSubscribersPage(
	ctx context.Context,
	dbtx db.DBTX,
	page, pageSize int32,
) (PaginatedSubscribers, error) {
	total, err := db.Stmts.QuerySubscribersCount(ctx, dbtx)
	if err != nil {
		return PaginatedSubscribers{}, err
	}

	offset := (page - 1) * pageSize
	dbSubs, err := db.Stmts.QuerySubscribersPage(
		ctx,
		dbtx,
		db.QuerySubscribersPageParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedSubscribers{}, err
	}

	subscribers := make([]Subscriber, len(dbSubs))
	for i, dbSub := range dbSubs {
		subscribers[i] = convertDBSubscriber(dbSub)
	}

	return PaginatedSubscribers{
		Subscribers: subscribers,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

func GetVerifiedSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Subscriber, error) {
	dbSubs, err := db.Stmts.QueryVerifiedSubscribers(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(dbSubs))
	for i, dbSub := range dbSubs {
		subscribers[i] = convertDBSubscriber(dbSub)
	}

	return subscribers, nil
}

func GetUnverifiedSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Subscriber, error) {
	dbSubs, err := db.Stmts.QueryUnverifiedSubscribers(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(dbSubs))
	for i, dbSub := range dbSubs {
		subscribers[i] = convertDBSubscriber(dbSub)
	}

	return subscribers, nil
}

func GetNewVerifiedSubsCurrentMonth(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	now := carbon.Now()

	count, err := db.Stmts.QueryVerifiedSubscriberCountByMonth(
		ctx,
		dbtx,
		db.QueryVerifiedSubscriberCountByMonthParams{
			StartMonth: pgtype.Timestamp{
				Time:  now.StartOfMonth().StdTime(),
				Valid: true,
			},
			EndMonth: pgtype.Timestamp{
				Time:  now.EndOfMonth().StdTime(),
				Valid: true,
			},
		},
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

type UpdateSubscriberPayload struct {
	ID           uuid.UUID `validate:"required,uuid"`
	Email        string    `validate:"required,email"`
	SubscribedAt time.Time `validate:"required"`
	Referer      string    `validate:"required"`
	IsVerified   bool
}

func UpdateSubscriber(
	ctx context.Context,
	dbtx db.DBTX,
	payload UpdateSubscriberPayload,
) (Subscriber, error) {
	if err := validate.Struct(payload); err != nil {
		return Subscriber{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpdateSubscriberParams{
		ID:        payload.ID,
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Email:     sql.NullString{String: payload.Email, Valid: true},
		SubscribedAt: pgtype.Timestamptz{
			Time:  payload.SubscribedAt,
			Valid: true,
		},
		Referer:    sql.NullString{String: payload.Referer, Valid: true},
		IsVerified: pgtype.Bool{Bool: payload.IsVerified, Valid: true},
	}

	updatedSubscriber, err := db.Stmts.UpdateSubscriber(ctx, dbtx, params)
	if err != nil {
		return Subscriber{}, err
	}

	return convertDBSubscriber(updatedSubscriber), nil
}

func UpdateSubscriberVerification(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
	isVerified bool,
) (Subscriber, error) {
	params := db.UpdateSubscriberVerificationParams{
		ID:         id,
		UpdatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		IsVerified: pgtype.Bool{Bool: isVerified, Valid: true},
	}

	verifiedSub, err := db.Stmts.UpdateSubscriberVerification(ctx, dbtx, params)
	if err != nil {
		return Subscriber{}, err
	}

	return convertDBSubscriber(verifiedSub), nil
}

func DeleteSubscriber(ctx context.Context, dbtx db.DBTX, id uuid.UUID) error {
	return db.Stmts.DeleteSubscriber(ctx, dbtx, id)
}

func DeleteSubscriberByEmail(
	ctx context.Context,
	dbtx db.DBTX,
	email string,
) error {
	return db.Stmts.DeleteSubscriberByEmail(
		ctx,
		dbtx,
		sql.NullString{String: email, Valid: true},
	)
}

func GetRecentSubscribers(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Subscriber, error) {
	rows, err := db.Stmts.QueryRecentSubscribers(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(rows))
	for i, dbSub := range rows {
		subscribers[i] = convertDBSubscriber(dbSub)
	}

	return subscribers, nil
}

func convertDBSubscriber(dbSub db.Subscriber) Subscriber {
	return Subscriber{
		ID:           dbSub.ID,
		CreatedAt:    dbSub.CreatedAt.Time,
		UpdatedAt:    dbSub.UpdatedAt.Time,
		Email:        dbSub.Email.String,
		SubscribedAt: dbSub.SubscribedAt.Time,
		Referer:      dbSub.Referer.String,
		IsVerified:   dbSub.IsVerified.Bool,
	}
}

func ClearOldUnverifiedSubs(ctx context.Context, dbtx db.DBTX) error {
	olderThan := time.Now().AddDate(0, -1, 0)
	return db.Stmts.DeleteSubsOlderThanMonth(
		ctx,
		dbtx,
		pgtype.Timestamp{
			Time:  olderThan,
			Valid: true,
		},
	)
}
