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
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	articles, err := models.GetArticlesPaginated(
		extractCtx(c),
		d.db.Pool,
		page,
		pageSize,
	)
	if err != nil {
		articles = models.PaginationResult{
			Articles:    []models.Article{},
			TotalCount:  0,
			Page:        1,
			PageSize:    pageSize,
			TotalPages:  0,
			HasNext:     false,
			HasPrevious: false,
		}
	}

	return dashboard.Home(articles).Render(renderArgs(c))
}

func (d Dashboard) NewArticle(c echo.Context) error {
	availableTags, err := models.GetArticleTags(extractCtx(c), d.db.Pool)
	if err != nil {
		return err
	}

	tagOptions := make([]dashboard.ArticleTagOption, len(availableTags))
	for i, tag := range availableTags {
		tagOptions[i] = dashboard.ArticleTagOption{
			ID:    tag.ID.String(),
			Title: tag.Title,
		}
	}

	formData := dashboard.NewArticleFormData{
		AvailableTags: tagOptions,
		CsrfToken:     csrf.Token(c.Request()),
		Errors:        make(map[string][]string),
	}

	return dashboard.NewArticle(formData).Render(renderArgs(c))
}

func (d Dashboard) StoreArticle(c echo.Context) error {
	type payload struct {
		Title           string   `form:"title"`
		Excerpt         string   `form:"excerpt"`
		MetaTitle       string   `form:"meta_title"`
		MetaDescription string   `form:"meta_description"`
		Publish         string   `form:"published"`
		Slug            string   `form:"slug"`
		ImageLink       string   `form:"image_link"`
		Content         string   `form:"content"`
		ReadTime        string   `form:"read_time"`
		TagIDs          []string `form:"tag_ids"`
		Action          string   `form:"action"`
	}

	var articlePayload payload
	if err := c.Bind(&articlePayload); err != nil {
		return err
	}

	var readTime int32
	if articlePayload.ReadTime != "" {
		if parsedTime, err := strconv.ParseInt(articlePayload.ReadTime, 10, 32); err == nil {
			readTime = int32(parsedTime)
		}
	}

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
			ReadTime:        readTime,
			TagIDs:          articlePayload.TagIDs,
		},
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
			"Failed to create article. Please try again.",
		); err != nil {
			return err
		}

		return dashboard.NewArticle(dashboard.NewArticleFormData{}).
			Render(renderArgs(c))
	}

	if articlePayload.Action == "publish" {
		_, err = models.PublishArticle(
			extractCtx(c),
			d.db.Pool,
			models.PublishArticlePayload{
				ID:  article.ID,
				Now: time.Now(),
			},
		)
		if err != nil {
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

	availableTags, err := models.GetArticleTags(extractCtx(c), d.db.Pool)
	if err != nil {
		return err
	}

	tagOptions := make([]dashboard.ArticleTagOption, len(availableTags))
	for i, tag := range availableTags {
		tagOptions[i] = dashboard.ArticleTagOption{
			ID:    tag.ID.String(),
			Title: tag.Title,
		}
	}

	selectedTagIDs := make([]string, len(article.Tags))
	for i, tag := range article.Tags {
		selectedTagIDs[i] = tag.ID.String()
	}

	formData := dashboard.EditArticleFormData{
		ID:              article.ID.String(),
		Title:           article.Title,
		Excerpt:         article.Excerpt,
		MetaTitle:       article.MetaTitle,
		MetaDescription: article.MetaDescription,
		IsPublished:     article.IsPublished,
		Slug:            article.Slug,
		ImageLink:       article.ImageLink,
		Content:         article.Content,
		ReadTime:        strconv.FormatInt(int64(article.ReadTime), 10),
		SelectedTagIDs:  selectedTagIDs,
		AvailableTags:   tagOptions,
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
		Title           string   `form:"title"`
		Excerpt         string   `form:"excerpt"`
		MetaTitle       string   `form:"meta_title"`
		MetaDescription string   `form:"meta_description"`
		Slug            string   `form:"slug"`
		ImageLink       string   `form:"image_link"`
		Content         string   `form:"content"`
		ReadTime        string   `form:"read_time"`
		TagIDs          []string `form:"tag_ids"`
		Action          string   `form:"action"`
		Published       string   `form:"published"`
	}

	var articlePayload payload
	if err := c.Bind(&articlePayload); err != nil {
		return err
	}

	slog.Info("IS IT PUBLISHED OR NOT", "p", articlePayload.Published)

	var readTime int32
	if articlePayload.ReadTime != "" {
		if parsedTime, err := strconv.ParseInt(articlePayload.ReadTime, 10, 32); err == nil {
			readTime = int32(parsedTime)
		}
	}

	published := articlePayload.Published == "on"

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
		ReadTime:        readTime,
		TagIDs:          articlePayload.TagIDs,
		IsPublished:     published,
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

	if published && article.FirstPublishedAt.IsZero() {
		_, err := models.PublishArticle(
			c.Request().Context(),
			d.db.Pool,
			models.PublishArticlePayload{ID: article.ID, Now: time.Now()},
		)
		if err != nil {
			if err := addFlash(
				c,
				contexts.FlashError,
				"Failed to update article. Please try again.",
			); err != nil {
				return err
			}
			return err
		}
		if err := addFlash(
			c,
			contexts.FlashSuccess,
			"Article published successfully!",
		); err != nil {
			return err
		}

		return c.Redirect(302, "/dashboard")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Article updated successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard")
}

func (d Dashboard) DeleteArticle(c echo.Context) error {
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

	// Check if article exists before deleting
	_, err = models.GetArticleByID(
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

	// Delete the article
	err = models.DeleteArticle(
		extractCtx(c),
		d.db.Pool,
		articleID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to delete article. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Article deleted successfully!",
	); err != nil {
		return err
	}

	// Redirect to dashboard
	return c.Redirect(302, "/dashboard")
}

// Tag management handlers
func (d Dashboard) Tags(c echo.Context) error {
	tags, err := models.GetArticleTags(extractCtx(c), d.db.Pool)
	if err != nil {
		return err
	}

	data := dashboard.TagsPageData{
		Tags:      tags,
		CsrfToken: csrf.Token(c.Request()),
		Errors:    make(map[string][]string),
	}

	return dashboard.Tags(data).Render(renderArgs(c))
}

func (d Dashboard) CreateTag(c echo.Context) error {
	type payload struct {
		Title string `form:"title"`
	}

	var tagPayload payload
	if err := c.Bind(&tagPayload); err != nil {
		return err
	}

	_, err := models.NewArticleTag(
		extractCtx(c),
		d.db.Pool,
		models.NewArticleTagPayload{
			Title: tagPayload.Title,
		},
	)
	if err != nil {
		if validationErrors := extractValidationErrors(err); validationErrors != nil {
			tags, _ := models.GetArticleTags(extractCtx(c), d.db.Pool)
			data := dashboard.TagsPageData{
				Tags:      tags,
				CsrfToken: csrf.Token(c.Request()),
				Errors:    validationErrors,
			}
			return dashboard.Tags(data).Render(renderArgs(c))
		}

		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to create tag. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/tags")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Tag created successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard/tags")
}

func (d Dashboard) UpdateTag(c echo.Context) error {
	tagID := c.Param("id")
	parsedID, err := uuid.Parse(tagID)
	if err != nil {
		return echo.NewHTTPError(404, "Tag not found")
	}

	type payload struct {
		Title string `form:"title"`
	}

	var tagPayload payload
	if err := c.Bind(&tagPayload); err != nil {
		return err
	}

	_, err = models.UpdateArticleTag(
		extractCtx(c),
		d.db.Pool,
		models.UpdateArticleTagPayload{
			ID:        parsedID,
			UpdatedAt: time.Now(),
			Title:     tagPayload.Title,
		},
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to update tag. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/tags")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Tag updated successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard/tags")
}

func (d Dashboard) DeleteTag(c echo.Context) error {
	tagID := c.Param("id")
	parsedID, err := uuid.Parse(tagID)
	if err != nil {
		return echo.NewHTTPError(404, "Tag not found")
	}

	// First delete all connections with this tag
	err = models.DeleteArticleTagConnectionsByTagID(
		extractCtx(c),
		d.db.Pool,
		parsedID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to delete tag. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/tags")
	}

	// Then delete the tag itself
	err = models.DeleteArticleTag(extractCtx(c), d.db.Pool, parsedID)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to delete tag. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/tags")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Tag deleted successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard/tags")
}
