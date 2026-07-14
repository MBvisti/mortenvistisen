package admin

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mortenvistisen/config"
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

type Projects struct {
	db                     storage.Pool
	service                services.Project
	maxProjectRequestBytes int64
}

func NewProjects(db storage.Pool, service services.Project, cfg config.Config) Projects {
	return Projects{
		db:                     db,
		service:                service,
		maxProjectRequestBytes: 2*cfg.R2.MaxUploadBytes + 5<<20,
	}
}

func (p Projects) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectIndex.Path(),
		Name:    routes.AdminProjectIndex.Name(),
		Handler: p.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectShow.Path(),
		Name:    routes.AdminProjectShow.Name(),
		Handler: p.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectNew.Path(),
		Name:    routes.AdminProjectNew.Name(),
		Handler: p.New,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.AdminProjectCreate.Path(),
		Name:    routes.AdminProjectCreate.Name(),
		Handler: p.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectEdit.Path(),
		Name:    routes.AdminProjectEdit.Name(),
		Handler: p.Edit,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPut,
		Path:    routes.AdminProjectUpdate.Path(),
		Name:    routes.AdminProjectUpdate.Name(),
		Handler: p.Update,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodDelete,
		Path:    routes.AdminProjectDestroy.Path(),
		Name:    routes.AdminProjectDestroy.Name(),
		Handler: p.Destroy,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type ProjectData struct {
	ID              int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Published       bool
	Title           string
	Slug            string
	StartedAt       time.Time
	Status          string
	Description     string
	MetaTitle       string
	MetaDescription string
	ImageLink       string
	MetaImageLink   string
	Content         string
	ProjectUrl      string
}

func newProjectData(entity models.ProjectEntity) ProjectData {
	return ProjectData{
		ID:              entity.ID,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
		Published:       entity.Published,
		Title:           entity.Title,
		Slug:            entity.Slug,
		StartedAt:       entity.StartedAt.Time,
		Status:          entity.Status,
		Description:     entity.Description,
		MetaTitle:       entity.MetaTitle,
		MetaDescription: entity.MetaDescription,
		ImageLink:       entity.ImageLink.String,
		MetaImageLink:   entity.MetaImageLink.String,
		Content:         entity.Content,
		ProjectUrl:      entity.ProjectUrl.String,
	}
}

func newProjectDataList(entities []models.ProjectEntity) []ProjectData {
	items := make([]ProjectData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newProjectData(entity))
	}

	return items
}

func (p Projects) Index(etx *echo.Context) error {
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

	projectsList, err := models.Project.Paginate(
		etx.Request().Context(),
		p.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Project/Index", inertia.Props{
		"items": newProjectDataList(projectsList.Projects),
	})
}

func (p Projects) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	project, err := models.Project.Find(etx.Request().Context(), p.db.Executor(), projectID)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Project/Show", inertia.Props{
		"item": newProjectData(project),
	})
}

func (p Projects) New(etx *echo.Context) error {
	return inertia.Page(etx, "Admin/Project/Create", inertia.Props{
		"validationRules": models.ProjectValidationRules(true),
	})
}

func (p Projects) bindProjectForm(etx *echo.Context, payload any) (imageFileHeaders, error) {
	request := etx.Request()
	request.Body = http.MaxBytesReader(
		etx.Response(),
		request.Body,
		p.maxProjectRequestBytes,
	)
	var form coverMultipartForm
	if err := etx.Bind(&form); err != nil {
		return imageFileHeaders{}, fmt.Errorf("bind project multipart form: %w", err)
	}
	if form.Payload == "" {
		return imageFileHeaders{}, errors.New("project form payload is required")
	}
	if err := json.Unmarshal([]byte(form.Payload), payload); err != nil {
		return imageFileHeaders{}, fmt.Errorf("decode project form payload: %w", err)
	}
	multipartForm, err := etx.MultipartForm()
	if err != nil {
		return imageFileHeaders{}, fmt.Errorf("read project multipart files: %w", err)
	}
	if multipartForm == nil {
		return imageFileHeaders{}, errors.New("project request must use multipart form data")
	}
	cover, err := multipartImageHeader(multipartForm, "cover", "project cover")
	if err != nil {
		return imageFileHeaders{}, err
	}
	metaImage, err := multipartImageHeader(multipartForm, "metaImage", "project social image")
	if err != nil {
		return imageFileHeaders{}, err
	}
	return imageFileHeaders{cover: cover, metaImage: metaImage}, nil
}

type CreateProjectFormPayload struct {
	Published       bool   `json:"published"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	StartedAt       string `json:"startedAt"`
	Status          string `json:"status"`
	Description     string `json:"description"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	Content         string `json:"content"`
	ProjectUrl      string `json:"projectUrl"`
}

func (p Projects) Create(etx *echo.Context) error {
	var payload CreateProjectFormPayload
	imageHeaders, err := p.bindProjectForm(etx, &payload)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateProjectFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	cover, closeCover, err := openCover(imageHeaders.cover)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if closeCover != nil {
		defer closeCover()
	}
	metaImage, closeMetaImage, err := openCover(imageHeaders.metaImage)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if closeMetaImage != nil {
		defer closeMetaImage()
	}

	data := models.CreateProjectData{
		Published: payload.Published,
		Title:     payload.Title,
		Slug:      payload.Slug,
		StartedAt: func() sql.NullTime {
			if payload.StartedAt == "" {
				return sql.NullTime{Valid: false}
			}
			if t, err := time.Parse("2006-01-02", payload.StartedAt); err == nil {
				return sql.NullTime{Time: t, Valid: true}
			}
			return sql.NullTime{Valid: false}
		}(),
		Status:          payload.Status,
		Description:     payload.Description,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		Content:         payload.Content,
		ProjectUrl:      sql.NullString{String: payload.ProjectUrl, Valid: true},
	}

	project, err := p.service.Create(
		etx.Request().Context(),
		data,
		cover,
		metaImage,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/Project/Create",
				inertia.Props{
					"validationRules": models.ProjectValidationRules(true),
				},
				inertia.WithValidationErrors(validationErrors.ToMap()),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create project: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminProjectNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Project created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectShow.URL(project.ID))
}

func (p Projects) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	project, err := models.Project.Find(etx.Request().Context(), p.db.Executor(), projectID)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Project/Edit", inertia.Props{
		"item":            newProjectData(project),
		"validationRules": models.ProjectValidationRules(true),
	})
}

type UpdateProjectFormPayload struct {
	Published       bool   `json:"published"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	StartedAt       string `json:"startedAt"`
	Status          string `json:"status"`
	Description     string `json:"description"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	Content         string `json:"content"`
	ProjectUrl      string `json:"projectUrl"`
	RemoveCover     bool   `json:"removeCover"`
	RemoveMetaImage bool   `json:"removeMetaImage"`
}

func (p Projects) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	var payload UpdateProjectFormPayload
	imageHeaders, err := p.bindProjectForm(etx, &payload)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateProjectFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	cover, closeCover, err := openCover(imageHeaders.cover)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if closeCover != nil {
		defer closeCover()
	}
	metaImage, closeMetaImage, err := openCover(imageHeaders.metaImage)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	if closeMetaImage != nil {
		defer closeMetaImage()
	}

	data := models.UpdateProjectData{
		ID:        projectID,
		Published: payload.Published,
		Title:     payload.Title,
		Slug:      payload.Slug,
		StartedAt: func() sql.NullTime {
			if payload.StartedAt == "" {
				return sql.NullTime{Valid: false}
			}
			if t, err := time.Parse("2006-01-02", payload.StartedAt); err == nil {
				return sql.NullTime{Time: t, Valid: true}
			}
			return sql.NullTime{Valid: false}
		}(),
		Status:          payload.Status,
		Description:     payload.Description,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		Content:         payload.Content,
		ProjectUrl:      sql.NullString{String: payload.ProjectUrl, Valid: true},
	}

	project, err := p.service.Update(
		etx.Request().Context(),
		data,
		cover,
		metaImage,
		payload.RemoveCover,
		payload.RemoveMetaImage,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			currentProject, findErr := models.Project.Find(
				etx.Request().Context(),
				p.db.Executor(),
				projectID,
			)
			if findErr != nil {
				return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
			}

			return inertia.Page(
				etx,
				"Admin/Project/Edit",
				inertia.Props{
					"item":            newProjectData(currentProject),
					"validationRules": models.ProjectValidationRules(true),
				},
				inertia.WithValidationErrors(validationErrors.ToMap()),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update project: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminProjectEdit.URL(projectID),
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Project updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectShow.URL(project.ID))
}

func (p Projects) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	err = models.Project.Destroy(etx.Request().Context(), p.db.Executor(), projectID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete project: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminProjectIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Project destroyed successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectIndex.URL())
}
