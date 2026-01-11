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

type Newsletters struct {
	db storage.Pool
}

func NewNewsletters(db storage.Pool) Newsletters {
	return Newsletters{db}
}

func (n Newsletters) Index(etx echo.Context) error {
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

	newslettersList, err := models.PaginateNewsletters(
		etx.Request().Context(),
		n.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.NewsletterIndex(newslettersList.Newsletters))
}

func (n Newsletters) Show(etx echo.Context) error {
	newsletterID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	newsletter, err := models.FindNewsletter(etx.Request().Context(), n.db.Conn(), newsletterID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.NewsletterShow(newsletter))
}

func (n Newsletters) New(etx echo.Context) error {
	return render(etx, views.NewsletterNew())
}

type CreateNewsletterFormPayload struct {
	Title           string `json:"title"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	IsPublished     bool   `json:"is_published"`
	ReleasedAt      string `json:"released_at"`
	Slug            string `json:"slug"`
	Content         string `json:"content"`
}

func (n Newsletters) Create(etx echo.Context) error {
	var payload CreateNewsletterFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateNewsletterFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.CreateNewsletterData{
		Title:           payload.Title,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		IsPublished:     payload.IsPublished,
		ReleasedAt: func() time.Time {
			if payload.ReleasedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.ReleasedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Slug:    payload.Slug,
		Content: payload.Content,
	}

	newsletter, err := models.CreateNewsletter(
		etx.Request().Context(),
		n.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to create newsletter: %v", err)); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.NewsletterNew.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Newsletter created successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.NewsletterShow.URL(newsletter.ID))
}

func (n Newsletters) Edit(etx echo.Context) error {
	newsletterID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	newsletter, err := models.FindNewsletter(etx.Request().Context(), n.db.Conn(), newsletterID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.NewsletterEdit(newsletter))
}

type UpdateNewsletterFormPayload struct {
	Title           string `json:"title"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	IsPublished     bool   `json:"is_published"`
	ReleasedAt      string `json:"released_at"`
	Slug            string `json:"slug"`
	Content         string `json:"content"`
}

func (n Newsletters) Update(etx echo.Context) error {
	newsletterID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	var payload UpdateNewsletterFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateNewsletterFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.UpdateNewsletterData{
		ID:              newsletterID,
		Title:           payload.Title,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		IsPublished:     payload.IsPublished,
		ReleasedAt: func() time.Time {
			if payload.ReleasedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.ReleasedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Slug:    payload.Slug,
		Content: payload.Content,
	}

	newsletter, err := models.UpdateNewsletter(
		etx.Request().Context(),
		n.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to update newsletter: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.NewsletterEdit.URL(newsletterID),
		)
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Newsletter updated successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.NewsletterShow.URL(newsletter.ID))
}

func (n Newsletters) Destroy(etx echo.Context) error {
	newsletterID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	err = models.DestroyNewsletter(etx.Request().Context(), n.db.Conn(), newsletterID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to delete newsletter: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.NewsletterIndex.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Newsletter destroyed successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.NewsletterIndex.URL())
}
