package handlers

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/router/contexts"
	"github.com/mbvisti/mortenvistisen/views/dashboard"
)

type Dashboard struct {
	db           psql.Postgres
	cacheManager *CacheManager
}

func newDashboard(db psql.Postgres, cacheManager *CacheManager) Dashboard {
	return Dashboard{
		db:           db,
		cacheManager: cacheManager,
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

	sortField := c.QueryParam("sort")
	sortOrder := c.QueryParam("order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc" // default
	}

	pageSize := 10
	articles, err := models.GetArticlesSorted(
		extractCtx(c),
		d.db.Pool,
		page,
		pageSize,
		models.SortConfig{
			Field: sortField,
			Order: sortOrder,
		},
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

	// Get complete statistics (not affected by pagination)
	publishedCount, err := models.CountPublishedArticles(extractCtx(c), d.db.Pool)
	if err != nil {
		publishedCount = 0
	}

	draftCount, err := models.CountDraftArticles(extractCtx(c), d.db.Pool)
	if err != nil {
		draftCount = 0
	}

	result := dashboard.DashboardSortableResult{
		Articles:         articles.Articles,
		TotalCount:       articles.TotalCount,
		Page:             articles.Page,
		PageSize:         articles.PageSize,
		TotalPages:       articles.TotalPages,
		HasNext:          articles.HasNext,
		HasPrevious:      articles.HasPrevious,
		CurrentSortField: sortField,
		CurrentSortOrder: sortOrder,
		PublishedCount:   publishedCount,
		DraftCount:       draftCount,
	}

	return dashboard.Home(result).Render(renderArgs(c))
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

	// Invalidate caches when article is created
	d.cacheManager.InvalidateLandingPage()

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

		// Invalidate caches when article is updated/published
		d.cacheManager.InvalidateLandingPage()
		d.cacheManager.InvalidateArticle(updatePayload.Slug)

		return c.Redirect(302, "/dashboard")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Article updated successfully!",
	); err != nil {
		return err
	}

	// Invalidate caches when article is updated
	d.cacheManager.InvalidateLandingPage()
	d.cacheManager.InvalidateArticle(updatePayload.Slug)

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

	// Check if article exists before deleting and get slug for cache invalidation
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

	// Invalidate caches when article is deleted
	d.cacheManager.InvalidateLandingPage()
	d.cacheManager.InvalidateArticle(article.Slug)

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
		Tags:   tags,
		Errors: make(map[string][]string),
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
				Tags:   tags,
				Errors: validationErrors,
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

func (d Dashboard) Subscribers(c echo.Context) error {
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	sortField := c.QueryParam("sort")
	sortOrder := c.QueryParam("order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc" // default
	}

	pageSize := 10
	subscribers, err := models.GetSubscribersSorted(
		extractCtx(c),
		d.db.Pool,
		page,
		pageSize,
		models.SortConfig{
			Field: sortField,
			Order: sortOrder,
		},
	)
	if err != nil {
		subscribers = models.SubscriberPaginationResult{
			Subscribers: []models.Subscriber{},
			TotalCount:  0,
			Page:        1,
			PageSize:    pageSize,
			TotalPages:  0,
			HasNext:     false,
			HasPrevious: false,
		}
	}

	monthlyCount, _ := models.CountMonthlySubscribers(extractCtx(c), d.db.Pool)
	verifiedCount, _ := models.CountVerifiedSubscribers(extractCtx(c), d.db.Pool)
	unverifiedCount := subscribers.TotalCount - verifiedCount

	result := dashboard.SubscribersSortableResult{
		Subscribers:      subscribers.Subscribers,
		TotalCount:       subscribers.TotalCount,
		Page:             subscribers.Page,
		PageSize:         subscribers.PageSize,
		TotalPages:       subscribers.TotalPages,
		HasNext:          subscribers.HasNext,
		HasPrevious:      subscribers.HasPrevious,
		CurrentSortField: sortField,
		CurrentSortOrder: sortOrder,
		MonthlyCount:     monthlyCount,
		VerifiedCount:    verifiedCount,
		UnverifiedCount:  unverifiedCount,
	}

	return dashboard.SubscribersSortable(result).Render(renderArgs(c))
}

func (d Dashboard) EditSubscriber(c echo.Context) error {
	idParam := c.Param("id")
	subscriberID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid subscriber ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/subscribers")
	}

	subscriber, err := models.GetSubscriber(
		extractCtx(c),
		d.db.Pool,
		subscriberID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Subscriber not found.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/subscribers")
	}

	formData := dashboard.EditSubscriberFormData{
		ID:         subscriber.ID.String(),
		Email:      subscriber.Email,
		Referer:    subscriber.Referer,
		IsVerified: subscriber.IsVerified,
		Errors:     make(map[string][]string),
	}

	return dashboard.EditSubscriber(formData).Render(renderArgs(c))
}

func (d Dashboard) UpdateSubscriber(c echo.Context) error {
	idParam := c.Param("id")
	subscriberID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid subscriber ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/subscribers")
	}

	type payload struct {
		Email      string `form:"email"`
		Referer    string `form:"referer"`
		IsVerified string `form:"is_verified"`
	}

	var subscriberPayload payload
	if err := c.Bind(&subscriberPayload); err != nil {
		return err
	}

	isVerified := subscriberPayload.IsVerified == "on"

	_, err = models.UpdateSubscriber(
		extractCtx(c),
		d.db.Pool,
		models.UpdateSubscriberPayload{
			ID:        subscriberID,
			UpdatedAt: time.Now(),
			Email:     subscriberPayload.Email,
			Referer:   subscriberPayload.Referer,
		},
	)
	if err != nil {
		if validationErrors := extractValidationErrors(err); validationErrors != nil {
			formData := dashboard.EditSubscriberFormData{
				ID:      subscriberID.String(),
				Email:   subscriberPayload.Email,
				Referer: subscriberPayload.Referer,
				Errors:  validationErrors,
			}
			return dashboard.EditSubscriber(formData).Render(renderArgs(c))
		}

		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to update subscriber. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/subscribers")
	}

	err = models.VerifySubscriber(
		extractCtx(c),
		d.db.Pool,
		models.VerifySubscriberPayload{
			ID:         subscriberID,
			UpdatedAt:  time.Now(),
			IsVerified: isVerified,
		},
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to update verification status. Please try again.",
		); err != nil {
			return err
		}
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Subscriber updated successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard/subscribers")
}

func (d Dashboard) DeleteSubscriber(c echo.Context) error {
	idParam := c.Param("id")
	subscriberID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid subscriber ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/subscribers")
	}

	_, err = models.GetSubscriber(
		extractCtx(c),
		d.db.Pool,
		subscriberID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Subscriber not found.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/subscribers")
	}

	err = models.DeleteSubscriber(
		extractCtx(c),
		d.db.Pool,
		subscriberID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to delete subscriber. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/subscribers")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Subscriber deleted successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard/subscribers")
}

// Newsletter management handlers
func (d Dashboard) Newsletters(c echo.Context) error {
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	newsletters, err := models.GetNewslettersPaginated(
		extractCtx(c),
		d.db.Pool,
		page,
		pageSize,
	)
	if err != nil {
		newsletters = models.NewsletterPaginationResult{
			Newsletters: []models.Newsletter{},
			TotalCount:  0,
			Page:        1,
			PageSize:    pageSize,
			TotalPages:  0,
			HasNext:     false,
			HasPrevious: false,
		}
	}

	return dashboard.Newsletters(newsletters).Render(renderArgs(c))
}

func (d Dashboard) NewNewsletter(c echo.Context) error {
	formData := dashboard.NewNewsletterFormData{
		Errors: make(map[string][]string),
	}

	return dashboard.NewNewsletter(formData).Render(renderArgs(c))
}

func (d Dashboard) StoreNewsletter(c echo.Context) error {
	type payload struct {
		Title   string `form:"title"`
		Slug    string `form:"slug"`
		Content string `form:"content"`
		Action  string `form:"action"`
	}

	var newsletterPayload payload
	if err := c.Bind(&newsletterPayload); err != nil {
		return err
	}

	newsletter, err := models.NewNewsletter(
		extractCtx(c),
		d.db.Pool,
		models.NewNewsletterPayload{
			Title:   newsletterPayload.Title,
			Slug:    newsletterPayload.Slug,
			Content: newsletterPayload.Content,
		},
	)
	if err != nil {
		if validationErrors := extractValidationErrors(err); validationErrors != nil {
			slog.ErrorContext(
				extractCtx(c),
				"could not validate newsletter payload",
				"error",
				err,
			)
		}

		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to create newsletter. Please try again.",
		); err != nil {
			return err
		}

		return dashboard.NewNewsletter(dashboard.NewNewsletterFormData{}).
			Render(renderArgs(c))
	}

	if newsletterPayload.Action == "publish" {
		_, err = models.PublishNewsletter(
			extractCtx(c),
			d.db.Pool,
			models.PublishNewsletterPayload{
				ID:  newsletter.ID,
				Now: time.Now(),
			},
		)
		if err != nil {
			if err := addFlash(
				c,
				contexts.FlashError,
				"Newsletter created as draft. Publishing failed.",
			); err != nil {
				return err
			}
		}

		if err := addFlash(
			c,
			contexts.FlashSuccess,
			"Newsletter published successfully!",
		); err != nil {
			return err
		}
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Newsletter saved as draft successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard/newsletters")
}

func (d Dashboard) EditNewsletter(c echo.Context) error {
	idParam := c.Param("id")
	newsletterID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid newsletter ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/newsletters")
	}

	newsletter, err := models.GetNewsletterByID(
		extractCtx(c),
		d.db.Pool,
		newsletterID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Newsletter not found.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/newsletters")
	}

	formData := dashboard.EditNewsletterFormData{
		ID:          newsletter.ID.String(),
		Title:       newsletter.Title,
		IsPublished: newsletter.IsPublished,
		Slug:        newsletter.Slug,
		Content:     newsletter.Content,
		Errors:      make(map[string][]string),
	}

	return dashboard.EditNewsletter(formData).Render(renderArgs(c))
}

func (d Dashboard) UpdateNewsletter(c echo.Context) error {
	idParam := c.Param("id")
	newsletterID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid newsletter ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/newsletters")
	}

	type payload struct {
		Title     string `form:"title"`
		Slug      string `form:"slug"`
		Content   string `form:"content"`
		Action    string `form:"action"`
		Published string `form:"published"`
	}

	var newsletterPayload payload
	if err := c.Bind(&newsletterPayload); err != nil {
		return err
	}

	published := newsletterPayload.Published == "on"

	updatePayload := models.UpdateNewsletterPayload{
		ID:          newsletterID,
		UpdatedAt:   time.Now(),
		Title:       newsletterPayload.Title,
		Slug:        newsletterPayload.Slug,
		Content:     newsletterPayload.Content,
		IsPublished: published,
	}

	newsletter, err := models.UpdateNewsletter(
		extractCtx(c),
		d.db.Pool,
		updatePayload,
	)
	if err != nil {
		if validationErrors := extractValidationErrors(err); validationErrors != nil {
			slog.ErrorContext(
				extractCtx(c),
				"could not validate newsletter payload",
				"error",
				err,
			)
		}

		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to update newsletter. Please try again.",
		); err != nil {
			return err
		}

		formData := dashboard.EditNewsletterFormData{
			ID:      newsletterID.String(),
			Title:   newsletterPayload.Title,
			Slug:    newsletterPayload.Slug,
			Content: newsletterPayload.Content,
			Errors:  make(map[string][]string),
		}

		return dashboard.EditNewsletter(formData).Render(renderArgs(c))
	}

	if published && newsletter.ReleasedAt.IsZero() {
		_, err := models.PublishNewsletter(
			c.Request().Context(),
			d.db.Pool,
			models.PublishNewsletterPayload{ID: newsletter.ID, Now: time.Now()},
		)
		if err != nil {
			if err := addFlash(
				c,
				contexts.FlashError,
				"Failed to update newsletter. Please try again.",
			); err != nil {
				return err
			}
			return err
		}
		if err := addFlash(
			c,
			contexts.FlashSuccess,
			"Newsletter published successfully!",
		); err != nil {
			return err
		}

		return c.Redirect(302, "/dashboard/newsletters")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Newsletter updated successfully!",
	); err != nil {
		return err
	}

	return c.Redirect(302, "/dashboard/newsletters")
}

func (d Dashboard) DeleteNewsletter(c echo.Context) error {
	idParam := c.Param("id")
	newsletterID, err := uuid.Parse(idParam)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Invalid newsletter ID.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/newsletters")
	}

	// Check if newsletter exists before deleting
	_, err = models.GetNewsletterByID(
		extractCtx(c),
		d.db.Pool,
		newsletterID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Newsletter not found.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/newsletters")
	}

	// Delete the newsletter
	err = models.DeleteNewsletter(
		extractCtx(c),
		d.db.Pool,
		newsletterID,
	)
	if err != nil {
		if err := addFlash(
			c,
			contexts.FlashError,
			"Failed to delete newsletter. Please try again.",
		); err != nil {
			return err
		}
		return c.Redirect(302, "/dashboard/newsletters")
	}

	if err := addFlash(
		c,
		contexts.FlashSuccess,
		"Newsletter deleted successfully!",
	); err != nil {
		return err
	}

	// Redirect to newsletters
	return c.Redirect(302, "/dashboard/newsletters")
}
