package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Api struct {
	db psql.Postgres
}

func newApi(
	db psql.Postgres,
) Api {
	return Api{db}
}

func (a *Api) AppHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "app is healthy and running")
}

func (a *Api) Collect(c echo.Context) error {
	type collectPayload struct {
		WebsiteID    uuid.UUID `json:"website_id"`
		Type         string    `json:"type"`
		URL          string    `json:"url"`
		Path         string    `json:"path"`
		Referrer     string    `json:"referrer"`
		Title        string    `json:"title"`
		Timestamp    time.Time `json:"timestamp"`
		Screen       string    `json:"screen"`
		Language     string    `json:"language"`
		VisitorID    uuid.UUID `json:"visitor_id"`
		SessionID    uuid.UUID `json:"session_id"`
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

	id, err := models.NewAnalytic(
		c.Request().Context(),
		models.NewAnalyticPayload{
			Timestamp:   payload.Timestamp,
			WebsiteID:   payload.WebsiteID,
			Type:        payload.Type,
			URL:         payload.URL,
			Path:        payload.Path,
			Referrer:    payload.Referrer,
			Title:       payload.Title,
			Screen:      payload.Screen,
			Language:    payload.Language,
			VisitorID:   payload.VisitorID,
			SessionID:   payload.SessionID,
			ScrollDepth: payload.ScrollDepth,
		},
		a.db.Pool,
	)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Collect", "error", err)
		return c.JSON(http.StatusOK, "could not collect analytics")
	}

	if payload.CustomData != "" {
		if err := models.NewAnalyticEvent(c.Request().Context(), models.NewAnalyticEventPayload{
			AnalyticsID:  id,
			ElementTag:   payload.ElementTag,
			ElementText:  payload.ElementText,
			ElementClass: payload.ElementClass,
			ElementID:    payload.ElementID,
			ElementHref:  payload.ElementHref,
			CustomData:   payload.CustomData,
		}, a.db.Pool); err != nil {
			slog.ErrorContext(c.Request().Context(), "Collect", "error", err)
			return c.JSON(http.StatusOK, "could not collect analytics")
		}
	}

	return c.JSON(http.StatusOK, "collected")
}
