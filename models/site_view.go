package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/db"
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
