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

	"github.com/labstack/echo/v5"
)

type Projects struct {
	db storage.Pool
}

func NewProjects(db storage.Pool) Projects {
	return Projects{db}
}

func (p Projects) Index(etx *echo.Context) error {
	page := int64(1)
	if rawPage := etx.QueryParam("page"); rawPage != "" {
		if parsed, err := strconv.Atoi(rawPage); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(10)
	if rawPerPage := etx.QueryParam("per_page"); rawPerPage != "" {
		if parsed, err := strconv.Atoi(rawPerPage); err == nil && parsed > 0 &&
			parsed <= 100 {
			perPage = int64(parsed)
		}
	}

	projectsList, err := models.PaginateProjects(
		etx.Request().Context(),
		p.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.ProjectIndex(projectsList))
}

func (p Projects) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	projectID := int32(parsed)

	project, err := models.FindProject(etx.Request().Context(), p.db.Conn(), projectID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.ProjectShow(project))
}

func (p Projects) New(etx *echo.Context) error {
	return render(etx, views.ProjectNew())
}

type CreateProjectFormPayload struct {
	Published   bool   `json:"published"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	StartedAt   string `json:"startedAt"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Content     string `json:"content"`
	ProjectURL  string `json:"projectUrl"`
}

func (p Projects) Create(etx *echo.Context) error {
	var payload CreateProjectFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateProjectFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.CreateProjectData{
		Published: payload.Published,
		Title:     payload.Title,
		Slug:      payload.Slug,
		StartedAt: func() time.Time {
			if payload.StartedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.StartedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Status:      payload.Status,
		Description: payload.Description,
		Content:     payload.Content,
		ProjectURL:  payload.ProjectURL,
	}

	project, err := models.CreateProject(
		etx.Request().Context(),
		p.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to create project: %v", err)); flashErr != nil {
			return flashErr
		}

		return etx.Redirect(http.StatusSeeOther, routes.ProjectNew.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Project created successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ProjectShow.URL(project.ID))
}

func (p Projects) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	projectID := int32(parsed)

	project, err := models.FindProject(etx.Request().Context(), p.db.Conn(), projectID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.ProjectUpdate(project))
}

type UpdateProjectFormPayload struct {
	Published   bool   `json:"published"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	StartedAt   string `json:"startedAt"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Content     string `json:"content"`
	ProjectURL  string `json:"projectUrl"`
}

func (p Projects) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	projectID := int32(parsed)

	var payload UpdateProjectFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateProjectFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.UpdateProjectData{
		ID:        projectID,
		Published: payload.Published,
		Title:     payload.Title,
		Slug:      payload.Slug,
		StartedAt: func() time.Time {
			if payload.StartedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.StartedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Status:      payload.Status,
		Description: payload.Description,
		Content:     payload.Content,
		ProjectURL:  payload.ProjectURL,
	}

	project, err := models.UpdateProject(
		etx.Request().Context(),
		p.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to update project: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}

		return etx.Redirect(
			http.StatusSeeOther,
			routes.ProjectEdit.URL(projectID),
		)
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Project updated successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ProjectShow.URL(project.ID))
}

func (p Projects) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	projectID := int32(parsed)

	err = models.DestroyProject(etx.Request().Context(), p.db.Conn(), projectID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to delete project: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}

		return etx.Redirect(http.StatusSeeOther, routes.ProjectIndex.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Project destroyed successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ProjectIndex.URL())
}
