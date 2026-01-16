package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"mortenvistisen/internal/hypermedia"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Articles struct {
	db storage.Pool
}

func NewArticles(db storage.Pool) Articles {
	return Articles{db}
}

func (a Articles) Index(etx echo.Context) error {
	articles, err := models.AllArticles(
		etx.Request().Context(),
		a.db.Conn(),
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.ArticleIndex(articles))
}

func (a Articles) Show(etx echo.Context) error {
	article, err := models.FindArticleBySlug(
		etx.Request().Context(),
		a.db.Conn(),
		etx.Param("slug"),
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not find article",
			"error",
			err,
			"article_slug",
			etx.Param("slug"),
		)

		return render(etx, views.NotFound())
	}

	return render(etx, views.ArticleShow(article))
}

func (a Articles) New(etx echo.Context) error {
	allTags, err := models.AllTags(etx.Request().Context(), a.db.Conn())
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not fetch all tags",
			"error",
			err,
		)
		return render(etx, views.InternalError())
	}

	return render(etx, views.ArticleNew(allTags))
}

type CreateArticleFormPayload struct {
	FirstPublishedAt string   `json:"firstPublishedAt"`
	Excerpt          string   `json:"excerpt"          validate:"omitempty,min=10,max=255"`
	Title            string   `json:"title"            validate:"omitempty,min=3,max=100"`
	MetaTitle        string   `json:"metaTitle"        validate:"omitempty,min=3,max=100"`
	MetaDescription  string   `json:"metaDescription"  validate:"omitempty,min=10,max=255"`
	ImageLink        string   `json:"imageLink"        validate:"omitempty,url"`
	ReadTime         int32    `json:"readTime"         validate:"omitempty,min=1"`
	Content          string   `json:"content"          validate:"omitempty,min=20"`
	Published        bool     `json:"published"`
	TagIDs           []string `json:"tagIds"`
}

func (a Articles) ValidateArticlePayload(etx echo.Context) error {
	var payload CreateArticleFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateArticleFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	if err := models.Validate.Struct(payload); err != nil {

		slog.Info("article payload validation failed", "error", err)
		var validationErrors validator.ValidationErrors
		if !errors.As(err, &validationErrors) {
			slog.ErrorContext(
				etx.Request().Context(),
				"could not parse validation errors for article payload",
				"error",
				err,
			)
			return render(etx, views.NotFound())
		}

		for _, ve := range validationErrors {
			return hypermedia.PatchElementTempl(
				etx,
				views.InputField(
					"article"+ve.Field(),
					"text",
					strings.ToLower(ve.Field()),
					ve.Field(),
					"",
					ve.Value().(string),
					"ERR",
					true,
				))
		}
	}

	slog.Info("article payload validated successfully")

	return nil
}

func (a Articles) Create(etx echo.Context) error {
	var payload CreateArticleFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateArticleFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	firstPublishedAt, _ := time.Parse("2006-01-02", payload.FirstPublishedAt)

	data := models.CreateArticleData{
		FirstPublishedAt: firstPublishedAt,
		Title:            payload.Title,
		Excerpt:          payload.Excerpt,
		MetaTitle:        payload.MetaTitle,
		Published:        payload.Published,
		MetaDescription:  payload.MetaDescription,
		ImageLink:        payload.ImageLink,
		ReadTime:         payload.ReadTime,
		Content:          payload.Content,
	}

	article, err := models.CreateArticle(
		etx.Request().Context(),
		a.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to create article: %v", err)); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.ArticleNew.URL())
	}

	// Associate tags with the article
	if len(payload.TagIDs) > 0 {
		tagUUIDs := make([]uuid.UUID, 0, len(payload.TagIDs))
		for _, tagIDStr := range payload.TagIDs {
			tagID, err := uuid.Parse(tagIDStr)
			if err != nil {
				slog.ErrorContext(
					etx.Request().Context(),
					"could not parse tag UUID",
					"error",
					err,
					"tagID",
					tagIDStr,
				)
				continue
			}
			tagUUIDs = append(tagUUIDs, tagID)
		}

		if len(tagUUIDs) > 0 {
			if err := models.AssociateTagsWithArticle(etx.Request().Context(), a.db.Conn(), article.ID, tagUUIDs); err != nil {
				slog.ErrorContext(
					etx.Request().Context(),
					"could not associate tags with article",
					"error",
					err,
				)
				// Don't fail the whole request, just log the error
			}
		}
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Article created successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return hypermedia.Redirect(etx, routes.ArticleEdit.URL(article.ID))
}

func (a Articles) Edit(etx echo.Context) error {
	articleID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	article, err := models.FindArticle(etx.Request().Context(), a.db.Conn(), articleID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	allTags, err := models.AllTags(etx.Request().Context(), a.db.Conn())
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not fetch all tags",
			"error",
			err,
		)
		return render(etx, views.InternalError())
	}

	articleTags, err := models.FindTagsForArticle(etx.Request().Context(), a.db.Conn(), articleID)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not fetch article tags",
			"error",
			err,
		)
		return render(etx, views.InternalError())
	}

	articleTagIDs := make([]uuid.UUID, len(articleTags))
	for i, tag := range articleTags {
		articleTagIDs[i] = tag.ID
	}

	return render(etx, views.ArticleEdit(article, allTags, articleTagIDs))
}

type UpdateArticleFormPayload struct {
	FirstPublishedAt string   `json:"firstPublishedAt"`
	Excerpt          string   `json:"excerpt"          validate:"omitempty,min=10,max=255"`
	Title            string   `json:"title"            validate:"omitempty,min=3,max=100"`
	MetaTitle        string   `json:"metaTitle"        validate:"omitempty,min=3,max=100"`
	MetaDescription  string   `json:"metaDescription"  validate:"omitempty,min=10,max=255"`
	ImageLink        string   `json:"imageLink"        validate:"omitempty,url"`
	ReadTime         int32    `json:"readTime"         validate:"omitempty,min=1"`
	Content          string   `json:"content"          validate:"omitempty,min=20"`
	Published        bool     `json:"published"`
	TagIDs           []string `json:"tagIds"`
}

func (a Articles) Update(etx echo.Context) error {
	articleID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	var payload UpdateArticleFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateArticleFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	firstPublishedAt, _ := time.Parse("2006-01-02", payload.FirstPublishedAt)

	data := models.UpdateArticleData{
		ID:               articleID,
		FirstPublishedAt: firstPublishedAt,
		Title:            payload.Title,
		Excerpt:          payload.Excerpt,
		MetaTitle:        payload.MetaTitle,
		Published:        payload.Published,
		MetaDescription:  payload.MetaDescription,
		ImageLink:        payload.ImageLink,
		ReadTime:         payload.ReadTime,
		Content:          payload.Content,
	}

	article, err := models.UpdateArticle(
		etx.Request().Context(),
		a.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to update article: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.ArticleEdit.URL(articleID),
		)
	}

	// Clear existing tag associations and re-associate
	if err := models.ClearArticleTagAssociations(etx.Request().Context(), a.db.Conn(), articleID); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not clear article tag associations",
			"error",
			err,
		)
		// Don't fail the whole request, just log the error
	}

	if len(payload.TagIDs) > 0 {
		tagUUIDs := make([]uuid.UUID, 0, len(payload.TagIDs))
		for _, tagIDStr := range payload.TagIDs {
			tagID, err := uuid.Parse(tagIDStr)
			if err != nil {
				slog.ErrorContext(
					etx.Request().Context(),
					"could not parse tag UUID",
					"error",
					err,
					"tagID",
					tagIDStr,
				)
				continue
			}
			tagUUIDs = append(tagUUIDs, tagID)
		}

		if len(tagUUIDs) > 0 {
			if err := models.AssociateTagsWithArticle(etx.Request().Context(), a.db.Conn(), article.ID, tagUUIDs); err != nil {
				slog.ErrorContext(
					etx.Request().Context(),
					"could not associate tags with article",
					"error",
					err,
				)
				// Don't fail the whole request, just log the error
			}
		}
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Article updated successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ArticleShow.URL(article.ID))
}

func (a Articles) Destroy(etx echo.Context) error {
	articleID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	err = models.DestroyArticle(etx.Request().Context(), a.db.Conn(), articleID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to delete article: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.ArticleIndex.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Article destroyed successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ArticleIndex.URL())
}
