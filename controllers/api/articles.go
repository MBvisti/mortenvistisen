package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"mortenvistisen/config"
	"mortenvistisen/internal/validation"
	"mortenvistisen/models"
	"mortenvistisen/router"
	"mortenvistisen/router/middleware"
	"mortenvistisen/router/routes"
	"mortenvistisen/services"

	"github.com/labstack/echo/v5"
)

type Articles struct {
	article services.Article
	cfg     config.Config
}

func NewArticles(article services.Article, cfg config.Config) Articles {
	return Articles{article: article, cfg: cfg}
}

func (a Articles) RegisterRoutes(r *router.Router) error {
	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodPatch,
		Path:    routes.ApiArticleUpdate.Path(),
		Name:    routes.ApiArticleUpdate.Name(),
		Handler: a.Update,
		Middlewares: []echo.MiddlewareFunc{
			middleware.APIBasicAuth(a.cfg),
		},
	})
	return err
}

type updateArticlePayload struct {
	Title           *string `json:"title"`
	Excerpt         *string `json:"excerpt"`
	MetaTitle       *string `json:"metaTitle"`
	MetaDescription *string `json:"metaDescription"`
}

type articleResponse struct {
	Slug            string `json:"slug"`
	Title           string `json:"title"`
	Excerpt         string `json:"excerpt"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
}

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func (a Articles) Update(etx *echo.Context) error {
	slug := strings.TrimSpace(etx.Param("slug"))
	if slug == "" {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "article slug is required"})
	}

	var payload updateArticlePayload
	if err := etx.Bind(&payload); err != nil {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "invalid JSON request body"})
	}

	patch := services.ArticlePatch{
		Title:           payload.Title,
		Excerpt:         payload.Excerpt,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
	}
	if patch.Empty() {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "at least one field is required"})
	}

	article, err := a.article.Patch(etx.Request().Context(), slug, patch)
	if err != nil {
		if fields, ok := validation.As(err); ok {
			return etx.JSON(http.StatusBadRequest, errorResponse{
				Error:  "validation failed",
				Fields: fields.ToMap(),
			})
		}
		if errors.Is(err, models.ErrNotFound) {
			return etx.JSON(http.StatusNotFound, errorResponse{Error: "article not found"})
		}

		slog.ErrorContext(etx.Request().Context(), "failed to patch article", "error", err)
		return etx.JSON(http.StatusInternalServerError, errorResponse{Error: "failed to update article"})
	}

	return etx.JSON(http.StatusOK, articleResponse{
		Slug:            article.Slug,
		Title:           article.Title,
		Excerpt:         article.Excerpt.String,
		MetaTitle:       article.MetaTitle.String,
		MetaDescription: article.MetaDescription.String,
	})
}
