package models

import (
	"context"
	"database/sql"
	"errors"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"
	"net/mail"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

const (
	subscriberEmailMaxLength   = 254
	subscriberRefererMaxLength = 2048
)

type SubscriberEntity struct {
	bun.BaseModel `bun:"table:subscribers,alias:subscribers"`
	ID            int32          `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time      `bun:"created_at"`
	UpdatedAt     time.Time      `bun:"updated_at"`
	Email         sql.NullString `bun:"email"`
	SubscribedAt  sql.NullTime   `bun:"subscribed_at"`
	Referer       sql.NullString `bun:"referer"`
	IsVerified    sql.NullBool   `bun:"is_verified"`
}

func (e *SubscriberEntity) validationBuilder() *validation.ValidationBuilder {
	b := validation.NewBuilder()
	b.Required("email", e.Email)
	b.MaxLen("email", e.Email, subscriberEmailMaxLength)
	b.AddRule("email", "email", "must be a valid email address")
	if e.Email.Valid {
		emailAddress := strings.TrimSpace(e.Email.String)
		if emailAddress != "" {
			parsedAddress, err := mail.ParseAddress(emailAddress)
			if err != nil || parsedAddress.Address != emailAddress {
				b.AddField("email", "email", "must be a valid email address")
			}
		}
	}
	b.Required("subscribedAt", e.SubscribedAt)
	b.MaxLen("referer", e.Referer, subscriberRefererMaxLength)

	return b
}

func (e *SubscriberEntity) Validate() error {
	return e.validationBuilder().Err()
}

func SubscriberValidationRules() validation.Rules {
	return (&SubscriberEntity{}).validationBuilder().Rules()
}

func (s subscriber) Find(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (SubscriberEntity, error) {
	var entity SubscriberEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SubscriberEntity{}, ErrNotFound
		}
		return SubscriberEntity{}, err
	}

	return entity, nil
}

type CreateSubscriberData struct {
	Email        sql.NullString
	SubscribedAt sql.NullTime
	Referer      sql.NullString
	IsVerified   sql.NullBool
}

func (s subscriber) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateSubscriberData,
) (SubscriberEntity, error) {
	email := data.Email
	if email.Valid {
		email.String = strings.ToLower(strings.TrimSpace(email.String))
	}

	entity := SubscriberEntity{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        email,
		SubscribedAt: data.SubscribedAt,
		Referer:      normalizeOptionalString(data.Referer),
		IsVerified:   data.IsVerified,
	}

	if err := validation.Validate(&entity); err != nil {
		return SubscriberEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return SubscriberEntity{}, err
	}

	return entity, nil
}

type UpdateSubscriberData struct {
	ID           int32
	UpdatedAt    time.Time
	Email        sql.NullString
	SubscribedAt sql.NullTime
	Referer      sql.NullString
	IsVerified   sql.NullBool
}

func (s subscriber) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateSubscriberData,
) (SubscriberEntity, error) {
	email := data.Email
	if email.Valid {
		email.String = strings.ToLower(strings.TrimSpace(email.String))
	}

	entity := SubscriberEntity{
		ID:           data.ID,
		UpdatedAt:    time.Now(),
		Email:        email,
		SubscribedAt: data.SubscribedAt,
		Referer:      normalizeOptionalString(data.Referer),
		IsVerified:   data.IsVerified,
	}

	if err := validation.Validate(&entity); err != nil {
		return SubscriberEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("email").
		Column("subscribed_at").
		Column("referer").
		Column("is_verified").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		return SubscriberEntity{}, err
	}

	return entity, nil
}

func (s subscriber) MarkVerified(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (SubscriberEntity, error) {
	var entity SubscriberEntity
	if err := db.NewUpdate().
		Model(&entity).
		Set("is_verified = TRUE").
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SubscriberEntity{}, ErrNotFound
		}
		return SubscriberEntity{}, err
	}

	return entity, nil
}

func (s subscriber) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*SubscriberEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (s subscriber) DestroyMany(
	ctx context.Context,
	db storage.Executor,
	ids []int32,
) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	result, err := db.NewDelete().
		Model((*SubscriberEntity)(nil)).
		Where("id IN (?)", bun.In(ids)).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s subscriber) All(ctx context.Context, db storage.Executor) ([]SubscriberEntity, error) {
	var entities []SubscriberEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

func (s subscriber) Verified(ctx context.Context, db storage.Executor) ([]SubscriberEntity, error) {
	var entities []SubscriberEntity
	if err := db.NewSelect().
		Model(&entities).
		Where("is_verified = TRUE").
		Where("subscribed_at IS NOT NULL").
		Where("email IS NOT NULL").
		Where("email <> ''").
		OrderExpr("id ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type PaginatedSubscribers struct {
	Subscribers []SubscriberEntity
	TotalCount  int64
	Page        int64
	PageSize    int64
	TotalPages  int64
}

func (s subscriber) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
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

	totalCount, err := db.NewSelect().
		Model(&SubscriberEntity{}).Count(ctx)
	if err != nil {
		return PaginatedSubscribers{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * pageSize

	entities := make([]SubscriberEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		OrderExpr("updated_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedSubscribers{}, err
	}

	return PaginatedSubscribers{
		Subscribers: entities,
		TotalCount:  int64(totalCount),
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}, nil
}

func (s subscriber) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateSubscriberData,
) (SubscriberEntity, error) {
	email := data.Email
	if email.Valid {
		email.String = strings.ToLower(strings.TrimSpace(email.String))
	}

	entity := SubscriberEntity{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        email,
		SubscribedAt: data.SubscribedAt,
		Referer:      normalizeOptionalString(data.Referer),
		IsVerified:   data.IsVerified,
	}

	if err := validation.Validate(&entity); err != nil {
		return SubscriberEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (email) WHERE email IS NOT NULL DO UPDATE").
		Set("updated_at = excluded.updated_at").
		Set("subscribed_at = excluded.subscribed_at").
		Set("referer = excluded.referer").
		Returning("*").
		Scan(ctx); err != nil {
		return SubscriberEntity{}, err
	}

	return entity, nil
}
