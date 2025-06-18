package models

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type Newsletter struct {
	ID                 uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ReleasedAt         time.Time
	IsPublished        bool
	Title              string
	Slug               string
	Content            string
	SendStatus         string
	TotalRecipients    int
	EmailsSent         int
	SendingStartedAt   time.Time
	SendingCompletedAt time.Time
}

func GetNewsletterByID(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) (Newsletter, error) {
	row, err := db.Stmts.QueryNewsletterByID(ctx, dbtx, id)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:                 row.ID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		ReleasedAt:         row.ReleasedAt.Time,
		IsPublished:        row.IsPublished.Bool,
		Title:              row.Title,
		Slug:               row.Slug.String,
		Content:            row.Content,
		SendStatus:         row.SendStatus,
		TotalRecipients:    int(row.TotalRecipients),
		EmailsSent:         int(row.EmailsSent),
		SendingStartedAt:   row.SendingStartedAt.Time,
		SendingCompletedAt: row.SendingCompletedAt.Time,
	}, nil
}

func GetNewsletterByTitle(
	ctx context.Context,
	dbtx db.DBTX,
	title string,
) (Newsletter, error) {
	row, err := db.Stmts.QueryNewsletterByTitle(ctx, dbtx, title)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:                 row.ID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		ReleasedAt:         row.ReleasedAt.Time,
		IsPublished:        row.IsPublished.Bool,
		Title:              row.Title,
		Slug:               row.Slug.String,
		Content:            row.Content,
		SendStatus:         row.SendStatus,
		TotalRecipients:    int(row.TotalRecipients),
		EmailsSent:         int(row.EmailsSent),
		SendingStartedAt:   row.SendingStartedAt.Time,
		SendingCompletedAt: row.SendingCompletedAt.Time,
	}, nil
}

func GetNewsletterBySlug(
	ctx context.Context,
	dbtx db.DBTX,
	slug string,
) (Newsletter, error) {
	row, err := db.Stmts.QueryNewsletterBySlug(
		ctx,
		dbtx,
		sql.NullString{String: slug, Valid: slug != ""},
	)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:                 row.ID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		ReleasedAt:         row.ReleasedAt.Time,
		IsPublished:        row.IsPublished.Bool,
		Title:              row.Title,
		Slug:               row.Slug.String,
		Content:            row.Content,
		SendStatus:         row.SendStatus,
		TotalRecipients:    int(row.TotalRecipients),
		EmailsSent:         int(row.EmailsSent),
		SendingStartedAt:   row.SendingStartedAt.Time,
		SendingCompletedAt: row.SendingCompletedAt.Time,
	}, nil
}

func GetNewsletters(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Newsletter, error) {
	rows, err := db.Stmts.QueryNewsletters(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = Newsletter{
			ID:                 row.ID,
			CreatedAt:          row.CreatedAt.Time,
			UpdatedAt:          row.UpdatedAt.Time,
			ReleasedAt:         row.ReleasedAt.Time,
			IsPublished:        row.IsPublished.Bool,
			Title:              row.Title,
			Slug:               row.Slug.String,
			Content:            row.Content,
			SendStatus:         row.SendStatus,
			TotalRecipients:    int(row.TotalRecipients),
			EmailsSent:         int(row.EmailsSent),
			SendingStartedAt:   row.SendingStartedAt.Time,
			SendingCompletedAt: row.SendingCompletedAt.Time,
		}
	}

	return newsletters, nil
}

func GetPublishedNewsletters(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Newsletter, error) {
	rows, err := db.Stmts.QueryPublishedNewsletters(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = Newsletter{
			ID:                 row.ID,
			CreatedAt:          row.CreatedAt.Time,
			UpdatedAt:          row.UpdatedAt.Time,
			ReleasedAt:         row.ReleasedAt.Time,
			IsPublished:        row.IsPublished.Bool,
			Title:              row.Title,
			Slug:               row.Slug.String,
			Content:            row.Content,
			SendStatus:         row.SendStatus,
			TotalRecipients:    int(row.TotalRecipients),
			EmailsSent:         int(row.EmailsSent),
			SendingStartedAt:   row.SendingStartedAt.Time,
			SendingCompletedAt: row.SendingCompletedAt.Time,
		}
	}

	return newsletters, nil
}

func GetDraftNewsletters(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Newsletter, error) {
	rows, err := db.Stmts.QueryDraftNewsletters(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = Newsletter{
			ID:                 row.ID,
			CreatedAt:          row.CreatedAt.Time,
			UpdatedAt:          row.UpdatedAt.Time,
			ReleasedAt:         row.ReleasedAt.Time,
			IsPublished:        row.IsPublished.Bool,
			Title:              row.Title,
			Slug:               row.Slug.String,
			Content:            row.Content,
			SendStatus:         row.SendStatus,
			TotalRecipients:    int(row.TotalRecipients),
			EmailsSent:         int(row.EmailsSent),
			SendingStartedAt:   row.SendingStartedAt.Time,
			SendingCompletedAt: row.SendingCompletedAt.Time,
		}
	}

	return newsletters, nil
}

type NewsletterPaginationResult struct {
	Newsletters []Newsletter
	TotalCount  int64
	Page        int
	PageSize    int
	TotalPages  int
	HasNext     bool
	HasPrevious bool
}

func GetNewslettersPaginated(
	ctx context.Context,
	dbtx db.DBTX,
	page int,
	pageSize int,
) (NewsletterPaginationResult, error) {
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

	totalCount, err := db.Stmts.CountNewsletters(ctx, dbtx)
	if err != nil {
		return NewsletterPaginationResult{}, err
	}

	// Get paginated newsletters
	rows, err := db.Stmts.QueryNewslettersPaginated(
		ctx,
		dbtx,
		db.QueryNewslettersPaginatedParams{
			//nolint:gosec // pageSize is bounded above
			Limit: int32(pageSize),
			//nolint:gosec // offset is calculated from bounded values
			Offset: int32(
				offset,
			),
		},
	)
	if err != nil {
		return NewsletterPaginationResult{}, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = Newsletter{
			ID:                 row.ID,
			CreatedAt:          row.CreatedAt.Time,
			UpdatedAt:          row.UpdatedAt.Time,
			ReleasedAt:         row.ReleasedAt.Time,
			IsPublished:        row.IsPublished.Bool,
			Title:              row.Title,
			Slug:               row.Slug.String,
			Content:            row.Content,
			SendStatus:         row.SendStatus,
			TotalRecipients:    int(row.TotalRecipients),
			EmailsSent:         int(row.EmailsSent),
			SendingStartedAt:   row.SendingStartedAt.Time,
			SendingCompletedAt: row.SendingCompletedAt.Time,
		}
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	return NewsletterPaginationResult{
		Newsletters: newsletters,
		TotalCount:  totalCount,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

type NewNewsletterPayload struct {
	Title   string `validate:"required,max=100"`
	Slug    string `validate:"required,max=255"`
	Content string
}

func NewNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	data NewNewsletterPayload,
) (Newsletter, error) {
	if err := validate.Struct(data); err != nil {
		slog.ErrorContext(
			ctx,
			"could not validate new newsletter payload",
			"error",
			err,
			"data",
			data,
		)
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}

	newsletter := Newsletter{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     data.Title,
		Slug:      data.Slug,
		Content:   data.Content,
	}

	_, err := db.Stmts.InsertNewsletter(ctx, dbtx, db.InsertNewsletterParams{
		ID: newsletter.ID,
		CreatedAt: pgtype.Timestamptz{
			Time:  newsletter.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  newsletter.UpdatedAt,
			Valid: true,
		},
		Title: newsletter.Title,
		Slug: sql.NullString{
			String: newsletter.Slug,
			Valid:  newsletter.Slug != "",
		},
		Content: newsletter.Content,
	})
	if err != nil {
		return Newsletter{}, err
	}

	return GetNewsletterByID(ctx, dbtx, newsletter.ID)
}

type UpdateNewsletterPayload struct {
	ID          uuid.UUID `validate:"required,uuid"`
	UpdatedAt   time.Time `validate:"required"`
	IsPublished bool
	Title       string `validate:"required,max=100"`
	Slug        string `validate:"required,max=255"`
	Content     string
}

func UpdateNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateNewsletterPayload,
) (Newsletter, error) {
	if err := validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}

	_, err := db.Stmts.UpdateNewsletter(ctx, dbtx, db.UpdateNewsletterParams{
		ID:        data.ID,
		UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		Title:     data.Title,
		Slug: sql.NullString{
			String: data.Slug,
			Valid:  data.Slug != "",
		},
		Content:     data.Content,
		IsPublished: pgtype.Bool{Bool: data.IsPublished, Valid: true},
	})
	if err != nil {
		return Newsletter{}, err
	}

	return GetNewsletterByID(ctx, dbtx, data.ID)
}

type UpdateNewsletterContentPayload struct {
	ID        uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
	Content   string
}

func UpdateNewsletterContent(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateNewsletterContentPayload,
) (Newsletter, error) {
	if err := validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.UpdateNewsletterContent(
		ctx,
		dbtx,
		db.UpdateNewsletterContentParams{
			ID:        data.ID,
			UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
			Content:   data.Content,
		},
	)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:                 row.ID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		ReleasedAt:         row.ReleasedAt.Time,
		IsPublished:        row.IsPublished.Bool,
		Title:              row.Title,
		Slug:               row.Slug.String,
		Content:            row.Content,
		SendStatus:         row.SendStatus,
		TotalRecipients:    int(row.TotalRecipients),
		EmailsSent:         int(row.EmailsSent),
		SendingStartedAt:   row.SendingStartedAt.Time,
		SendingCompletedAt: row.SendingCompletedAt.Time,
	}, nil
}

type PublishNewsletterPayload struct {
	ID  uuid.UUID `validate:"required,uuid"`
	Now time.Time
}

func PublishNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	data PublishNewsletterPayload,
) (Newsletter, error) {
	if err := validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := db.Stmts.PublishNewsletter(
		ctx,
		dbtx,
		db.PublishNewsletterParams{
			ID:        data.ID,
			UpdatedAt: pgtype.Timestamptz{Time: data.Now, Valid: true},
			ReleasedAt: pgtype.Timestamptz{
				Time:  data.Now,
				Valid: true,
			},
			IsPublished: pgtype.Bool{Bool: true, Valid: true},
		},
	)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:                 row.ID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		ReleasedAt:         row.ReleasedAt.Time,
		IsPublished:        row.IsPublished.Bool,
		Title:              row.Title,
		Slug:               row.Slug.String,
		Content:            row.Content,
		SendStatus:         row.SendStatus,
		TotalRecipients:    int(row.TotalRecipients),
		EmailsSent:         int(row.EmailsSent),
		SendingStartedAt:   row.SendingStartedAt.Time,
		SendingCompletedAt: row.SendingCompletedAt.Time,
	}, nil
}

func DeleteNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	id uuid.UUID,
) error {
	return db.Stmts.DeleteNewsletter(ctx, dbtx, id)
}

func GetNewslettersReadyToSend(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Newsletter, error) {
	rows, err := db.Stmts.QueryNewslettersReadyToSend(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	newsletters := make([]Newsletter, len(rows))
	for i, row := range rows {
		newsletters[i] = Newsletter{
			ID:                 row.ID,
			CreatedAt:          row.CreatedAt.Time,
			UpdatedAt:          row.UpdatedAt.Time,
			ReleasedAt:         row.ReleasedAt.Time,
			IsPublished:        row.IsPublished.Bool,
			Title:              row.Title,
			Slug:               row.Slug.String,
			Content:            row.Content,
			SendStatus:         row.SendStatus,
			TotalRecipients:    int(row.TotalRecipients),
			EmailsSent:         int(row.EmailsSent),
			SendingStartedAt:   row.SendingStartedAt.Time,
			SendingCompletedAt: row.SendingCompletedAt.Time,
		}
	}

	return newsletters, nil
}

type MarkNewsletterReadyToSendPayload struct {
	ID              uuid.UUID `validate:"required,uuid"`
	Now             time.Time
	TotalRecipients int
}

func MarkNewsletterReadyToSend(
	ctx context.Context,
	dbtx db.DBTX,
	data MarkNewsletterReadyToSendPayload,
) (Newsletter, error) {
	if err := validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}

	// Safe conversion for TotalRecipients
	if data.TotalRecipients > 2147483647 { // max int32
		return Newsletter{}, errors.New("total recipients exceeds maximum allowed value")
	}

	row, err := db.Stmts.MarkNewsletterReadyToSend(
		ctx,
		dbtx,
		db.MarkNewsletterReadyToSendParams{
			ID:              data.ID,
			UpdatedAt:       pgtype.Timestamptz{Time: data.Now, Valid: true},
			TotalRecipients: int32(data.TotalRecipients),
		},
	)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:                 row.ID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		ReleasedAt:         row.ReleasedAt.Time,
		IsPublished:        row.IsPublished.Bool,
		Title:              row.Title,
		Slug:               row.Slug.String,
		Content:            row.Content,
		SendStatus:         row.SendStatus,
		TotalRecipients:    int(row.TotalRecipients),
		EmailsSent:         int(row.EmailsSent),
		SendingStartedAt:   row.SendingStartedAt.Time,
		SendingCompletedAt: row.SendingCompletedAt.Time,
	}, nil
}

type UpdateNewsletterSendStatusPayload struct {
	ID                 uuid.UUID `validate:"required,uuid"`
	Now                time.Time
	SendStatus         string `validate:"required,oneof=draft ready_to_send sending sent"`
	SendingStartedAt   *time.Time
	SendingCompletedAt *time.Time
	EmailsSent         int
}

func UpdateNewsletterSendStatus(
	ctx context.Context,
	dbtx db.DBTX,
	data UpdateNewsletterSendStatusPayload,
) (Newsletter, error) {
	if err := validate.Struct(data); err != nil {
		return Newsletter{}, errors.Join(ErrDomainValidation, err)
	}

	// Safe conversion for EmailsSent
	if data.EmailsSent > 2147483647 { // max int32
		return Newsletter{}, errors.New("emails sent exceeds maximum allowed value")
	}

	params := db.UpdateNewsletterSendStatusParams{
		ID:         data.ID,
		UpdatedAt:  pgtype.Timestamptz{Time: data.Now, Valid: true},
		SendStatus: data.SendStatus,
		EmailsSent: int32(data.EmailsSent),
	}

	if data.SendingStartedAt != nil {
		params.SendingStartedAt = pgtype.Timestamptz{Time: *data.SendingStartedAt, Valid: true}
	}

	if data.SendingCompletedAt != nil {
		params.SendingCompletedAt = pgtype.Timestamptz{Time: *data.SendingCompletedAt, Valid: true}
	}

	row, err := db.Stmts.UpdateNewsletterSendStatus(ctx, dbtx, params)
	if err != nil {
		return Newsletter{}, err
	}

	return Newsletter{
		ID:                 row.ID,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
		ReleasedAt:         row.ReleasedAt.Time,
		IsPublished:        row.IsPublished.Bool,
		Title:              row.Title,
		Slug:               row.Slug.String,
		Content:            row.Content,
		SendStatus:         row.SendStatus,
		TotalRecipients:    int(row.TotalRecipients),
		EmailsSent:         int(row.EmailsSent),
		SendingStartedAt:   row.SendingStartedAt.Time,
		SendingCompletedAt: row.SendingCompletedAt.Time,
	}, nil
}
