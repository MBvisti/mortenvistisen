package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Tags struct {
	db storage.Pool
}

func NewTags(db storage.Pool) Tags {
	return Tags{db}
}

func (t Tags) Index(etx echo.Context) error {
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

	tagsList, err := models.PaginateTags(
		etx.Request().Context(),
		t.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.TagIndex(tagsList.Tags))
}

func (t Tags) Show(etx echo.Context) error {
	tagID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	tag, err := models.FindTag(etx.Request().Context(), t.db.Conn(), tagID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.TagShow(tag))
}

func (t Tags) New(etx echo.Context) error {
	return render(etx, views.TagNew())
}

type CreateTagFormPayload struct {
	Title string `json:"title"`
}

func (t Tags) Create(etx echo.Context) error {
	var payload CreateTagFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateTagFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.CreateTagData{
		Title: payload.Title,
	}

	tag, err := models.CreateTag(
		etx.Request().Context(),
		t.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to create tag: %v", err)); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.TagNew.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Tag created successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.TagShow.URL(tag.ID))
}

func (t Tags) Edit(etx echo.Context) error {
	tagID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	tag, err := models.FindTag(etx.Request().Context(), t.db.Conn(), tagID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.TagEdit(tag))
}

type UpdateTagFormPayload struct {
	Title string `json:"title"`
}

func (t Tags) Update(etx echo.Context) error {
	tagID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	var payload UpdateTagFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateTagFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.UpdateTagData{
		ID:    tagID,
		Title: payload.Title,
	}

	tag, err := models.UpdateTag(
		etx.Request().Context(),
		t.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to update tag: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.TagEdit.URL(tagID),
		)
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Tag updated successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.TagShow.URL(tag.ID))
}

func (t Tags) Destroy(etx echo.Context) error {
	tagID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	err = models.DestroyTag(etx.Request().Context(), t.db.Conn(), tagID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to delete tag: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.TagIndex.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Tag destroyed successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.TagIndex.URL())
}
