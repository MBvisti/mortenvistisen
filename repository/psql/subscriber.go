package psql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p Postgres) QuerySubscriberByID(
	ctx context.Context,
	id uuid.UUID,
) (domain.Subscriber, error) {
	subscriber, err := p.Queries.QuerySubscriberByID(ctx, id)
	if err != nil {
		return domain.Subscriber{}, err
	}

	return domain.Subscriber{
		ID:           subscriber.ID,
		CreatedAt:    subscriber.CreatedAt.Time,
		UpdatedAt:    subscriber.UpdatedAt.Time,
		Email:        subscriber.Email.String,
		SubscribedAt: subscriber.SubscribedAt.Time,
		Referer:      subscriber.Referer.String,
		IsVerified:   subscriber.IsVerified.Bool,
	}, nil
}

func (p Postgres) QuerySubscriberByEmail(
	ctx context.Context,
	email string,
) (domain.Subscriber, error) {
	subscriber, err := p.Queries.QuerySubscriberByEmail(
		ctx,
		sql.NullString{String: email, Valid: true},
	)
	if err != nil {
		return domain.Subscriber{}, err
	}

	return domain.Subscriber{
		ID:           subscriber.ID,
		CreatedAt:    subscriber.CreatedAt.Time,
		UpdatedAt:    subscriber.UpdatedAt.Time,
		Email:        subscriber.Email.String,
		SubscribedAt: subscriber.SubscribedAt.Time,
		Referer:      subscriber.Referer.String,
		IsVerified:   subscriber.IsVerified.Bool,
	}, nil
}

func (p Postgres) InsertSubscriber(
	ctx context.Context,
	data domain.Subscriber,
) (domain.Subscriber, error) {
	createdAt := pgtype.Timestamptz{
		Time:  data.CreatedAt,
		Valid: true,
	}
	updatedAt := pgtype.Timestamptz{
		Time:  data.UpdatedAt,
		Valid: true,
	}
	subscribedAt := pgtype.Timestamptz{
		Time:  data.SubscribedAt,
		Valid: true,
	}

	newSubscriber, err := p.Queries.InsertSubscriber(ctx, database.InsertSubscriberParams{
		ID:        data.ID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Email: sql.NullString{
			String: data.Email,
			Valid:  true,
		},
		SubscribedAt: subscribedAt,
		Referer: sql.NullString{
			String: data.Referer,
			Valid:  true,
		},
		IsVerified: pgtype.Bool{
			Bool:  data.IsVerified,
			Valid: true,
		},
	})
	if err != nil {
		return domain.Subscriber{}, errors.Join(ErrInternalDBErr, err)
	}

	return domain.Subscriber{
		ID:           newSubscriber.ID,
		CreatedAt:    newSubscriber.CreatedAt.Time,
		UpdatedAt:    newSubscriber.UpdatedAt.Time,
		Email:        newSubscriber.Email.String,
		SubscribedAt: newSubscriber.SubscribedAt.Time,
		Referer:      newSubscriber.Referer.String,
		IsVerified:   newSubscriber.IsVerified.Bool,
	}, nil
}

func (p Postgres) UpdateSubscriber(
	ctx context.Context,
	data domain.Subscriber,
) (domain.Subscriber, error) {
	subscriberToUpdate, err := p.Queries.QuerySubscriberByID(ctx, data.ID)
	if err != nil {
		return domain.Subscriber{}, err
	}

	updateParams := database.UpdateSubscriberParams{
		UpdatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		Email:        subscriberToUpdate.Email,
		SubscribedAt: subscriberToUpdate.SubscribedAt,
		Referer:      subscriberToUpdate.Referer,
		IsVerified:   subscriberToUpdate.IsVerified,
	}

	switch {
	case data.Email != "":
		updateParams.Email = sql.NullString{String: data.Email, Valid: true}
	case data.SubscribedAt != time.Time{}:
		updateParams.SubscribedAt = pgtype.Timestamptz{
			Time:  data.SubscribedAt,
			Valid: true,
		}
	case data.Referer != "":
		updateParams.Referer = sql.NullString{String: data.Referer, Valid: true}
	}

	updatedSubscriber, err := p.Queries.UpdateSubscriber(ctx, updateParams)
	if err != nil {
		return domain.Subscriber{}, errors.Join(ErrInternalDBErr, err)
	}

	return domain.Subscriber{
		ID:           updatedSubscriber.ID,
		CreatedAt:    updatedSubscriber.CreatedAt.Time,
		UpdatedAt:    updatedSubscriber.UpdatedAt.Time,
		Email:        updatedSubscriber.Email.String,
		SubscribedAt: updatedSubscriber.SubscribedAt.Time,
		Referer:      updatedSubscriber.Referer.String,
		IsVerified:   updatedSubscriber.IsVerified.Bool,
	}, nil
}

func (p Postgres) QueryNewSubscribersByMonth(ctx context.Context) ([]domain.Subscriber, error) {
	subs, err := p.Queries.QueryNewSubscribersInCurrentMonth(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Join(ErrInternalDBErr, err)
	}

	subscribers := make([]domain.Subscriber, len(subs))
	for i, sub := range subs {
		subscribers[i] = domain.Subscriber{
			ID:           sub.ID,
			CreatedAt:    sub.CreatedAt.Time,
			UpdatedAt:    sub.UpdatedAt.Time,
			Email:        sub.Email.String,
			SubscribedAt: sub.SubscribedAt.Time,
			Referer:      sub.Referer.String,
			IsVerified:   sub.IsVerified.Bool,
		}
	}

	return subscribers, nil
}

func (p Postgres) ListSubscribers(
	ctx context.Context,
	filters models.QueryFilters,
	opts ...models.PaginationOption,
) ([]domain.Subscriber, error) {
	options := &models.PaginationOptions{}

	for _, opt := range opts {
		opt(options)
	}

	params := database.QuerySubscribersParams{
		Offset: sql.NullInt32{Int32: options.Offset, Valid: true},
		Limit:  sql.NullInt32{Int32: options.Limit, Valid: true},
	}

	for k, v := range filters {
		if k == "IsVerified" {
			val, ok := v.(bool)
			if ok {
				params.IsVerified = pgtype.Bool{Bool: val, Valid: true}
			}
		}
	}

	subs, err := p.Queries.QuerySubscribers(ctx, params)
	if err != nil {
		return nil, err
	}

	subscribers := make([]domain.Subscriber, len(subs))
	for i, sub := range subs {
		subscribers[i] = domain.Subscriber{
			ID:           sub.ID,
			CreatedAt:    sub.CreatedAt.Time,
			UpdatedAt:    sub.UpdatedAt.Time,
			Email:        sub.Email.String,
			SubscribedAt: sub.SubscribedAt.Time,
			Referer:      sub.Referer.String,
			IsVerified:   sub.IsVerified.Bool,
		}
	}

	return subscribers, nil
}

func (p Postgres) CountSubscribers(
	ctx context.Context,
) (int64, error) {
	count, err := p.Queries.QuerySubscriberCount(ctx)
	if err != nil {
		return 0, errors.Join(ErrInternalDBErr, err)
	}

	return count, nil
}

func (p Postgres) CountSubscribersByStatus(
	ctx context.Context,
	verified bool,
) (int64, error) {
	count, err := p.Queries.QuerySubscriberCountByStatus(ctx, pgtype.Bool{
		Bool:  verified,
		Valid: true,
	})
	if err != nil {
		return 0, errors.Join(ErrInternalDBErr, err)
	}

	return count, nil
}

func (p Postgres) DeleteSubscriber(ctx context.Context, subscriberID uuid.UUID) error {
	if err := p.Queries.DeleteSubscriberTokenBySubscriberID(ctx, subscriberID); err != nil {
		return err
	}

	if err := p.Queries.DeleteSubscriber(ctx, subscriberID); err != nil {
		return err
	}

	return nil
}
