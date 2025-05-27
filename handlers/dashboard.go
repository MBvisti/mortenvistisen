package handlers

import (
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
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

func (d Dashboard) Index(ctx echo.Context) error {
	// Get page parameter from query string, default to 1
	pageStr := ctx.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get articles with pagination
	pageSize := 10
	result, err := models.GetArticlesPaginated(
		setAppCtx(ctx),
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

	return dashboard.Home(result).Render(renderArgs(ctx))
}

func (d Dashboard) NewArticle(ctx echo.Context) error {
	formData := dashboard.NewArticleFormData{
		CsrfToken: csrf.Token(ctx.Request()),
		Errors:    make(map[string][]string),
	}

	return dashboard.NewArticle(formData).Render(renderArgs(ctx))
}

func (d Dashboard) StoreArticle(ctx echo.Context) error {
	// Parse form data
	title := ctx.FormValue("title")
	excerpt := ctx.FormValue("excerpt")
	metaTitle := ctx.FormValue("meta_title")
	metaDescription := ctx.FormValue("meta_description")
	slug := ctx.FormValue("slug")
	imageLink := ctx.FormValue("image_link")
	content := ctx.FormValue("content")
	action := ctx.FormValue("action") // "save_draft" or "publish"

	// Prepare form data for re-rendering on error
	formData := dashboard.NewArticleFormData{
		Title:           title,
		Excerpt:         excerpt,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		Slug:            slug,
		ImageLink:       imageLink,
		Content:         content,
		CsrfToken:       csrf.Token(ctx.Request()),
		Errors:          make(map[string][]string),
	}

	// Prepare payload
	payload := models.NewArticlePayload{
		Title:           title,
		Excerpt:         excerpt,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		Slug:            slug,
		Content:         &content,
	}

	// Handle optional image link
	if imageLink != "" {
		payload.ImageLink = &imageLink
	}

	// Create the article
	article, err := models.NewArticle(setAppCtx(ctx), d.db.Pool, payload)
	if err != nil {
		// Handle validation errors
		if validationErrors := extractValidationErrors(err); validationErrors != nil {
			formData.Errors = validationErrors
			return dashboard.NewArticle(formData).Render(renderArgs(ctx))
		}

		// Handle other errors (like duplicate slug)
		_ = addFlash(ctx, "error", "Failed to create article. Please try again.")
		return dashboard.NewArticle(formData).Render(renderArgs(ctx))
	}

	// If action is "publish", publish the article
	if action == "publish" {
		now := time.Now()
		_, err = models.PublishArticle(setAppCtx(ctx), d.db.Pool, models.PublishArticlePayload{
			ID:          article.ID,
			UpdatedAt:   now,
			PublishedAt: now,
		})
		if err != nil {
			// Article was created but publishing failed
			_ = addFlash(ctx, "warning", "Article created as draft. Publishing failed.")
		} else {
			_ = addFlash(ctx, "success", "Article published successfully!")
		}
	} else {
		_ = addFlash(ctx, "success", "Article saved as draft!")
	}

	// Redirect to dashboard
	return ctx.Redirect(302, "/dashboard")
}
