package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Subscribers struct {
	db storage.Pool
}

func NewSubscribers(db storage.Pool) Subscribers {
	return Subscribers{db}
}

func (s Subscribers) Index(etx echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(25)
	if pp := etx.QueryParam("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 &&
			parsed <= 100 {
			perPage = int64(parsed)
		}
	}

	subscribersList, err := models.PaginateSubscribers(
		etx.Request().Context(),
		s.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.SubscriberIndex(subscribersList.Subscribers))
}

func (s Subscribers) Show(etx echo.Context) error {
	subscriberID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	subscriber, err := models.FindSubscriber(etx.Request().Context(), s.db.Conn(), subscriberID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.SubscriberShow(subscriber))
}

func (s Subscribers) New(etx echo.Context) error {
	return render(etx, views.SubscriberNew())
}

type CreateSubscriberFormPayload struct {
	Email        string `json:"email"`
	SubscribedAt string `json:"subscribed_at"`
	Referer      string `json:"referer"`
	IsVerified   bool   `json:"is_verified"`
}

func (s Subscribers) Create(etx echo.Context) error {
	var payload CreateSubscriberFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateSubscriberFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.CreateSubscriberData{
		Email: payload.Email,
		SubscribedAt: func() time.Time {
			if payload.SubscribedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.SubscribedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Referer:    payload.Referer,
		IsVerified: payload.IsVerified,
	}

	subscriber, err := models.CreateSubscriber(
		etx.Request().Context(),
		s.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to create subscriber: %v", err)); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.SubscriberNew.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Subscriber created successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.SubscriberShow.URL(subscriber.ID))
}

func (s Subscribers) Edit(etx echo.Context) error {
	subscriberID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	subscriber, err := models.FindSubscriber(etx.Request().Context(), s.db.Conn(), subscriberID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.SubscriberEdit(subscriber))
}

type UpdateSubscriberFormPayload struct {
	Email        string `json:"email"`
	SubscribedAt string `json:"subscribed_at"`
	Referer      string `json:"referer"`
	IsVerified   bool   `json:"is_verified"`
}

func (s Subscribers) Update(etx echo.Context) error {
	subscriberID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	var payload UpdateSubscriberFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateSubscriberFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.UpdateSubscriberData{
		ID:    subscriberID,
		Email: payload.Email,
		SubscribedAt: func() time.Time {
			if payload.SubscribedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.SubscribedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Referer:    payload.Referer,
		IsVerified: payload.IsVerified,
	}

	subscriber, err := models.UpdateSubscriber(
		etx.Request().Context(),
		s.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to update subscriber: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.SubscriberEdit.URL(subscriberID),
		)
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Subscriber updated successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.SubscriberShow.URL(subscriber.ID))
}

func (s Subscribers) Destroy(etx echo.Context) error {
	subscriberID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	err = models.DestroySubscriber(etx.Request().Context(), s.db.Conn(), subscriberID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to delete subscriber: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.SubscriberIndex.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Subscriber destroyed successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.SubscriberIndex.URL())
}
