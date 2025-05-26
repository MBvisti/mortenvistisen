package models

import (
	"context"
	"database/sql"
	"net"
	"strings"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oschwald/geoip2-golang"
)

type SiteSession struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Hostname     string
	Browser      string
	OS           string
	Device       string
	Screen       string
	Lang         string
	Country      string
	Subdivision1 string
	Subdivision2 string
	City         string
}

type NewSiteSessionData struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Hostname  string
	Browser   string
	OS        string
	Device    string
	Screen    string
	Lang      string
	IPAddress string
	Finger    string
}

func GOCSiteSession(
	ctx context.Context,
	dbtx db.DBTX,
	payload NewSiteSessionData,
) (SiteSession, error) {
	if siteSession, err := db.Stmts.QuerySiteSession(ctx, dbtx, payload.ID); err == nil {
		return SiteSession{
			ID:           siteSession.ID,
			CreatedAt:    siteSession.CreatedAt.Time,
			Hostname:     siteSession.Hostname.String,
			Browser:      siteSession.Browser.String,
			OS:           siteSession.Os.String,
			Device:       siteSession.Device.String,
			Screen:       siteSession.Screen.String,
			Lang:         siteSession.Lang.String,
			Country:      siteSession.Country.String,
			Subdivision1: siteSession.Subdivision1.String,
			Subdivision2: siteSession.Subdivision2.String,
			City:         siteSession.City.String,
		}, nil
	}

	geoDB, err := geoip2.Open("data/geolite_2_city.mmdb")
	if err != nil {
		return SiteSession{}, err
	}
	defer geoDB.Close()

	ip := net.ParseIP(payload.IPAddress)
	rec, err := geoDB.City(ip)
	if err != nil {
		return SiteSession{}, err
	}

	params := db.InsertSiteSessionParams{
		ID: payload.ID,
		CreatedAt: pgtype.Timestamptz{
			Time:  payload.CreatedAt,
			Valid: true,
		},
		Hostname: sql.NullString{
			String: payload.Hostname,
			Valid:  true,
		},
		Browser: sql.NullString{
			String: strings.ToLower(payload.Browser),
			Valid:  true,
		},
		Os: sql.NullString{
			String: strings.ToLower(payload.OS),
			Valid:  true,
		},
		Device: sql.NullString{
			String: payload.Device,
			Valid:  true,
		},
		Screen: sql.NullString{
			String: payload.Screen,
			Valid:  true,
		},
		Lang: sql.NullString{
			String: payload.Lang,
			Valid:  true,
		},
		Country: sql.NullString{
			String: rec.Country.IsoCode,
			Valid:  true,
		},
		Finger: sql.NullString{
			String: payload.Finger,
			Valid:  true,
		},
	}

	if rec.City.Names["en"] != "" {
		params.City = sql.NullString{
			String: rec.City.Names["en"],
			Valid:  true,
		}
	}

	switch len(rec.Subdivisions) {
	case 1:
		params.Subdivision1 = sql.NullString{
			String: rec.Subdivisions[0].Names["en"],
			Valid:  true,
		}
	case 2:
		params.Subdivision1 = sql.NullString{
			String: rec.Subdivisions[0].Names["en"],
			Valid:  true,
		}
		params.Subdivision2 = sql.NullString{
			String: rec.Subdivisions[1].Names["en"],
			Valid:  true,
		}

	}

	newSession, err := db.Stmts.InsertSiteSession(
		ctx,
		dbtx,
		params,
	)
	if err != nil {
		return SiteSession{}, err
	}

	return SiteSession{
		ID:           newSession.ID,
		CreatedAt:    newSession.CreatedAt.Time,
		Hostname:     newSession.Hostname.String,
		Browser:      newSession.Browser.String,
		OS:           newSession.Os.String,
		Device:       newSession.Device.String,
		Screen:       newSession.Screen.String,
		Lang:         newSession.Lang.String,
		Country:      newSession.Country.String,
		Subdivision1: newSession.Subdivision1.String,
		Subdivision2: newSession.Subdivision2.String,
		City:         newSession.City.String,
	}, nil
}
