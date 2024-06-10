package models

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type subscriberStorage interface {
	QuerySubscriberByID(ctx context.Context, id uuid.UUID) (domain.Subscriber, error)
	QuerySubscriberByEmail(ctx context.Context, email string) (domain.Subscriber, error)
	ListSubscribers(
		ctx context.Context,
		opts ...PaginationOption,
	) ([]domain.Subscriber, error)
	CountSubscribers(
		ctx context.Context,
	) (int64, error)
	CountSubscribersByStatus(
		ctx context.Context,
		verified bool,
	) (int64, error)
	QueryNewSubscribersByMonth(ctx context.Context) ([]domain.Subscriber, error)
	UpdateSubscriber(
		ctx context.Context,
		id uuid.UUID,
		data domain.Subscriber,
	) (domain.Subscriber, error)
	InsertSubscriber(
		ctx context.Context,
		data domain.Subscriber,
	) (domain.Subscriber, error)
}

type subscriberEmailService interface {
	SendNewSubscriberEmail(
		ctx context.Context,
		subscriberEmail string,
		activationToken, unsubscribeToken domain.Token,
	) error
}

type subscriberTokenService interface {
	CreateSubscriptionToken(
		ctx context.Context,
		subscriberID uuid.UUID,
	) (domain.Token, error)
	CreateUnsubscribeToken(
		ctx context.Context,
		subscriberID uuid.UUID,
	) (domain.Token, error)
}

type SubscriberService struct {
	emailService subscriberEmailService
	tknService   subscriberTokenService
	validator    *validator.Validate
	storage      subscriberStorage
}

func NewSubscriberSvc(
	emailService subscriberEmailService,
	tknService subscriberTokenService,
	validator *validator.Validate,
	storage subscriberStorage,
) SubscriberService {
	return SubscriberService{
		emailService,
		tknService,
		validator,
		storage,
	}
}

func (svc *SubscriberService) ByID(ctx context.Context, id uuid.UUID) (domain.Subscriber, error) {
	subscriber, err := svc.storage.QuerySubscriberByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscriber{}, ErrNoRowWithIdentifier
		}

		return domain.Subscriber{}, err
	}

	return domain.Subscriber{
		ID:           subscriber.ID,
		CreatedAt:    subscriber.CreatedAt,
		UpdatedAt:    subscriber.UpdatedAt,
		Email:        subscriber.Email,
		SubscribedAt: subscriber.SubscribedAt,
		Referer:      subscriber.Referer,
		IsVerified:   subscriber.IsVerified,
	}, nil
}

func (svc *SubscriberService) ByEmail(
	ctx context.Context,
	email string,
) (domain.Subscriber, error) {
	subscriber, err := svc.storage.QuerySubscriberByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscriber{}, ErrNoRowWithIdentifier
		}

		return domain.Subscriber{}, err
	}

	return domain.Subscriber{
		ID:           subscriber.ID,
		CreatedAt:    subscriber.CreatedAt,
		UpdatedAt:    subscriber.UpdatedAt,
		Email:        subscriber.Email,
		SubscribedAt: subscriber.SubscribedAt,
		Referer:      subscriber.Referer,
		IsVerified:   subscriber.IsVerified,
	}, nil
}

func (svc *SubscriberService) List(
	ctx context.Context,
	offset int32,
	limit int32,
) ([]domain.Subscriber, error) {
	subs, err := svc.storage.ListSubscribers(ctx, WithPagination(limit, offset))
	if err != nil {
		return nil, err
	}

	subscribers := make([]domain.Subscriber, len(subs))
	for i, sub := range subs {
		subscribers[i] = domain.Subscriber{
			ID:           sub.ID,
			CreatedAt:    sub.CreatedAt,
			UpdatedAt:    sub.UpdatedAt,
			Email:        sub.Email,
			SubscribedAt: sub.SubscribedAt,
			Referer:      sub.Referer,
			IsVerified:   sub.IsVerified,
		}
	}

	return subscribers, nil
}

func (svc *SubscriberService) Count(ctx context.Context) (int64, error) {
	count, err := svc.storage.CountSubscribers(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (svc *SubscriberService) NewForCurrentMonth(ctx context.Context) ([]domain.Subscriber, error) {
	subs, err := svc.storage.QueryNewSubscribersByMonth(ctx)
	if err != nil {
		return nil, err
	}

	subscribers := make([]domain.Subscriber, len(subs))
	for i, sub := range subs {
		subscribers[i] = domain.Subscriber{
			ID:           sub.ID,
			CreatedAt:    sub.CreatedAt,
			UpdatedAt:    sub.UpdatedAt,
			Email:        sub.Email,
			SubscribedAt: sub.SubscribedAt,
			Referer:      sub.Referer,
			IsVerified:   sub.IsVerified,
		}
	}

	return subscribers, nil
}

func (svc *SubscriberService) UnverifiedCount(ctx context.Context) (int64, error) {
	count, err := svc.storage.CountSubscribersByStatus(ctx, false)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (svc *SubscriberService) Verify(ctx context.Context) error {
	return nil
}

// New creates a new subscriber and sends them an verification email
// TODO: create job to create and send email on queue
func (svc *SubscriberService) New(ctx context.Context, email, articleTitle string) error {
	// 1. create type
	// 2. store subscriber
	// 3. create token
	// 4. create email & send email

	// _, err := svc.ByEmail(ctx, email)
	err := errors.New("something happened")
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("could not query subscriber by email", "error", err, "email", email)
		return errors.Join(ErrUnrecoverableEvent, err)
	}
	if err == nil {
		return ErrSubscriberExists
	}

	// 1. create type
	subscriber, err := domain.NewSubscriber(email, articleTitle, time.Now(), false, svc.validator)
	if err != nil {
		return err
	}

	// 2. store subscriber
	if _, err := svc.storage.InsertSubscriber(ctx, subscriber); err != nil {
		slog.Error("could not insert subscriber into database", "error", err)
		return errors.Join(ErrUnrecoverableEvent, err)
	}

	// 3. create token
	activationTkn, err := svc.tknService.CreateSubscriptionToken(ctx, subscriber.ID)
	if err != nil {
		return errors.Join(ErrUnrecoverableEvent, err)
	}

	// 3. create token
	unsubscribeTkn, err := svc.tknService.CreateUnsubscribeToken(ctx, subscriber.ID)
	if err != nil {
		return errors.Join(ErrUnrecoverableEvent, err)
	}

	// 4. create email & send email
	if err := svc.emailService.SendNewSubscriberEmail(ctx, subscriber.Email, activationTkn, unsubscribeTkn); err != nil {
		return errors.Join(
			ErrUnrecoverableEvent,
			err,
		) // could be retried if something goes wrong, a little TODO
	}

	return nil
}
