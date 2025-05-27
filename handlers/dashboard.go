package handlers

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/router/contexts"
	"github.com/mbvisti/mortenvistisen/views/dashboard"
)

type Dashboard struct {
	db psql.Postgres
}

func newDashboard(db psql.Postgres) Dashboard {
	return Dashboard{
		db: db,
	}
}

func extractCtx(c echo.Context) context.Context {
	return c.Request().Context()
}

func (d Dashboard) Index(c echo.Context) error {
	// Get page parameter from query string, default to 1
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get articles with pagination
	pageSize := 10
	result, err := models.GetArticlesPaginated(
		extractCtx(c),
		d.db.Pool,
		page,
		pageSize,
	)
	if err != nil {
		// Log error and show empty state
		result = models.PaginationResult{
			Articles:    []models.Article{},
			TotalCount:  0,
			Page:        1,
			PageSize:    pageSize,
			TotalPages:  0,
			HasNext:     false,
			HasPrevious: false,
		}
	}

	return dashboard.Home(result).Render(renderArgs(c))
}

func (d Dashboard) NewArticle(c echo.Context) error {
	formData := dashboard.NewArticleFormData{
		CsrfToken: csrf.Token(c.Request()),
		Errors:    make(map[string][]string),
	}

	return dashboard.NewArticle(formData).Render(renderArgs(c))
}

func (d Dashboard) StoreArticle(c echo.Context) error {
	type payload struct {
		Title           string `form:"title"`
		Excerpt         string `form:"excerpt"`
		MetaTitle       string `form:"meta_title"`
		MetaDescription string `form:"meta_description"`
		Slug            string `form:"slug"`
		ImageLink       string `form:"image_link"`
		Content         string `form:"content"`
		Action          string `form:"action"`
	}

	var articlePayload payload
	if err := c.Bind(&articlePayload); err != nil {
		return err
	}

	slog.Info(
		"################### CONTENT ###################",
		"content",
		articlePayload.Content,
	)

	article, err := models.NewArticle(
		extractCtx(c),
		d.db.Pool,
		models.NewArticlePayload{
			Title:           articlePayload.Title,
			Excerpt:         articlePayload.Excerpt,
			MetaTitle:       articlePayload.MetaTitle,
			MetaDescription: articlePayload.MetaDescription,
			Slug:            articlePayload.Slug,
			ImageLink:       articlePayload.ImageLink,
			Content:         articlePayload.Content,
		},
	)
	if err != nil {
		if validationErrors := extractValidationErrors(err); validationErrors != nil {
			// formData.Errors = validationErrors
			// return dashboard.NewArticle(formData).Render(renderArgs(c))
			slog.ErrorContext(
				extractCtx(c),
				"could not validate article payload",
				"error",
				err,
			)
		}

		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to create article. Please try again.",
		); err != nil {
			return err
		}

		return dashboard.NewArticle(dashboard.NewArticleFormData{}).
			Render(renderArgs(c))
	}

	if articlePayload.Action == "on" {
		now := time.Now()
		_, err = models.PublishArticle(
			extractCtx(c),
			d.db.Pool,
			models.PublishArticlePayload{
				ID:          article.ID,
				UpdatedAt:   now,
				PublishedAt: now,
			},
		)
		if err != nil {
			// Article was created but publishing failed
			if err := addFlash(
				c,
				contexts.FlashError,
				"Article created as draft. Publishing failed.",
			); err != nil {
				return err
			}
		}

		if err := addFlash(
			c,
			contexts.FlashSuccess,
			"Article published successfully!",
		); err != nil {
			return err
		}
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Article saved as draft successfully!",
	); err != nil {
		return err
	}

	// Redirect to dashboard
	return c.Redirect(302, "/dashboard")
}

func (d Dashboard) EditArticle(c echo.Context) error {
	idParam := c.Param("id")
	articleID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid article ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard")
	}

	article, err := models.GetArticleByID(
		extractCtx(c),
		d.db.Pool,
		articleID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Article not found.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard")
	}

	formData := dashboard.EditArticleFormData{
		ID:              article.ID.String(),
		Title:           article.Title,
		Excerpt:         article.Excerpt,
		MetaTitle:       article.MetaTitle,
		MetaDescription: article.MetaDescription,
		Slug:            article.Slug,
		ImageLink:       article.ImageLink,
		Content:         article.Content,
		CsrfToken:       csrf.Token(c.Request()),
		Errors:          make(map[string][]string),
	}

	return dashboard.EditArticle(formData).Render(renderArgs(c))
}

func (d Dashboard) UpdateArticle(c echo.Context) error {
	idParam := c.Param("id")
	articleID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid article ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard")
	}

	type payload struct {
		Title           string `form:"title"`
		Excerpt         string `form:"excerpt"`
		MetaTitle       string `form:"meta_title"`
		MetaDescription string `form:"meta_description"`
		Slug            string `form:"slug"`
		ImageLink       string `form:"image_link"`
		Content         string `form:"content"`
		Action          string `form:"action"`
	}

	var articlePayload payload
	if err := c.Bind(&articlePayload); err != nil {
		return err
	}

	updatePayload := models.UpdateArticlePayload{
		ID:              articleID,
		UpdatedAt:       time.Now(),
		Title:           articlePayload.Title,
		Excerpt:         articlePayload.Excerpt,
		MetaTitle:       articlePayload.MetaTitle,
		MetaDescription: articlePayload.MetaDescription,
		Slug:            articlePayload.Slug,
		ImageLink:       articlePayload.ImageLink,
		Content:         articlePayload.Content,
	}

	if articlePayload.Action == "publish" {
		updatePayload.PublishedAt = time.Now()
	}

	article, err := models.UpdateArticle(
		extractCtx(c),
		d.db.Pool,
		updatePayload,
	)
	if err != nil {
		if validationErrors := extractValidationErrors(err); validationErrors != nil {
			slog.ErrorContext(
				extractCtx(c),
				"could not validate article payload",
				"error",
				err,
			)
		}

		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to update article. Please try again.",
		); err != nil {
			return err
		}

		formData := dashboard.EditArticleFormData{
			ID:              articleID.String(),
			Title:           articlePayload.Title,
			Excerpt:         articlePayload.Excerpt,
			MetaTitle:       articlePayload.MetaTitle,
			MetaDescription: articlePayload.MetaDescription,
			Slug:            articlePayload.Slug,
			ImageLink:       articlePayload.ImageLink,
			Content:         articlePayload.Content,
			CsrfToken:       csrf.Token(c.Request()),
			Errors:          make(map[string][]string),
		}

		return dashboard.EditArticle(formData).Render(renderArgs(c))
	}

	var successMsg string
	if article.IsPublished() {
		successMsg = "Article updated and published successfully!"
	} else {
		successMsg = "Article updated as draft successfully!"
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		successMsg,
	); err != nil {
		return err
	}

	// Redirect to dashboard
	return c.Redirect(302, "/dashboard")
}
