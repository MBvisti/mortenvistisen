package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvisti/mortenvistisen/models/internal/db"
)

type NewsletterEmailSend struct {
	ID           uuid.UUID
	NewsletterID uuid.UUID
	SubscriberID uuid.UUID
	EmailAddress string
	Status       string
	SentAt       time.Time
	FailedAt     time.Time
	ErrorMessage string
	RiverJobID   int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}


func CreateNewsletterEmailSend(
	ctx context.Context,
	dbtx db.DBTX,
	newsletterID uuid.UUID,
	subscriberID uuid.UUID,
	emailAddress string,
	riverJobID int64,
) (NewsletterEmailSend, error) {
	row, err := db.Stmts.InsertNewsletterEmailSend(ctx, dbtx, db.InsertNewsletterEmailSendParams{
		NewsletterID: newsletterID,
		SubscriberID: subscriberID,
		EmailAddress: emailAddress,
		Status:       "pending",
		RiverJobID: sql.NullInt64{
			Int64: riverJobID,
			Valid: riverJobID > 0,
		},
	})
	if err != nil {
		return NewsletterEmailSend{}, err
	}

	return newsletterEmailSendRowToModel(row), nil
}

func UpdateNewsletterEmailSendStatus(
	ctx context.Context,
	dbtx db.DBTX,
	newsletterID uuid.UUID,
	subscriberID uuid.UUID,
	status string,
	errorMessage string,
) (NewsletterEmailSend, error) {
	row, err := db.Stmts.UpdateNewsletterEmailSendStatus(ctx, dbtx, db.UpdateNewsletterEmailSendStatusParams{
		NewsletterID: newsletterID,
		SubscriberID: subscriberID,
		Status:       status,
		ErrorMessage: sql.NullString{
			String: errorMessage,
			Valid:  errorMessage != "",
		},
	})
	if err != nil {
		return NewsletterEmailSend{}, err
	}

	return newsletterEmailSendRowToModel(row), nil
}

func GetNewsletterSendStats(
	ctx context.Context,
	dbtx db.DBTX,
	newsletterID uuid.UUID,
) (NewsletterSendStats, error) {
	row, err := db.Stmts.GetNewsletterSendStats(ctx, dbtx, newsletterID)
	if err != nil {
		return NewsletterSendStats{}, err
	}

	return NewsletterSendStats{
		NewsletterID:   row.NewsletterID,
		TotalEmails:    row.TotalEmails,
		SentEmails:     row.SentEmails,
		FailedEmails:   row.FailedEmails,
		BouncedEmails:  row.BouncedEmails,
		PendingEmails:  row.PendingEmails,
		CompletionRate: parseNumericToFloat64(row.CompletionRate),
	}, nil
}

func GetAllNewsletterSendStats(
	ctx context.Context,
	dbtx db.DBTX,
) (map[uuid.UUID]NewsletterSendStats, error) {
	rows, err := db.Stmts.GetAllNewsletterSendStats(ctx, dbtx)
	if err != nil {
		return nil, err
	}

	stats := make(map[uuid.UUID]NewsletterSendStats, len(rows))
	for _, row := range rows {
		stats[row.NewsletterID] = NewsletterSendStats{
			NewsletterID:   row.NewsletterID,
			TotalEmails:    row.TotalEmails,
			SentEmails:     row.SentEmails,
			FailedEmails:   row.FailedEmails,
			BouncedEmails:  row.BouncedEmails,
			PendingEmails:  row.PendingEmails,
			CompletionRate: parseNumericToFloat64(row.CompletionRate),
		}
	}

	return stats, nil
}

func GetNewsletterEmailSendsByNewsletter(
	ctx context.Context,
	dbtx db.DBTX,
	newsletterID uuid.UUID,
) ([]NewsletterEmailSend, error) {
	rows, err := db.Stmts.GetNewsletterEmailSendsByNewsletter(ctx, dbtx, newsletterID)
	if err != nil {
		return nil, err
	}

	sends := make([]NewsletterEmailSend, len(rows))
	for i, row := range rows {
		sends[i] = newsletterEmailSendRowToModel(row)
	}

	return sends, nil
}

func DeleteNewsletterEmailSends(
	ctx context.Context,
	dbtx db.DBTX,
	newsletterID uuid.UUID,
) error {
	return db.Stmts.DeleteNewsletterEmailSends(ctx, dbtx, newsletterID)
}

func parseNumericToFloat64(numeric pgtype.Numeric) float64 {
	if !numeric.Valid {
		return 0.0
	}
	
	return 0.0
}

func newsletterEmailSendRowToModel(row db.NewsletterEmailSend) NewsletterEmailSend {
	return NewsletterEmailSend{
		ID:           row.ID,
		NewsletterID: row.NewsletterID,
		SubscriberID: row.SubscriberID,
		EmailAddress: row.EmailAddress,
		Status:       row.Status,
		SentAt:       row.SentAt.Time,
		FailedAt:     row.FailedAt.Time,
		ErrorMessage: row.ErrorMessage.String,
		RiverJobID:   row.RiverJobID.Int64,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}