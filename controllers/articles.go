package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/email"

	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/queue/jobs"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v5"
	"github.com/riverqueue/river"
)

type Articles struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
}

func NewArticles(db storage.Pool, insertOnly queue.InsertOnly) Articles {
	return Articles{
		db:         db,
		insertOnly: insertOnly,
	}
}

func (a Articles) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(10)
	if pp := etx.QueryParam("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 &&
			parsed <= 100 {
			perPage = int64(parsed)
		}
	}

	articlesList, err := models.PaginateArticles(
		etx.Request().Context(),
		a.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.ArticleIndex(articlesList))
}

func (a Articles) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	articleID := int32(parsed)

	article, err := models.FindArticle(etx.Request().Context(), a.db.Conn(), articleID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.ArticleShow(article))
}

func (a Articles) New(etx *echo.Context) error {
	tags, err := models.AllTags(etx.Request().Context(), a.db.Conn())
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.ArticleNew(tags, nil))
}

type CreateArticleFormPayload struct {
	Published       bool            `json:"published"`
	Title           string          `json:"title"           validate:"omitempty,max=100"`
	Excerpt         string          `json:"excerpt"         validate:"omitempty,max=255"`
	MetaTitle       string          `json:"metaTitle"       validate:"omitempty,max=100"`
	MetaDescription string          `json:"metaDescription" validate:"omitempty,max=160"`
	ImageLink       string          `json:"imageLink"       validate:"omitempty,url"`
	ReadTime        int32           `json:"readTime"        validate:"gt=0"`
	Content         string          `json:"content"`
	TagSelections   map[string]bool `json:"tagSelections"`
}

func (a Articles) Create(etx *echo.Context) error {
	ctx := etx.Request().Context()

	var payload CreateArticleFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			ctx,
			"could not parse CreateArticleFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.CreateArticleData{
		Published:       payload.Published,
		Title:           payload.Title,
		Excerpt:         payload.Excerpt,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		ImageLink:       payload.ImageLink,
		ReadTime:        payload.ReadTime,
		Content:         payload.Content,
	}

	tx, err := a.db.BeginTx(ctx)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to create article"); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.ArticleNew.URL())
	}
	defer tx.Rollback(ctx)

	article, err := models.CreateArticle(
		ctx,
		tx,
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to create article: %v", err)); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.ArticleNew.URL())
	}

	tagIDs := make([]int32, 0, len(payload.TagSelections))
	for rawID, selected := range payload.TagSelections {
		if !selected {
			continue
		}
		parsedID, err := strconv.ParseInt(rawID, 10, 32)
		if err != nil {
			continue
		}
		tagIDs = append(tagIDs, int32(parsedID))
	}
	if len(tagIDs) > 0 {
		if err := models.AttachTagsToArticle(ctx, tx, article.ID, tagIDs); err != nil {
			if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Article created, but tags could not be associated: %v", err)); flashErr != nil {
				return render(etx, views.InternalError())
			}
			return etx.Redirect(http.StatusSeeOther, routes.ArticleShow.URL(article.ID))
		}
	}

	scheduledJobs := 0
	becameFirstPublished := article.Published && !article.FirstPublishedAt.IsZero()
	if becameFirstPublished {
		scheduledJobs, err = a.scheduleArticleReleaseEmails(ctx, tx, article)
		if err != nil {
			if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to schedule article release emails: %v", err)); flashErr != nil {
				return render(etx, views.InternalError())
			}
			return etx.Redirect(http.StatusSeeOther, routes.ArticleNew.URL())
		}
	}

	if err := a.db.CommitTx(ctx, tx); err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to create article"); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.ArticleNew.URL())
	}

	successMessage := "Article created successfully"
	if scheduledJobs > 0 {
		successMessage = fmt.Sprintf(
			"Article created and %d release emails were scheduled",
			scheduledJobs,
		)
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, successMessage); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ArticleShow.URL(article.ID))
}

func (a Articles) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	articleID := int32(parsed)

	article, err := models.FindArticle(etx.Request().Context(), a.db.Conn(), articleID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	tags, err := models.AllTags(etx.Request().Context(), a.db.Conn())
	if err != nil {
		return render(etx, views.InternalError())
	}

	selectedTagIDsList, err := models.TagIDsForArticle(
		etx.Request().Context(),
		a.db.Conn(),
		articleID,
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	selectedTagIDs := make(map[int32]bool, len(selectedTagIDsList))
	for _, tagID := range selectedTagIDsList {
		selectedTagIDs[tagID] = true
	}

	return render(etx, views.ArticleUpdate(article, tags, selectedTagIDs))
}

type UpdateArticleFormPayload struct {
	Published       bool            `json:"published"`
	Title           string          `json:"title"`
	Excerpt         string          `json:"excerpt"`
	MetaTitle       string          `json:"metaTitle"`
	MetaDescription string          `json:"metaDescription"`
	Slug            string          `json:"slug"`
	ImageLink       string          `json:"imageLink"`
	ReadTime        int32           `json:"readTime"`
	Content         string          `json:"content"`
	TagSelections   map[string]bool `json:"tagSelections"`
}

func (a Articles) Update(etx *echo.Context) error {
	ctx := etx.Request().Context()

	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	articleID := int32(parsed)

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
		ID:              articleID,
		Published:       payload.Published,
		Title:           payload.Title,
		Excerpt:         payload.Excerpt,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		Slug:            payload.Slug,
		ImageLink:       payload.ImageLink,
		ReadTime:        payload.ReadTime,
		Content:         payload.Content,
	}

	tx, err := a.db.BeginTx(ctx)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to update article"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.ArticleEdit.URL(articleID),
		)
	}
	defer tx.Rollback(ctx)

	currentArticle, err := models.FindArticle(ctx, tx, articleID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to load article"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.ArticleEdit.URL(articleID),
		)
	}

	article, err := models.UpdateArticle(
		ctx,
		tx,
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

	tagIDs := make([]int32, 0, len(payload.TagSelections))
	for rawID, selected := range payload.TagSelections {
		if !selected {
			continue
		}
		parsedID, err := strconv.ParseInt(rawID, 10, 32)
		if err != nil {
			continue
		}
		tagIDs = append(tagIDs, int32(parsedID))
	}
	if err := models.ReplaceTagsForArticle(ctx, tx, articleID, tagIDs); err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Article updated, but tags could not be associated: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.ArticleShow.URL(article.ID))
	}

	scheduledJobs := 0
	becameFirstPublished := currentArticle.FirstPublishedAt.IsZero() &&
		!article.FirstPublishedAt.IsZero()
	if becameFirstPublished {
		scheduledJobs, err = a.scheduleArticleReleaseEmails(ctx, tx, article)
		if err != nil {
			if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to schedule article release emails: %v", err)); flashErr != nil {
				return render(etx, views.InternalError())
			}
			return etx.Redirect(
				http.StatusSeeOther,
				routes.ArticleEdit.URL(articleID),
			)
		}
	}

	if err := a.db.CommitTx(ctx, tx); err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to update article"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.ArticleEdit.URL(articleID),
		)
	}

	successMessage := "Article updated successfully"
	if scheduledJobs > 0 {
		successMessage = fmt.Sprintf(
			"Article updated and %d release emails were scheduled",
			scheduledJobs,
		)
	}
	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, successMessage); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ArticleShow.URL(article.ID))
}

func (a Articles) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	articleID := int32(parsed)

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

func (a Articles) scheduleArticleReleaseEmails(
	ctx context.Context,
	tx pgx.Tx,
	article models.Article,
) (int, error) {
	subscribers, err := models.AllSubscribers(ctx, tx)
	if err != nil {
		return 0, err
	}

	sort.Slice(subscribers, func(i, j int) bool {
		return subscribers[i].ID < subscribers[j].ID
	})

	articleEmail := email.NewArticleNotification{
		ArticleTitle: article.Title,
		Summary:      article.Excerpt,
		ArticleURL:   articlePublicURL(article),
	}

	htmlBody, err := articleEmail.ToHTML()
	if err != nil {
		return 0, err
	}

	textBody, err := articleEmail.ToText()
	if err != nil {
		return 0, err
	}

	scheduleBase := time.Now().UTC().Add(10 * time.Second).Truncate(time.Second)

	var insertParams []river.InsertManyParams
	sendIndex := 0
	for _, subscriber := range subscribers {
		if !subscriber.IsVerified {
			continue
		}

		emailAddress := strings.TrimSpace(subscriber.Email)
		if emailAddress == "" {
			continue
		}

		scheduledAt := scheduledNewsletterSendTime(scheduleBase, sendIndex)
		insertParams = append(insertParams, river.InsertManyParams{
			Args: jobs.SendMarketingEmailArgs{
				Data: email.MarketingData{
					To:             []string{emailAddress},
					From:           "newsletter@mortenvistisen.com",
					Subject:        article.Title,
					HTMLBody:       htmlBody,
					TextBody:       textBody,
					UnsubscribeURL: newsletterUnsubscribeURL(),
					Tags:           []string{"article_release"},
					Metadata: map[string]string{
						"article_id":    strconv.Itoa(int(article.ID)),
						"subscriber_id": strconv.Itoa(int(subscriber.ID)),
					},
				},
			},
			InsertOpts: &river.InsertOpts{
				ScheduledAt: scheduledAt,
			},
		})

		sendIndex++
	}

	if len(insertParams) == 0 {
		return 0, nil
	}

	if _, err := a.insertOnly.InsertManyTx(ctx, tx, insertParams); err != nil {
		return 0, err
	}

	return len(insertParams), nil
}

func articlePublicURL(article models.Article) string {
	baseURL := strings.TrimRight(config.BaseURL, "/")
	if article.Slug != "" {
		return fmt.Sprintf("%s%s", baseURL, routes.Article.URL(article.Slug))
	}

	return fmt.Sprintf("%s%s", baseURL, routes.ArticleShow.URL(article.ID))
}
