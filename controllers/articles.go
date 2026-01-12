package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

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
	articleID, err := uuid.Parse(etx.Param("id"))
	if err != nil {
		return render(etx, views.BadRequest())
	}

	article, err := models.FindArticle(etx.Request().Context(), a.db.Conn(), articleID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.ArticleShow(article))
}

func (a Articles) New(etx echo.Context) error {
	return render(etx, views.ArticleNew())
}

type CreateArticleFormPayload struct {
	FirstPublishedAt string `json:"first_published_at"`
	Title            string `json:"title"`
	Excerpt          string `json:"excerpt"`
	MetaTitle        string `json:"meta_title"`
	MetaDescription  string `json:"meta_description"`
	ImageLink        string `json:"image_link"`
	ReadTime         int32  `json:"read_time"`
	Content          string `json:"content"`
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

	data := models.CreateArticleData{
		FirstPublishedAt: func() time.Time {
			if payload.FirstPublishedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.FirstPublishedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Title:           payload.Title,
		Excerpt:         payload.Excerpt,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		ImageLink:       payload.ImageLink,
		ReadTime:        payload.ReadTime,
		Content:         payload.Content,
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

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Article created successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ArticleShow.URL(article.ID))
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

	return render(etx, views.ArticleEdit(article))
}

type UpdateArticleFormPayload struct {
	FirstPublishedAt string `json:"first_published_at"`
	Title            string `json:"title"`
	Excerpt          string `json:"excerpt"`
	MetaTitle        string `json:"meta_title"`
	MetaDescription  string `json:"meta_description"`
	Slug             string `json:"slug"`
	ImageLink        string `json:"image_link"`
	ReadTime         int32  `json:"read_time"`
	Content          string `json:"content"`
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

	data := models.UpdateArticleData{
		ID: articleID,
		FirstPublishedAt: func() time.Time {
			if payload.FirstPublishedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.FirstPublishedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Title:           payload.Title,
		Excerpt:         payload.Excerpt,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		Slug:            payload.Slug,
		ImageLink:       payload.ImageLink,
		ReadTime:        payload.ReadTime,
		Content:         payload.Content,
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
