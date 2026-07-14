package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"mortenvistisen/internal/inertia"
	"mortenvistisen/internal/storage"
	"mortenvistisen/internal/validation"
	"mortenvistisen/models"
	"mortenvistisen/router"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/services"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type Subscribers struct {
	db                     storage.Pool
	subscriberConfirmation services.SubscriberConfirmation
}

func NewSubscribers(
	db storage.Pool,
	subscriberConfirmation services.SubscriberConfirmation,
) Subscribers {
	return Subscribers{
		db:                     db,
		subscriberConfirmation: subscriberConfirmation,
	}
}

func (s Subscribers) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminSubscriberIndex.Path(),
		Name:    routes.AdminSubscriberIndex.Name(),
		Handler: s.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.AdminSubscriberConfirmationCreate.Path(),
		Name:    routes.AdminSubscriberConfirmationCreate.Name(),
		Handler: s.SendConfirmation,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodDelete,
		Path:    routes.AdminSubscriberBulkDestroy.Path(),
		Name:    routes.AdminSubscriberBulkDestroy.Name(),
		Handler: s.BulkDestroy,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminSubscriberShow.Path(),
		Name:    routes.AdminSubscriberShow.Name(),
		Handler: s.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminSubscriberNew.Path(),
		Name:    routes.AdminSubscriberNew.Name(),
		Handler: s.New,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.AdminSubscriberCreate.Path(),
		Name:    routes.AdminSubscriberCreate.Name(),
		Handler: s.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminSubscriberEdit.Path(),
		Name:    routes.AdminSubscriberEdit.Name(),
		Handler: s.Edit,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPut,
		Path:    routes.AdminSubscriberUpdate.Path(),
		Name:    routes.AdminSubscriberUpdate.Name(),
		Handler: s.Update,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodDelete,
		Path:    routes.AdminSubscriberDestroy.Path(),
		Name:    routes.AdminSubscriberDestroy.Name(),
		Handler: s.Destroy,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type SubscriberData struct {
	ID           int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string
	SubscribedAt time.Time
	Referer      string
	IsVerified   bool
}

func newSubscriberData(entity models.SubscriberEntity) SubscriberData {
	return SubscriberData{
		ID:           entity.ID,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		Email:        entity.Email.String,
		SubscribedAt: entity.SubscribedAt.Time,
		Referer:      entity.Referer.String,
		IsVerified:   entity.IsVerified.Bool,
	}
}

func newSubscriberDataList(entities []models.SubscriberEntity) []SubscriberData {
	items := make([]SubscriberData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newSubscriberData(entity))
	}

	return items
}

type SubscriberPaginationData struct {
	Page       int64
	PageSize   int64
	TotalCount int64
	TotalPages int64
}

const subscriberPageSize int64 = 10

func (s Subscribers) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	subscribersList, err := models.Subscriber.Paginate(
		etx.Request().Context(),
		s.db.Executor(),
		page,
		subscriberPageSize,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Subscriber/Index", inertia.Props{
		"items": newSubscriberDataList(subscribersList.Subscribers),
		"pagination": SubscriberPaginationData{
			Page:       subscribersList.Page,
			PageSize:   subscribersList.PageSize,
			TotalCount: subscribersList.TotalCount,
			TotalPages: subscribersList.TotalPages,
		},
	})
}

func (s Subscribers) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	subscriberID := int32(parsed)

	subscriber, err := models.Subscriber.Find(
		etx.Request().Context(),
		s.db.Executor(),
		subscriberID,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Subscriber/Show", inertia.Props{
		"item": newSubscriberData(subscriber),
	})
}

func (s Subscribers) SendConfirmation(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	subscriberID := int32(parsed)

	if err := s.subscriberConfirmation.Send(etx.Request().Context(), subscriberID); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to send subscriber confirmation email",
			"subscriber_id", subscriberID,
			"error", err,
		)

		message := "Failed to send confirmation email"
		if errors.Is(err, services.ErrSubscriberAlreadyVerified) {
			message = "This subscriber is already verified"
		}
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, message); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminSubscriberShow.URL(subscriberID),
			http.StatusSeeOther,
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Confirmation email sent successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(
		etx,
		routes.AdminSubscriberShow.URL(subscriberID),
		http.StatusSeeOther,
	)
}

func (s Subscribers) New(etx *echo.Context) error {
	return inertia.Page(etx, "Admin/Subscriber/Create", inertia.Props{
		"validationRules": models.SubscriberValidationRules(),
	})
}

type CreateSubscriberFormPayload struct {
	Email        string `json:"email"`
	SubscribedAt string `json:"subscribedAt"`
	Referer      string `json:"referer"`
	IsVerified   bool   `json:"isVerified"`
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

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	data := models.CreateSubscriberData{
		Email: sql.NullString{String: payload.Email, Valid: true},
		SubscribedAt: func() sql.NullTime {
			if payload.SubscribedAt == "" {
				return sql.NullTime{Valid: false}
			}
			if t, err := time.Parse("2006-01-02", payload.SubscribedAt); err == nil {
				return sql.NullTime{Time: t, Valid: true}
			}
			return sql.NullTime{Valid: false}
		}(),
		Referer:    sql.NullString{String: payload.Referer, Valid: true},
		IsVerified: sql.NullBool{Bool: payload.IsVerified, Valid: true},
	}

	subscriber, err := models.Subscriber.Create(
		etx.Request().Context(),
		s.db.Executor(),
		data,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/Subscriber/Create",
				inertia.Props{
					"validationRules": models.SubscriberValidationRules(),
				},
				inertia.WithValidationErrors(validationErrors.ToMap()),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create subscriber: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminSubscriberNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Subscriber created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminSubscriberShow.URL(subscriber.ID))
}

func (s Subscribers) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	subscriberID := int32(parsed)

	subscriber, err := models.Subscriber.Find(
		etx.Request().Context(),
		s.db.Executor(),
		subscriberID,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Subscriber/Edit", inertia.Props{
		"item":            newSubscriberData(subscriber),
		"validationRules": models.SubscriberValidationRules(),
	})
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
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
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

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	data := models.UpdateSubscriberData{
		ID:    subscriberID,
		Email: sql.NullString{String: payload.Email, Valid: true},
		SubscribedAt: func() sql.NullTime {
			if payload.SubscribedAt == "" {
				return sql.NullTime{Valid: false}
			}
			if t, err := time.Parse("2006-01-02", payload.SubscribedAt); err == nil {
				return sql.NullTime{Time: t, Valid: true}
			}
			return sql.NullTime{Valid: false}
		}(),
		Referer:    sql.NullString{String: payload.Referer, Valid: true},
		IsVerified: sql.NullBool{Bool: payload.IsVerified, Valid: true},
	}

	subscriber, err := models.Subscriber.Update(
		etx.Request().Context(),
		s.db.Executor(),
		data,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			currentSubscriber, findErr := models.Subscriber.Find(
				etx.Request().Context(),
				s.db.Executor(),
				subscriberID,
			)
			if findErr != nil {
				return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
			}

			return inertia.Page(
				etx,
				"Admin/Subscriber/Edit",
				inertia.Props{
					"item":            newSubscriberData(currentSubscriber),
					"validationRules": models.SubscriberValidationRules(),
				},
				inertia.WithValidationErrors(validationErrors.ToMap()),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update subscriber: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminSubscriberEdit.URL(subscriberID),
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Subscriber updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminSubscriberShow.URL(subscriber.ID))
}

func (s Subscribers) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	subscriberID := int32(parsed)

	err = models.Subscriber.Destroy(etx.Request().Context(), s.db.Executor(), subscriberID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete subscriber: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminSubscriberIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Subscriber deleted successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminSubscriberIndex.URL())
}

type BulkDestroySubscribersPayload struct {
	SubscriberIDs []int32 `json:"subscriberIds"`
}

func (s Subscribers) BulkDestroy(etx *echo.Context) error {
	var payload BulkDestroySubscribersPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse BulkDestroySubscribersPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	if len(payload.SubscriberIDs) == 0 || len(payload.SubscriberIDs) > int(subscriberPageSize) {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			"Select between 1 and 10 subscribers to delete",
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}

		return inertia.Redirect(etx, routes.AdminSubscriberIndex.URL())
	}

	deletedCount, err := models.Subscriber.DestroyMany(
		etx.Request().Context(),
		s.db.Executor(),
		payload.SubscriberIDs,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete subscribers: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminSubscriberIndex.URL())
	}

	message := fmt.Sprintf("%d subscribers deleted successfully", deletedCount)
	if deletedCount == 1 {
		message = "1 subscriber deleted successfully"
	}
	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, message); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Redirect(etx, routes.AdminSubscriberIndex.URL())
}
