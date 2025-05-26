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

type EventType string

var (
	PageView  EventType = "page_view"
	PageLeave EventType = "page_leave"
	Click     EventType = "click"
)

func GetEventType(name string) EventType {
	switch name {
	case "page_view":
		return PageView
	case "page_leave":
		return PageLeave
	case "click":
		return Click
	}

	// TODO: fix
	return "SHIT"
}

type SiteView struct {
	ID             int64
	CreatedAt      time.Time
	UrlPath        string
	UrlQuery       string
	ReferrerPath   string
	ReferrerQuery  string
	ReferrerDomain string
	PageTitle      string
	EventType      string
	EventName      string
	VisitorID      uuid.UUID
	Session        SiteSession
}

type NewSiteViewData struct {
	SessionID      uuid.UUID
	VisitorID      uuid.UUID
	CreatedAt      time.Time
	UrlPath        string
	UrlQuery       string
	ReferrerPath   string
	ReferrerQuery  string
	ReferrerDomain string
	PageTitle      string
	EventType      EventType
	EventName      string
}

func NewSiteView(
	ctx context.Context,
	dbtx db.DBTX,
	payload NewSiteViewData,
) error {
	params := db.InsertSiteViewParams{
		SessionID: pgtype.UUID{
			Bytes: payload.SessionID,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  payload.CreatedAt,
			Valid: true,
		},
		UrlPath: sql.NullString{
			String: payload.UrlPath,
			Valid:  true,
		},
		UrlQuery: sql.NullString{
			String: payload.UrlQuery,
			Valid:  true,
		},
		ReferrerPath: sql.NullString{
			String: payload.ReferrerPath,
			Valid:  true,
		},
		ReferrerQuery: sql.NullString{
			String: payload.ReferrerQuery,
			Valid:  true,
		},
		ReferrerDomain: sql.NullString{
			String: payload.ReferrerDomain,
			Valid:  true,
		},
		PageTitle: sql.NullString{
			String: payload.PageTitle,
			Valid:  true,
		},
		EventName: sql.NullString{
			String: payload.EventName,
			Valid:  true,
		},
		VisitorID: pgtype.UUID{
			Bytes: payload.VisitorID,
			Valid: true,
		},
	}

	switch payload.EventType {
	case PageView:
		params.EventType = db.NullSiteEvent{
			SiteEvent: db.SiteEventPageView,
			Valid:     true,
		}
	case PageLeave:
		params.EventType = db.NullSiteEvent{
			SiteEvent: db.SiteEventPageLeave,
			Valid:     true,
		}
	case Click:
		params.EventType = db.NullSiteEvent{
			SiteEvent: db.SiteEventClick,
			Valid:     true,
		}
	}

	return db.Stmts.InsertSiteView(ctx, dbtx, params)
}

func GetSiteViewsByDate(
	ctx context.Context,
	dbtx db.DBTX,
	start time.Time,
	end time.Time,
) ([]SiteView, error) {
	rows, err := db.Stmts.QueryViewsByDate(
		ctx,
		dbtx,
		db.QueryViewsByDateParams{
			StartDate: pgtype.Timestamp{
				Time:  start,
				Valid: true,
			},
			EndDate: pgtype.Timestamp{
				Time:  end,
				Valid: true,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	siteViews := make([]SiteView, len(rows))
	for i, row := range rows {
		siteViews[i] = SiteView{
			ID:             row.SiteView.ID,
			CreatedAt:      row.SiteView.CreatedAt.Time,
			UrlPath:        row.SiteView.UrlPath.String,
			UrlQuery:       row.SiteView.UrlQuery.String,
			ReferrerPath:   row.SiteView.ReferrerPath.String,
			ReferrerQuery:  row.SiteView.ReferrerQuery.String,
			ReferrerDomain: row.SiteView.ReferrerDomain.String,
			PageTitle:      row.SiteView.PageTitle.String,
			EventType:      string(row.SiteView.EventType.SiteEvent),
			EventName:      row.SiteView.EventName.String,
			VisitorID:      row.SiteView.VisitorID.Bytes,
			Session: SiteSession{
				ID:           row.SiteSession.ID,
				CreatedAt:    row.SiteSession.CreatedAt.Time,
				Hostname:     row.SiteSession.Hostname.String,
				Browser:      row.SiteSession.Browser.String,
				OS:           row.SiteSession.Os.String,
				Device:       row.SiteSession.Device.String,
				Screen:       row.SiteSession.Screen.String,
				Lang:         row.SiteSession.Lang.String,
				Country:      row.SiteSession.Country.String,
				Subdivision1: row.SiteSession.Subdivision1.String,
				Subdivision2: row.SiteSession.Subdivision2.String,
				City:         row.SiteSession.City.String,
			},
		}
	}

	return siteViews, nil
}

type TrafficCount struct {
	Visits int64
	Views  int64
}

func GetHourlyTrafficCounts(
	ctx context.Context,
	dbtx db.DBTX,
) (map[time.Time]TrafficCount, error) {
	end := carbon.Now().EndOfHour()
	// start := now.StartOfDay()
	start := end.SubHours(24)

	rows, err := db.Stmts.QueryTrafficCountsByDate(
		ctx,
		dbtx,
		db.QueryTrafficCountsByDateParams{
			StartDate: pgtype.Timestamp{
				Time:  start.StdTime(),
				Valid: true,
			},
			EndDate: pgtype.Timestamp{
				Time:  end.StdTime(),
				Valid: true,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	hourlyVisitCounts := make(map[time.Time]TrafficCount, len(rows))
	for _, r := range rows {
		hourlyVisitCounts[r.Hour.Time] = TrafficCount{
			Visits: r.VisitorCount,
			Views:  r.Views,
		}
	}

	return hourlyVisitCounts, nil
}
