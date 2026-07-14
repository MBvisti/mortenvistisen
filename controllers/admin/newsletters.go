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

type Newsletters struct {
	db                        storage.Pool
	service                   services.Newsletter
	maxNewsletterRequestBytes int64
}

func NewNewsletters(db storage.Pool, service services.Newsletter, cfg config.Config) Newsletters {
	return Newsletters{
		db:                        db,
		service:                   service,
		maxNewsletterRequestBytes: 2*cfg.R2.MaxUploadBytes + 5<<20,
	}
}

func (n Newsletters) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminNewsletterIndex.Path(),
		Name:    routes.AdminNewsletterIndex.Name(),
		Handler: n.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminNewsletterShow.Path(),
		Name:    routes.AdminNewsletterShow.Name(),
		Handler: n.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminNewsletterNew.Path(),
		Name:    routes.AdminNewsletterNew.Name(),
		Handler: n.New,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.AdminNewsletterCreate.Path(),
		Name:    routes.AdminNewsletterCreate.Name(),
		Handler: n.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminNewsletterEdit.Path(),
		Name:    routes.AdminNewsletterEdit.Name(),
		Handler: n.Edit,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPut,
		Path:    routes.AdminNewsletterUpdate.Path(),
		Name:    routes.AdminNewsletterUpdate.Name(),
		Handler: n.Update,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodDelete,
		Path:    routes.AdminNewsletterDestroy.Path(),
		Name:    routes.AdminNewsletterDestroy.Name(),
		Handler: n.Destroy,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type NewsletterData struct {
	ID              int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Title           string
	Slug            string
	MetaTitle       string
	MetaDescription string
	IsPublished     bool
	ReleasedAt      time.Time
	ImageLink       string
	MetaImageLink   string
	Content         string
}

func newNewsletterData(entity models.NewsletterEntity) NewsletterData {
	return NewsletterData{
		ID:              entity.ID,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
		Title:           entity.Title,
		Slug:            entity.Slug.String,
		MetaTitle:       entity.MetaTitle,
		MetaDescription: entity.MetaDescription,
		IsPublished:     entity.IsPublished.Bool,
		ReleasedAt:      entity.ReleasedAt.Time,
		ImageLink:       entity.ImageLink.String,
		MetaImageLink:   entity.MetaImageLink.String,
		Content:         entity.Content.String,
	}
}

func newNewsletterDataList(entities []models.NewsletterEntity) []NewsletterData {
	items := make([]NewsletterData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newNewsletterData(entity))
	}

	return items
}

type NewsletterPaginationData struct {
	Page       int64
	PageSize   int64
	TotalCount int64
	TotalPages int64
}

const newsletterPageSize int64 = 10

func (n Newsletters) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	newslettersList, err := models.Newsletter.Paginate(
		etx.Request().Context(),
		n.db.Executor(),
		page,
		newsletterPageSize,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Newsletter/Index", inertia.Props{
		"items": newNewsletterDataList(newslettersList.Newsletters),
		"pagination": NewsletterPaginationData{
			Page:       newslettersList.Page,
			PageSize:   newslettersList.PageSize,
			TotalCount: newslettersList.TotalCount,
			TotalPages: newslettersList.TotalPages,
		},
	})
}

func (n Newsletters) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	newsletterID := int32(parsed)

	newsletter, err := models.Newsletter.Find(
		etx.Request().Context(),
		n.db.Executor(),
		newsletterID,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Newsletter/Show", inertia.Props{
		"item": newNewsletterData(newsletter),
	})
}

func (n Newsletters) New(etx *echo.Context) error {
	return inertia.Page(etx, "Admin/Newsletter/Create", inertia.Props{
		"validationRules": models.NewsletterValidationRules(true),
	})
}

func (n Newsletters) bindNewsletterForm(etx *echo.Context, payload any) (imageFileHeaders, error) {
	request := etx.Request()
	request.Body = http.MaxBytesReader(
		etx.Response(),
		request.Body,
		n.maxNewsletterRequestBytes,
	)
	var form coverMultipartForm
	if err := etx.Bind(&form); err != nil {
		return imageFileHeaders{}, fmt.Errorf("bind newsletter multipart form: %w", err)
	}
	if form.Payload == "" {
		return imageFileHeaders{}, errors.New("newsletter form payload is required")
	}
	if err := json.Unmarshal([]byte(form.Payload), payload); err != nil {
		return imageFileHeaders{}, fmt.Errorf("decode newsletter form payload: %w", err)
	}
	multipartForm, err := etx.MultipartForm()
	if err != nil {
		return imageFileHeaders{}, fmt.Errorf("read newsletter multipart files: %w", err)
	}
	if multipartForm == nil {
		return imageFileHeaders{}, errors.New("newsletter request must use multipart form data")
	}
	cover, err := multipartImageHeader(multipartForm, "cover", "newsletter cover")
	if err != nil {
		return imageFileHeaders{}, err
	}
	metaImage, err := multipartImageHeader(multipartForm, "metaImage", "newsletter social image")
	if err != nil {
		return imageFileHeaders{}, err
	}
	return imageFileHeaders{cover: cover, metaImage: metaImage}, nil
}

type CreateNewsletterFormPayload struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	IsPublished     bool   `json:"isPublished"`
	Content         string `json:"content"`
}

func (n Newsletters) Create(etx *echo.Context) error {
	var payload CreateNewsletterFormPayload
	imageHeaders, err := n.bindNewsletterForm(etx, &payload)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateNewsletterFormPayload",
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

	data := models.CreateNewsletterData{
		Title:           payload.Title,
		Slug:            sql.NullString{String: payload.Slug, Valid: true},
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		IsPublished:     sql.NullBool{Bool: payload.IsPublished, Valid: true},
		Content:         sql.NullString{String: payload.Content, Valid: true},
	}

	newsletter, err := n.service.Create(
		etx.Request().Context(),
		data,
		cover,
		metaImage,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/Newsletter/Create",
				inertia.Props{
					"validationRules": models.NewsletterValidationRules(true),
				},
				inertia.WithValidationErrors(validationErrors.ToMap()),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create newsletter: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminNewsletterNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Newsletter created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminNewsletterShow.URL(newsletter.ID))
}

func (n Newsletters) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	newsletterID := int32(parsed)

	newsletter, err := models.Newsletter.Find(
		etx.Request().Context(),
		n.db.Executor(),
		newsletterID,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Newsletter/Edit", inertia.Props{
		"item":            newNewsletterData(newsletter),
		"validationRules": models.NewsletterValidationRules(true),
	})
}

type UpdateNewsletterFormPayload struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	IsPublished     bool   `json:"isPublished"`
	Content         string `json:"content"`
	RemoveCover     bool   `json:"removeCover"`
	RemoveMetaImage bool   `json:"removeMetaImage"`
}

func (n Newsletters) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	newsletterID := int32(parsed)

	var payload UpdateNewsletterFormPayload
	imageHeaders, err := n.bindNewsletterForm(etx, &payload)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateNewsletterFormPayload",
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

	data := models.UpdateNewsletterData{
		ID:              newsletterID,
		Title:           payload.Title,
		Slug:            sql.NullString{String: payload.Slug, Valid: true},
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		IsPublished:     sql.NullBool{Bool: payload.IsPublished, Valid: true},
		Content:         sql.NullString{String: payload.Content, Valid: true},
	}

	newsletter, err := n.service.Update(
		etx.Request().Context(),
		data,
		cover,
		metaImage,
		payload.RemoveCover,
		payload.RemoveMetaImage,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			currentNewsletter, findErr := models.Newsletter.Find(
				etx.Request().Context(),
				n.db.Executor(),
				newsletterID,
			)
			if findErr != nil {
				return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
			}

			return inertia.Page(
				etx,
				"Admin/Newsletter/Edit",
				inertia.Props{
					"item":            newNewsletterData(currentNewsletter),
					"validationRules": models.NewsletterValidationRules(true),
				},
				inertia.WithValidationErrors(validationErrors.ToMap()),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update newsletter: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminNewsletterEdit.URL(newsletterID),
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Newsletter updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminNewsletterShow.URL(newsletter.ID))
}

func (n Newsletters) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	newsletterID := int32(parsed)

	err = models.Newsletter.Destroy(etx.Request().Context(), n.db.Executor(), newsletterID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete newsletter: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminNewsletterIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Newsletter destroyed successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminNewsletterIndex.URL())
}
