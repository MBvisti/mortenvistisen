package admin

import (
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
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type Tags struct {
	db storage.Pool
}

func NewTags(db storage.Pool) Tags {
	return Tags{db}
}

func (t Tags) RegisterRoutes(r *router.Router) error {
	var errs []error
	for _, route := range []echo.Route{
		{
			Method:  http.MethodGet,
			Path:    routes.AdminTagIndex.Path(),
			Name:    routes.AdminTagIndex.Name(),
			Handler: t.Index,
		},
		{
			Method:  http.MethodPost,
			Path:    routes.AdminTagCreate.Path(),
			Name:    routes.AdminTagCreate.Name(),
			Handler: t.Create,
		},
		{
			Method:  http.MethodPut,
			Path:    routes.AdminTagUpdate.Path(),
			Name:    routes.AdminTagUpdate.Name(),
			Handler: t.Update,
		},
		{
			Method:  http.MethodDelete,
			Path:    routes.AdminTagDestroy.Path(),
			Name:    routes.AdminTagDestroy.Name(),
			Handler: t.Destroy,
		},
	} {
		if _, err := r.AddRoute(route); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

type TagData struct {
	ID           int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Title        string
	ArticleCount int64
}

func newTagData(entity models.TagEntity) TagData {
	return TagData{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Title:     entity.Title,
	}
}

func newTagDataList(entities []models.TagEntity) []TagData {
	items := make([]TagData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newTagData(entity))
	}

	return items
}

func newTagUsageDataList(entities []models.TagWithUsage) []TagData {
	items := make([]TagData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, TagData{
			ID:           entity.ID,
			CreatedAt:    entity.CreatedAt,
			UpdatedAt:    entity.UpdatedAt,
			Title:        entity.Title,
			ArticleCount: entity.ArticleCount,
		})
	}

	return items
}

type TagPaginationData struct {
	Page       int64
	PageSize   int64
	TotalCount int64
	TotalPages int64
}

const tagPageSize int64 = 10

func tagIndexURL(etx *echo.Context) string {
	if value := etx.QueryParam("page"); value != "" {
		if page, err := strconv.Atoi(value); err == nil && page > 1 {
			return fmt.Sprintf("%s?page=%d", routes.AdminTagIndex.URL(), page)
		}
	}

	return routes.AdminTagIndex.URL()
}

func (t Tags) renderIndex(etx *echo.Context, opts ...inertia.PageOption) error {
	page := int64(1)
	if value := etx.QueryParam("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	tags, err := models.Tag.PaginateWithUsage(
		etx.Request().Context(),
		t.db.Executor(),
		page,
		tagPageSize,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Tag/Index", inertia.Props{
		"items": newTagUsageDataList(tags.Tags),
		"pagination": TagPaginationData{
			Page:       tags.Page,
			PageSize:   tags.PageSize,
			TotalCount: tags.TotalCount,
			TotalPages: tags.TotalPages,
		},
		"validationRules": models.TagValidationRules(),
	}, opts...)
}

func (t Tags) Index(etx *echo.Context) error {
	return t.renderIndex(etx)
}

type TagFormPayload struct {
	Title string `json:"title"`
}

func (t Tags) Create(etx *echo.Context) error {
	var payload TagFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(etx.Request().Context(), "could not parse tag payload", "error", err)
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	_, err := models.Tag.Create(
		etx.Request().Context(),
		t.db.Executor(),
		models.CreateTagData{Title: payload.Title},
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			return t.renderIndex(etx, inertia.WithValidationErrors(validationErrors.ToMap()))
		}
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create tag: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, tagIndexURL(etx))
	}

	if err := cookies.AddFlash(etx, cookies.FlashSuccess, "Tag created successfully"); err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, tagIndexURL(etx))
}

func (t Tags) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	var payload TagFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(etx.Request().Context(), "could not parse tag payload", "error", err)
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	_, err = models.Tag.Update(
		etx.Request().Context(),
		t.db.Executor(),
		models.UpdateTagData{ID: int32(parsed), Title: payload.Title},
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			return t.renderIndex(etx, inertia.WithValidationErrors(validationErrors.ToMap()))
		}
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update tag: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, tagIndexURL(etx))
	}

	if err := cookies.AddFlash(etx, cookies.FlashSuccess, "Tag updated successfully"); err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, tagIndexURL(etx))
}

func (t Tags) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	if err := models.Tag.Destroy(
		etx.Request().Context(),
		t.db.Executor(),
		int32(parsed),
	); err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete tag: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, tagIndexURL(etx))
	}

	if err := cookies.AddFlash(etx, cookies.FlashSuccess, "Tag deleted successfully"); err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, tagIndexURL(etx))
}
