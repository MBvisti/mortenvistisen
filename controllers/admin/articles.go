package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
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

type Articles struct {
	db                     storage.Pool
	service                services.Article
	maxArticleRequestBytes int64
}

func NewArticles(db storage.Pool, service services.Article, cfg config.Config) Articles {
	return Articles{
		db:                     db,
		service:                service,
		maxArticleRequestBytes: 2*cfg.R2.MaxUploadBytes + 5<<20,
	}
}

func (a Articles) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminArticleIndex.Path(),
		Name:    routes.AdminArticleIndex.Name(),
		Handler: a.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminArticleShow.Path(),
		Name:    routes.AdminArticleShow.Name(),
		Handler: a.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminArticleNew.Path(),
		Name:    routes.AdminArticleNew.Name(),
		Handler: a.New,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.AdminArticleCreate.Path(),
		Name:    routes.AdminArticleCreate.Name(),
		Handler: a.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminArticleEdit.Path(),
		Name:    routes.AdminArticleEdit.Name(),
		Handler: a.Edit,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPut,
		Path:    routes.AdminArticleUpdate.Path(),
		Name:    routes.AdminArticleUpdate.Name(),
		Handler: a.Update,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodDelete,
		Path:    routes.AdminArticleDestroy.Path(),
		Name:    routes.AdminArticleDestroy.Name(),
		Handler: a.Destroy,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type ArticleData struct {
	ID               int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
	FirstPublishedAt time.Time
	Published        bool
	Title            string
	Excerpt          string
	MetaTitle        string
	MetaDescription  string
	Slug             string
	ImageLink        string
	MetaImageLink    string
	ReadTime         int32
	Content          string
	Tags             []TagData
}

func newArticleData(entity models.ArticleEntity) ArticleData {
	return ArticleData{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		FirstPublishedAt: entity.FirstPublishedAt.Time,
		Published:        entity.Published,
		Title:            entity.Title,
		Excerpt:          entity.Excerpt.String,
		MetaTitle:        entity.MetaTitle.String,
		MetaDescription:  entity.MetaDescription.String,
		Slug:             entity.Slug,
		ImageLink:        entity.ImageLink.String,
		MetaImageLink:    entity.MetaImageLink.String,
		ReadTime:         entity.ReadTime.Int32,
		Content:          entity.Content.String,
	}
}

func (a Articles) articleDataWithTags(
	ctx context.Context,
	entity models.ArticleEntity,
) (ArticleData, error) {
	tags, err := models.ArticleTagConnection.TagsForArticle(ctx, a.db.Executor(), entity.ID)
	if err != nil {
		return ArticleData{}, err
	}

	item := newArticleData(entity)
	item.Tags = newTagDataList(tags)

	return item, nil
}

func (a Articles) tagOptions(ctx context.Context) ([]TagData, error) {
	tags, err := models.Tag.All(ctx, a.db.Executor())
	if err != nil {
		return nil, err
	}

	return newTagDataList(tags), nil
}

func newArticleDataList(entities []models.ArticleEntity) []ArticleData {
	items := make([]ArticleData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newArticleData(entity))
	}

	return items
}

type ArticlePaginationData struct {
	Page       int64
	PageSize   int64
	TotalCount int64
	TotalPages int64
}

type ArticleFiltersData struct {
	Status string
}

const articlePageSize int64 = 10

func articleStatusFilter(status string) models.ArticleFilter {
	published := status == "published"
	if published || status == "draft" {
		return models.ArticleFilter{Published: &published}
	}
	return models.ArticleFilter{}
}

func (a Articles) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	status := etx.QueryParam("status")
	filter := articleStatusFilter(status)
	if filter.Published == nil {
		status = ""
	}

	articlesList, err := models.Article.Paginate(
		etx.Request().Context(),
		a.db.Executor(),
		filter,
		page,
		articlePageSize,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Article/Index", inertia.Props{
		"items": newArticleDataList(articlesList.Articles),
		"pagination": ArticlePaginationData{
			Page:       articlesList.Page,
			PageSize:   articlesList.PageSize,
			TotalCount: articlesList.TotalCount,
			TotalPages: articlesList.TotalPages,
		},
		"filters": ArticleFiltersData{Status: status},
	})
}

func (a Articles) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	articleID := int32(parsed)

	article, err := models.Article.Find(etx.Request().Context(), a.db.Executor(), articleID)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	item, err := a.articleDataWithTags(etx.Request().Context(), article)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Article/Show", inertia.Props{"item": item})
}

func (a Articles) New(etx *echo.Context) error {
	tagOptions, err := a.tagOptions(etx.Request().Context())
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Article/Create", inertia.Props{
		"validationRules": models.ArticleValidationRules(true),
		"tagOptions":      tagOptions,
	})
}

type coverMultipartForm struct {
	Payload string `form:"payload"`
}

type imageFileHeaders struct {
	cover     *multipart.FileHeader
	metaImage *multipart.FileHeader
}

func multipartImageHeader(
	form *multipart.Form,
	field, label string,
) (*multipart.FileHeader, error) {
	files := form.File[field]
	if len(files) > 1 {
		return nil, fmt.Errorf("only one %s may be uploaded", label)
	}
	if len(files) == 1 {
		return files[0], nil
	}
	return nil, nil
}

func (a Articles) bindArticleForm(etx *echo.Context, payload any) (imageFileHeaders, error) {
	request := etx.Request()
	request.Body = http.MaxBytesReader(
		etx.Response(),
		request.Body,
		a.maxArticleRequestBytes,
	)

	var form coverMultipartForm
	if err := etx.Bind(&form); err != nil {
		return imageFileHeaders{}, fmt.Errorf("bind article multipart form: %w", err)
	}
	if form.Payload == "" {
		return imageFileHeaders{}, errors.New("article form payload is required")
	}
	if err := json.Unmarshal([]byte(form.Payload), payload); err != nil {
		return imageFileHeaders{}, fmt.Errorf("decode article form payload: %w", err)
	}
	multipartForm, err := etx.MultipartForm()
	if err != nil {
		return imageFileHeaders{}, fmt.Errorf("read article multipart files: %w", err)
	}
	if multipartForm == nil {
		return imageFileHeaders{}, errors.New("article request must use multipart form data")
	}
	cover, err := multipartImageHeader(multipartForm, "cover", "article cover")
	if err != nil {
		return imageFileHeaders{}, err
	}
	metaImage, err := multipartImageHeader(multipartForm, "metaImage", "article social image")
	if err != nil {
		return imageFileHeaders{}, err
	}
	return imageFileHeaders{cover: cover, metaImage: metaImage}, nil
}

func openCover(
	header *multipart.FileHeader,
) (*services.Cover, func() error, error) {
	if header == nil {
		return nil, nil, nil
	}

	file, err := header.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("open uploaded article cover: %w", err)
	}

	return &services.Cover{
		Filename: header.Filename,
		Size:     header.Size,
		Body:     file,
	}, file.Close, nil
}

type CreateArticleFormPayload struct {
	Published       bool     `json:"published"`
	Title           string   `json:"title"`
	Excerpt         string   `json:"excerpt"`
	MetaTitle       string   `json:"metaTitle"`
	MetaDescription string   `json:"metaDescription"`
	Slug            string   `json:"slug"`
	ReadTime        int32    `json:"readTime"`
	Content         string   `json:"content"`
	TagIDs          []int32  `json:"tagIds"`
	NewTagTitles    []string `json:"newTags"`
}

func (a Articles) Create(etx *echo.Context) error {
	var payload CreateArticleFormPayload
	imageHeaders, err := a.bindArticleForm(etx, &payload)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateArticleFormPayload",
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

	data := models.CreateArticleData{
		Published:       payload.Published,
		Title:           payload.Title,
		Excerpt:         sql.NullString{String: payload.Excerpt, Valid: true},
		MetaTitle:       sql.NullString{String: payload.MetaTitle, Valid: true},
		MetaDescription: sql.NullString{String: payload.MetaDescription, Valid: true},
		Slug:            payload.Slug,
		ReadTime:        sql.NullInt32{Int32: payload.ReadTime, Valid: true},
		Content:         sql.NullString{String: payload.Content, Valid: true},
	}

	article, err := a.service.Create(
		etx.Request().Context(),
		data,
		services.ArticleTagSelection{
			TagIDs:       payload.TagIDs,
			NewTagTitles: payload.NewTagTitles,
		},
		cover,
		metaImage,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			tagOptions, optionsErr := a.tagOptions(etx.Request().Context())
			if optionsErr != nil {
				return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
			}
			return inertia.Page(
				etx,
				"Admin/Article/Create",
				inertia.Props{
					"validationRules": models.ArticleValidationRules(true),
					"tagOptions":      tagOptions,
				},
				inertia.WithValidationErrors(
					validationErrors.ToMap(),
				),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create article: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminArticleNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Article created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminArticleShow.URL(article.ID))
}

func (a Articles) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	articleID := int32(parsed)

	article, err := models.Article.Find(etx.Request().Context(), a.db.Executor(), articleID)
	if err != nil {
		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	item, err := a.articleDataWithTags(etx.Request().Context(), article)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	tagOptions, err := a.tagOptions(etx.Request().Context())
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Article/Edit", inertia.Props{
		"item":            item,
		"validationRules": models.ArticleValidationRules(true),
		"tagOptions":      tagOptions,
	})
}

type UpdateArticleFormPayload struct {
	Published       bool     `json:"published"`
	Title           string   `json:"title"`
	Excerpt         string   `json:"excerpt"`
	MetaTitle       string   `json:"metaTitle"`
	MetaDescription string   `json:"metaDescription"`
	Slug            string   `json:"slug"`
	ReadTime        int32    `json:"readTime"`
	Content         string   `json:"content"`
	TagIDs          []int32  `json:"tagIds"`
	NewTagTitles    []string `json:"newTags"`
	RemoveCover     bool     `json:"removeCover"`
	RemoveMetaImage bool     `json:"removeMetaImage"`
}

func (a Articles) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	articleID := int32(parsed)

	var payload UpdateArticleFormPayload
	imageHeaders, err := a.bindArticleForm(etx, &payload)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateArticleFormPayload",
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

	data := models.UpdateArticleData{
		ID:              articleID,
		Published:       payload.Published,
		Title:           payload.Title,
		Excerpt:         sql.NullString{String: payload.Excerpt, Valid: true},
		MetaTitle:       sql.NullString{String: payload.MetaTitle, Valid: true},
		MetaDescription: sql.NullString{String: payload.MetaDescription, Valid: true},
		Slug:            payload.Slug,
		ReadTime:        sql.NullInt32{Int32: payload.ReadTime, Valid: true},
		Content:         sql.NullString{String: payload.Content, Valid: true},
	}

	article, err := a.service.Update(
		etx.Request().Context(),
		data,
		services.ArticleTagSelection{
			TagIDs:       payload.TagIDs,
			NewTagTitles: payload.NewTagTitles,
		},
		cover,
		metaImage,
		payload.RemoveCover,
		payload.RemoveMetaImage,
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			currentArticle, findErr := models.Article.Find(
				etx.Request().Context(),
				a.db.Executor(),
				articleID,
			)
			if findErr != nil {
				return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
			}

			item, itemErr := a.articleDataWithTags(etx.Request().Context(), currentArticle)
			if itemErr != nil {
				return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
			}
			tagOptions, optionsErr := a.tagOptions(etx.Request().Context())
			if optionsErr != nil {
				return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
			}

			return inertia.Page(
				etx,
				"Admin/Article/Edit",
				inertia.Props{
					"item":            item,
					"validationRules": models.ArticleValidationRules(true),
					"tagOptions":      tagOptions,
				},
				inertia.WithValidationErrors(
					validationErrors.ToMap(),
				),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update article: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminArticleEdit.URL(articleID),
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Article updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminArticleShow.URL(article.ID))
}

func (a Articles) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	articleID := int32(parsed)

	err = models.Article.Destroy(etx.Request().Context(), a.db.Executor(), articleID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete article: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminArticleIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Article destroyed successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminArticleIndex.URL())
}
