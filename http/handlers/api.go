package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
	"github.com/mileusna/useragent"
)

type Api struct {
	db    psql.Postgres
	cache otter.Cache[string, Session]
}

type Session struct {
	SessionID uuid.UUID
	VisitorID uuid.UUID
}

func newApi(
	db psql.Postgres,
) Api {
	cache, err := otter.MustBuilder[string, Session](10_000).
		WithTTL(30 * time.Minute).
		Build()
	if err != nil {
		panic(err)
	}

	return Api{db, cache}
}

func (a *Api) AppHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "app is healthy and running")
}

func generateFingerprint(ip, userAgent string) string {
	hash := sha256.Sum256([]byte(ip + userAgent))
	return hex.EncodeToString(hash[:])
}

func (a *Api) Collect(c echo.Context) error {
	type collectPayload struct {
		Type         string    `json:"type"`
		URL          string    `json:"url"`
		Path         string    `json:"path"`
		Referrer     string    `json:"referrer"`
		Title        string    `json:"title"`
		Timestamp    time.Time `json:"timestamp"`
		Screen       string    `json:"screen"`
		Language     string    `json:"language"`
		UserAgent    string    `json:"user_agent"`
		ScrollDepth  int32     `json:"scroll_depth"`
		ElementTag   string    `json:"element_tag"`
		ElementText  string    `json:"element_text"`
		ElementClass string    `json:"element_class"`
		ElementID    string    `json:"element_id"`
		ElementHref  string    `json:"element_href"`
		CustomData   string    `json:"custom_data"`
	}

	var payload collectPayload
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return c.JSON(http.StatusOK, "could not collect analytics")
	}

	var sessID uuid.UUID
	var visID uuid.UUID

	finger := generateFingerprint(c.RealIP(), payload.UserAgent)
	cachedSession, exists := a.cache.Get(finger)
	if exists {
		sessID = cachedSession.SessionID
		visID = cachedSession.VisitorID
	}

	if !exists {
		sessID = uuid.New()
		visID = uuid.New()

		if ok := a.cache.Set(finger, Session{
			SessionID: sessID,
			VisitorID: visID,
		}); !ok {
			return c.JSON(http.StatusOK, "could not collect analytics")
		}
	}

	ua := useragent.Parse(payload.UserAgent)

	if ua.Bot {
		return c.JSON(http.StatusOK, "bot visit")
	}

	device := "unknown"
	if ua.Desktop {
		device = "desktop"
	}
	if ua.Mobile {
		device = "mobile"
	}
	if ua.Tablet {
		device = "table"
	}

	session, err := models.GOCSiteSession(
		c.Request().Context(),
		a.db.Pool,
		models.NewSiteSessionData{
			ID:        sessID,
			CreatedAt: payload.Timestamp,
			Hostname:  c.Request().Host,
			Browser:   ua.Name,
			OS:        ua.OS,
			Device:    device,
			Screen:    payload.Screen,
			Lang:      payload.Language,
			IPAddress: c.RealIP(),
			Finger:    finger,
		},
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"COLLECT HANDLER",
			"error",
			err,
		)
		return c.JSON(http.StatusOK, "could not collect analytics")
	}

	parsedURL, err := url.Parse(payload.URL)
	if err != nil {
		return c.JSON(http.StatusOK, "could not collect analytics")
	}

	parsedReferrerURL, err := url.Parse(payload.Referrer)
	if err != nil {
		return c.JSON(http.StatusOK, "could not collect analytics")
	}

	if err := models.NewSiteView(c.Request().Context(), a.db.Pool, models.NewSiteViewData{
		SessionID:      session.ID,
		CreatedAt:      payload.Timestamp,
		UrlPath:        payload.URL,
		VisitorID:      visID,
		UrlQuery:       parsedURL.RawQuery,
		ReferrerPath:   payload.Referrer,
		ReferrerQuery:  parsedReferrerURL.RawQuery,
		ReferrerDomain: parsedReferrerURL.Host,
		PageTitle:      payload.Title,
		EventType:      models.GetEventType(payload.Type),
		EventName:      payload.CustomData,
	}); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"COLLECT HANDLER",
			"error",
			err,
		)
		return c.JSON(http.StatusOK, "could not collect analytics")
	}

	return c.JSON(http.StatusOK, "")
}
