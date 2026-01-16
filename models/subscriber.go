package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models/internal/db"
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

func FindSubscriber(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) (Subscriber, error) {
	row, err := queries.QuerySubscriberByID(ctx, exec, id)
	if err != nil {
		return Subscriber{}, err
	}

	return rowToSubscriber(row), nil
}

type CreateSubscriberData struct {
	Email        string
	SubscribedAt time.Time
	Referer      string
	IsVerified   bool
}

func CreateSubscriber(
	ctx context.Context,
	exec storage.Executor,
	data CreateSubscriberData,
) (Subscriber, error) {
	if err := Validate.Struct(data); err != nil {
		return Subscriber{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.InsertSubscriberParams{
		ID:           uuid.New(),
		Email:        pgtype.Text{String: data.Email, Valid: true},
		SubscribedAt: pgtype.Timestamptz{Time: data.SubscribedAt, Valid: true},
		Referer:      pgtype.Text{String: data.Referer, Valid: true},
		IsVerified:   pgtype.Bool{Bool: data.IsVerified, Valid: true},
	}
	row, err := queries.InsertSubscriber(ctx, exec, params)
	if err != nil {
		return Subscriber{}, err
	}

	return rowToSubscriber(row), nil
}

type UpdateSubscriberData struct {
	ID           uuid.UUID
	UpdatedAt    time.Time
	Email        string
	SubscribedAt time.Time
	Referer      string
	IsVerified   bool
}

func UpdateSubscriber(
	ctx context.Context,
	exec storage.Executor,
	data UpdateSubscriberData,
) (Subscriber, error) {
	if err := Validate.Struct(data); err != nil {
		return Subscriber{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpdateSubscriberParams{
		ID:           data.ID,
		Email:        pgtype.Text{String: data.Email, Valid: true},
		SubscribedAt: pgtype.Timestamptz{Time: data.SubscribedAt, Valid: true},
		Referer:      pgtype.Text{String: data.Referer, Valid: true},
		IsVerified:   pgtype.Bool{Bool: data.IsVerified, Valid: true},
	}

	row, err := queries.UpdateSubscriber(ctx, exec, params)
	if err != nil {
		return Subscriber{}, err
	}

	return rowToSubscriber(row), nil
}

func DestroySubscriber(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) error {
	return queries.DeleteSubscriber(ctx, exec, id)
}

func AllSubscribers(
	ctx context.Context,
	exec storage.Executor,
) ([]Subscriber, error) {
	rows, err := queries.QuerySubscribers(ctx, exec)
	if err != nil {
		return nil, err
	}

	subscribers := make([]Subscriber, len(rows))
	for i, row := range rows {
		subscribers[i] = rowToSubscriber(row)
	}

	return subscribers, nil
}

type PaginatedSubscribers struct {
	Subscribers []Subscriber
	TotalCount  int64
	Page        int64
	PageSize    int64
	TotalPages  int64
}

func PaginateSubscribers(
	ctx context.Context,
	exec storage.Executor,
	page int64,
	pageSize int64,
) (PaginatedSubscribers, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := queries.CountSubscribers(ctx, exec)
	if err != nil {
		return PaginatedSubscribers{}, err
	}

	rows, err := queries.QueryPaginatedSubscribers(
		ctx,
		exec,
		db.QueryPaginatedSubscribersParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedSubscribers{}, err
	}

	subscribers := make([]Subscriber, len(rows))
	for i, row := range rows {
		subscribers[i] = rowToSubscriber(row)
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedSubscribers{
		Subscribers: subscribers,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

func UpsertSubscriber(
	ctx context.Context,
	exec storage.Executor,
	data CreateSubscriberData,
) (Subscriber, error) {
	if err := Validate.Struct(data); err != nil {
		return Subscriber{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpsertSubscriberParams{
		ID:           uuid.New(),
		Email:        pgtype.Text{String: data.Email, Valid: true},
		SubscribedAt: pgtype.Timestamptz{Time: data.SubscribedAt, Valid: true},
		Referer:      pgtype.Text{String: data.Referer, Valid: true},
		IsVerified:   pgtype.Bool{Bool: data.IsVerified, Valid: true},
	}
	row, err := queries.UpsertSubscriber(ctx, exec, params)
	if err != nil {
		return Subscriber{}, err
	}

	return rowToSubscriber(row), nil
}

func rowToSubscriber(row db.Subscriber) Subscriber {
	return Subscriber{
		ID:           row.ID,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		Email:        row.Email.String,
		SubscribedAt: row.SubscribedAt.Time,
		Referer:      row.Referer.String,
		IsVerified:   row.IsVerified.Bool,
	}
}
