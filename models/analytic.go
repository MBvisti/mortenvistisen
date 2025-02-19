package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Analytic struct {
	ID          uuid.UUID
	WebsiteID   uuid.UUID
	Type        string
	URL         string
	Path        string
	Referrer    string
	Title       string
	Timestamp   time.Time
	Screen      string
	Language    string
	VisitorID   uuid.UUID
	SessionID   uuid.UUID
	ScrollDepth int64
}

type AnalyticEvent struct {
	ID           uuid.UUID
	WebsiteID    uuid.UUID
	Timestamp    time.Time
	VisitorID    uuid.UUID
	SessionID    uuid.UUID
	ElementTag   string
	ElementText  string
	ElementClass string
	ElementID    string
	ElementHref  string
	CustomData   string
}

// func GetAnalytic(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	dbtx db.DBTX,
// ) (Analytic, error) {
// 	a, err := db.Stmts.QueryAnalyticsByID(ctx, dbtx, id)
// 	if err != nil {
// 		return Analytic{}, err
// 	}
//
// 	return Analytic{
// 		ID:          a.ID,
// 		CreatedAt:   a.CreatedAt.Time,
// 		UpdatedAt:   a.UpdatedAt.Time,
// 		WebsiteID:   a.WebsiteID.Bytes,
// 		Type:        a.Type.String,
// 		URL:         a.Url.String,
// 		Path:        a.Path.String,
// 		Referrer:    a.Referrer.String,
// 		Title:       a.Title.String,
// 		Timestamp:   a.Timestamp.Time,
// 		Screen:      a.Screen.String,
// 		Language:    a.Language.String,
// 		VisitorID:   a.VisitorID.Bytes,
// 		SessionID:   a.SessionID.Bytes,
// 		ScrollDepth: stringToInt64(a.ScrollDepth.String),
// 	}, nil
// }

type NewAnalyticPayload struct {
	WebsiteID   uuid.UUID
	Type        string
	URL         string
	Path        string
	Referrer    string
	Title       string
	Timestamp   time.Time
	Screen      string
	Language    string
	VisitorID   uuid.UUID
	SessionID   uuid.UUID
	ScrollDepth int32
	RealIP      string
}

func NewAnalytic(
	ctx context.Context,
	data NewAnalyticPayload,
	dbtx db.DBTX,
) (uuid.UUID, error) {
	id := uuid.New()
	if err := db.Stmts.InsertAnalytic(ctx, dbtx, db.InsertAnalyticParams{
		ID: id,
		WebsiteID: pgtype.UUID{
			Bytes: data.WebsiteID,
			Valid: true,
		},
		Type: sql.NullString{
			String: data.Type,
			Valid:  true,
		},
		Url: sql.NullString{
			String: data.URL,
			Valid:  true,
		},
		Path: sql.NullString{
			String: data.Path,
			Valid:  true,
		},
		Referrer: sql.NullString{
			String: data.Referrer,
			Valid:  true,
		},
		Title: sql.NullString{
			String: data.Title,
			Valid:  true,
		},
		Timestamp: pgtype.Timestamptz{
			Time:  data.Timestamp,
			Valid: true,
		},
		Screen: sql.NullString{
			String: data.Screen,
			Valid:  true,
		},
		Language: sql.NullString{
			String: data.Language,
			Valid:  true,
		},
		VisitorID: pgtype.UUID{Bytes: data.VisitorID, Valid: true},
		SessionID: pgtype.UUID{Bytes: data.SessionID, Valid: true},
		ScrollDepth: sql.NullInt32{
			Int32: data.ScrollDepth,
			Valid: true,
		},
		RealIp: sql.NullString{
			String: data.RealIP,
			Valid:  true,
		},
	}); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

type NewAnalyticEventPayload struct {
	AnalyticsID  uuid.UUID
	ElementTag   string
	ElementText  string
	ElementClass string
	ElementID    string
	ElementHref  string
	CustomData   string
}

func NewAnalyticEvent(
	ctx context.Context,
	data NewAnalyticEventPayload,
	dbtx db.DBTX,
) error {
	return db.Stmts.InsertAnalyticEvent(ctx, dbtx, db.InsertAnalyticEventParams{
		ID: uuid.New(),
		AnalyticsID: pgtype.UUID{
			Bytes: data.AnalyticsID,
			Valid: true,
		},
		ElementTag: sql.NullString{
			String: data.ElementTag,
			Valid:  true,
		},
		ElementText: sql.NullString{
			String: data.ElementText,
			Valid:  true,
		},
		ElementClass: sql.NullString{
			String: data.ElementClass,
			Valid:  true,
		},
		ElementID: sql.NullString{
			String: data.ElementID,
			Valid:  true,
		},
		ElementHref: sql.NullString{
			String: data.ElementHref,
			Valid:  true,
		},
		CustomData: sql.NullString{
			String: data.CustomData,
			Valid:  true,
		},
	})
}

func GetDailyVisits(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	visits, err := db.Stmts.QueryDailyVisits(
		ctx,
		dbtx,
		db.QueryDailyVisitsParams{
			WebsiteID: pgtype.UUID{
				Bytes: uuid.MustParse("0210debc-df55-4c4c-a2f2-01c268a04911"),
				Valid: true,
			},
			Date: time.Now(),
		},
	)
	if err != nil {
		return 0, err
	}

	return visits, nil
}

func GetDailyViews(
	ctx context.Context,
	dbtx db.DBTX,
) (int64, error) {
	views, err := db.Stmts.QueryDailyViews(
		ctx,
		dbtx,
		db.QueryDailyViewsParams{
			WebsiteID: pgtype.UUID{
				Bytes: uuid.MustParse("0210debc-df55-4c4c-a2f2-01c268a04911"),
				Valid: true,
			},
			Date: time.Now(),
		},
	)
	if err != nil {
		return 0, err
	}

	return views, nil
}

type Stat struct {
	Hour   string
	Views  int64
	Visits int64
}

func GetDailyStats(
	ctx context.Context,
	dbtx db.DBTX,
) ([]Stat, error) {
	statRows, err := db.Stmts.QueryHourlyStats(
		ctx,
		dbtx,
	)
	if err != nil {
		return nil, err
	}

	stats := make([]Stat, len(statRows))
	for i, s := range statRows {
		h := carbon.CreateFromStdTime(s.Hour.Time).
			ToKitchenString("Europe/Berlin")

		stats[i] = Stat{
			Hour:   h,
			Views:  s.Views,
			Visits: s.Visits,
		}
	}

	return stats, nil
}

// func DeleteAnalytic(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	dbtx db.DBTX,
// ) error {
// 	return db.Stmts.DeleteAnalytics(ctx, dbtx, id)
// }
//
// func GetAnalyticsByWebsite(
// 	ctx context.Context,
// 	websiteID uuid.UUID,
// 	dbtx db.DBTX,
// ) ([]Analytic, error) {
// 	analytics, err := db.Stmts.QueryAnalyticsByWebsiteID(
// 		ctx,
// 		dbtx,
// 		pgtype.UUID{Bytes: websiteID, Valid: true},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	result := make([]Analytic, len(analytics))
// 	for i, a := range analytics {
// 		result[i] = Analytic{
// 			ID:          a.ID,
// 			CreatedAt:   a.CreatedAt.Time,
// 			UpdatedAt:   a.UpdatedAt.Time,
// 			WebsiteID:   a.WebsiteID.Bytes,
// 			Type:        a.Type.String,
// 			URL:         a.Url.String,
// 			Path:        a.Path.String,
// 			Referrer:    a.Referrer.String,
// 			Title:       a.Title.String,
// 			Timestamp:   a.Timestamp.Time,
// 			Screen:      a.Screen.String,
// 			Language:    a.Language.String,
// 			VisitorID:   a.VisitorID.Bytes,
// 			SessionID:   a.SessionID.Bytes,
// 			ScrollDepth: stringToInt64(a.ScrollDepth.String),
// 		}
// 	}
//
// 	return result, nil
// }
//
// func nullString(s string) sql.NullString {
// 	return sql.NullString{String: s, Valid: s != ""}
// }
//
// func stringToInt64(s string) int64 {
// 	if s == "" {
// 		return 0
// 	}
// 	i, err := strconv.ParseInt(s, 10, 64)
// 	if err != nil {
// 		return 0
// 	}
// 	return i
// }
