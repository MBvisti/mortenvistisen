package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/services"
	"mortenvistisen/views"

	"github.com/labstack/echo/v5"
)

type Subscribers struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
	cfg        config.Config
}

func NewSubscribers(
	db storage.Pool,
	insertOnly queue.InsertOnly,
	cfg config.Config,
) Subscribers {
	return Subscribers{
		db:         db,
		insertOnly: insertOnly,
		cfg:        cfg,
	}
}

func (s Subscribers) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(10)
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

	return render(etx, views.SubscriberIndex(subscribersList))
}

func (s Subscribers) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	subscriberID := int32(parsed)

	subscriber, err := models.FindSubscriber(etx.Request().Context(), s.db.Conn(), subscriberID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.SubscriberShow(subscriber))
}

func (s Subscribers) New(etx *echo.Context) error {
	return render(etx, views.SubscriberNew())
}

type CreateSubscriberFormPayload struct {
	Email        string `json:"email"`
	SubscribedAt string `json:"subscribedAt"`
	Referer      string `json:"referer"`
	IsVerified   bool   `json:"isVerified"`
}

func (s Subscribers) Signup(etx *echo.Context) error {
	var payload struct {
		Email string `json:"email"`
	}

	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse newsletter signup payload",
			"error",
			err,
		)

		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Could not process newsletter signup"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.HomePage.URL())
	}

	slog.Info("PAYLOAD", "email", payload.Email)

	err := services.RequestSubscriberVerification(
		etx.Request().Context(),
		s.db,
		s.insertOnly,
		s.cfg.Auth.Pepper,
		services.RequestSubscriberVerificationData{
			Email:   payload.Email,
			Referer: etx.Request().Referer(),
		},
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to request subscriber verification",
			"error",
			err,
		)

		if errors.Is(err, services.ErrSubscriberAlreadyVerified) {
			if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Could not sign up for the newsletter"); flashErr != nil {
				return render(etx, views.InternalError())
			}
			return etx.Redirect(http.StatusSeeOther, routes.HomePage.URL())
		}

		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Could not sign up for the newsletter"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.HomePage.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Check your email for the verification code"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.HomePage.URL())
}

func (s Subscribers) VerificationNew(etx *echo.Context) error {
	return render(etx, views.SubscriberConfirmationForm())
}

func (s Subscribers) VerificationCreate(etx *echo.Context) error {
	var payload struct {
		Code string `json:"code"`
	}

	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse subscriber verification payload",
			"error",
			err,
		)
		return render(etx, views.BadRequest())
	}

	_, err := services.VerifySubscriber(
		etx.Request().Context(),
		s.db,
		s.cfg.Auth.Pepper,
		services.VerifySubscriberData{
			Code: strings.ToUpper(strings.TrimSpace(payload.Code)),
		},
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to verify subscriber",
			"error",
			err,
		)

		var errorMsg string
		switch err {
		case services.ErrSubscriberVerificationInvalidCode:
			errorMsg = "Invalid verification code"
		case services.ErrSubscriberVerificationExpiredCode:
			errorMsg = "Verification code has expired"
		default:
			errorMsg = "Failed to verify subscription"
		}

		if flashErr := cookies.AddFlash(etx, cookies.FlashError, errorMsg); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.SubscriberVerificationNew.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Newsletter subscription verified"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.HomePage.URL())
}

func (s Subscribers) Create(etx *echo.Context) error {
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

func (s Subscribers) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	subscriberID := int32(parsed)

	subscriber, err := models.FindSubscriber(etx.Request().Context(), s.db.Conn(), subscriberID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.SubscriberEdit(subscriber))
}

type UpdateSubscriberFormPayload struct {
	Email        string `json:"email"`
	SubscribedAt string `json:"subscribedAt"`
	Referer      string `json:"referer"`
	IsVerified   bool   `json:"isVerified"`
}

func (s Subscribers) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	subscriberID := int32(parsed)

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

func (s Subscribers) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	subscriberID := int32(parsed)

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
