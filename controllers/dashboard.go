package controllers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/MBvisti/mortenvistisen/entity"
	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/MBvisti/mortenvistisen/views/validation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (c *Controller) DashboardIndex(ctx echo.Context) error {
	return dashboard.Index().Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardResendVerificationMail(ctx echo.Context) error {
	subscriberID := ctx.Param("id")

	subscriberUUID, err := uuid.Parse(subscriberID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	subscriber, err := c.db.QuerySubscriber(ctx.Request().Context(), subscriberUUID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if subscriber.IsVerified.Bool {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if err := c.db.DeleteSubscriberTokenBySubscriberID(ctx.Request().Context(), subscriberUUID); err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	generatedTkn, err := c.tknManager.GenerateToken()
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	activationToken := tokens.CreateActivationToken(
		generatedTkn.PlainTextToken,
		generatedTkn.HashedToken,
	)

	if err := c.db.InsertSubscriberToken(ctx.Request().Context(), database.InsertSubscriberTokenParams{
		ID:           uuid.New(),
		CreatedAt:    database.ConvertToPGTimestamptz(time.Now()),
		Hash:         activationToken.Hash,
		ExpiresAt:    database.ConvertToPGTimestamptz(activationToken.GetExpirationTime()),
		Scope:        activationToken.GetScope(),
		SubscriberID: subscriberUUID,
	}); err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	newsletterMail := templates.NewsletterWelcomeMail{
		ConfirmationLink: fmt.Sprintf(
			"%s://%s/verify-subscriber?token=%s",
			c.cfg.App.AppScheme,
			c.cfg.App.AppHost,
			activationToken.GetPlainText(),
		),
		UnsubscribeLink: "",
	}
	textVersion, err := newsletterMail.GenerateTextVersion()
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	_, err = c.queueClient.Insert(ctx.Request().Context(), queue.EmailJobArgs{
		To:          subscriber.Email.String,
		From:        "noreply@mortenvistisen.com",
		Subject:     "Thanks for signing up!",
		TextVersion: textVersion,
		HtmlVersion: htmlVersion,
	}, nil)
	if err != nil {
		return err
	}

	return dashboard.SuccessMsg("Verification mail send").Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardSubscribers(ctx echo.Context) error {
	subs, err := c.db.QueryAllSubscribers(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.SubscriberViewData, 0, len(subs))
	for _, sub := range subs {
		viewData = append(viewData, dashboard.SubscriberViewData{
			Email:        sub.Email.String,
			ID:           sub.ID.String(),
			Verified:     sub.IsVerified.Bool,
			SubscribedAt: sub.SubscribedAt.Time.String(),
			Refererer:    sub.Referer.String,
		})
	}

	return dashboard.Subscribers(viewData, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticles(ctx echo.Context) error {
	articles, err := c.db.QueryAllPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.ArticleViewData, 0, len(articles))
	for _, article := range articles {
		viewData = append(viewData, dashboard.ArticleViewData{
			ID:         article.ID.String(),
			Title:      article.Title,
			Draft:      article.Draft,
			ReleasedAt: article.ReleasedAt.Time.String(),
			Slug:       article.Slug,
		})
	}

	return dashboard.Articles(viewData, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticleCreate(ctx echo.Context) error {
	tags, err := c.db.QueryAllTags(ctx.Request().Context())
	if err != nil {
		return err
	}

	var keywords []dashboard.Keyword
	for _, tag := range tags {
		keywords = append(keywords, dashboard.Keyword{
			ID:    tag.ID.String(),
			Value: tag.Name,
		})
	}

	filenames, err := posts.GetAllFiles()
	if err != nil {
		return err
	}

	usedFilenames, err := c.db.QueryAllFilenames(ctx.Request().Context())
	if err != nil {
		return err
	}

	var unusedFileNames []string
	for _, filename := range filenames {
		if !slices.Contains(usedFilenames, filename) {
			unusedFileNames = append(unusedFileNames, filename)
		}
	}

	return dashboard.CreateArticle(
		dashboard.CreateArticleViewData{
			Keywords:  keywords,
			Filenames: unusedFileNames,
		},
		csrf.Token(ctx.Request())).Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticleEdit(ctx echo.Context) error {
	slug := ctx.Param("slug")

	post, err := c.db.GetPostBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return err
	}

	postContent, err := c.postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	tags, err := c.db.GetTagsForPost(ctx.Request().Context(), post.ID)
	if err != nil {
		return err
	}

	var keywords string
	for i, kw := range tags {
		if i == len(tags)-1 {
			keywords = keywords + kw.Name
		} else {
			keywords = keywords + kw.Name + ", "
		}
	}

	return dashboard.ArticleEdit(dashboard.ArticleEditViewData{
		ID:          post.ID.String(),
		CreatedAt:   post.CreatedAt.Time,
		UpdatedAt:   post.UpdatedAt.Time,
		Title:       post.Title,
		HeaderTitle: post.HeaderTitle.String,
		Filename:    post.Filename,
		Slug:        post.Slug,
		Excerpt:     post.Excerpt,
		Draft:       post.Draft,
		ReleasedAt:  post.ReleasedAt.Time,
		ReadTime:    post.ReadTime.Int32,
		Content:     postContent,
		Keywords:    keywords,
	}, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DeleteSubscriber(ctx echo.Context) error {
	id := ctx.Param("ID")

	parsedID := uuid.MustParse(id)

	if err := c.db.DeleteSubscriber(ctx.Request().Context(), parsedID); err != nil {
		return err
	}

	ctx.Response().Header().Add("HX-Refresh", "true")

	return ctx.Redirect(http.StatusSeeOther, "/dashboard/subscribers")
}

type NewTagFormPayload struct {
	Name             string   `form:"tag-name"`
	SelectedKeywords []string `form:"selected-keyword"`
}

func (c *Controller) DashboadTagStore(ctx echo.Context) error {
	var tagPayload NewTagFormPayload
	if err := ctx.Bind(&tagPayload); err != nil {
		return err
	}

	if err := c.db.InsertTag(ctx.Request().Context(), database.InsertTagParams{
		ID:   uuid.New(),
		Name: tagPayload.Name,
	}); err != nil {
		return err
	}

	tags, err := c.db.QueryAllTags(ctx.Request().Context())
	if err != nil {
		return err
	}

	var keywords []dashboard.Keyword
	for _, tag := range tags {
		selected := false
		for _, selectedKW := range tagPayload.SelectedKeywords {
			if selectedKW == tag.ID.String() {
				selected = true
			}
		}

		kw := dashboard.Keyword{
			ID:       tag.ID.String(),
			Value:    tag.Name,
			Selected: selected,
		}
		keywords = append(keywords, kw)
	}

	return dashboard.KeywordsGrid(keywords).Render(views.ExtractRenderDeps(ctx))
}

type NewPostFormPayload struct {
	Title             string   `form:"title"`
	HeaderTitle       string   `form:"header-title"`
	Excerpt           string   `form:"excerpt"`
	EstimatedReadTime string   `form:"estimated-read-time"`
	Filename          string   `form:"filename"`
	SelectedKeywords  []string `form:"selected-keyword"`
	Release           string   `form:"release"`
}

func (c *Controller) DashboadPostStore(ctx echo.Context) error {
	var postPayload NewPostFormPayload
	if err := ctx.Bind(&postPayload); err != nil {
		return err
	}

	var releaseNow bool
	if postPayload.Release == "on" {
		releaseNow = true
	}

	estimatedReadTime, err := strconv.Atoi(postPayload.EstimatedReadTime)
	if err != nil {
		return err
	}

	if err := services.NewPost(ctx.Request().Context(), &c.db, c.validate, entity.NewPost{
		Title:             postPayload.Title,
		HeaderTitle:       postPayload.HeaderTitle,
		Excerpt:           postPayload.Excerpt,
		ReleaseNow:        releaseNow,
		EstimatedReadTime: int64(estimatedReadTime),
		Filename:          postPayload.Filename,
	}, postPayload.SelectedKeywords); err != nil {
		e, ok := err.(validator.ValidationErrors)
		if !ok {
			return c.InternalError(ctx)
		}

		tags, err := c.db.QueryAllTags(ctx.Request().Context())
		if err != nil {
			return err
		}

		var keywords []dashboard.Keyword
		for _, tag := range tags {
			keywords = append(keywords, dashboard.Keyword{
				ID:    tag.ID.String(),
				Value: tag.Name,
			})
		}

		filenames, err := posts.GetAllFiles()
		if err != nil {
			return err
		}

		usedFilenames, err := c.db.QueryAllFilenames(ctx.Request().Context())
		if err != nil {
			return err
		}

		var unusedFileNames []string
		for _, filename := range filenames {
			if !slices.Contains(usedFilenames, filename) {
				unusedFileNames = append(unusedFileNames, filename)
			}
		}

		props := dashboard.CreateArticleFormProps{
			Title: validation.InputField{
				OldValue: postPayload.Title,
			},
			HeaderTitle: validation.InputField{
				OldValue: postPayload.HeaderTitle,
			},
			Excerpt: validation.InputField{
				OldValue: postPayload.Excerpt,
			},
			Filename: validation.InputField{
				OldValue: postPayload.Filename,
			},
			EstimatedReadTime: validation.InputField{
				OldValue: postPayload.EstimatedReadTime,
			},
			ReleaseNow: releaseNow,
		}

		for _, validationError := range e {
			switch validationError.StructField() {
			case "Title":
				props.Title.Invalid = true
				props.Title.InvalidMsg = "Title is not long enough"
			case "HeaderTitle":
				props.HeaderTitle.Invalid = true
				props.HeaderTitle.InvalidMsg = "Header title is not long enough"
			case "Excerpt":
				props.Excerpt.Invalid = true
				props.Excerpt.InvalidMsg = "Excerpt has to be between 130 and 160 chars long"
			case "EstimatedReadTime":
				props.EstimatedReadTime.Invalid = true
				props.EstimatedReadTime.InvalidMsg = "Est. time cannot be 0"
			case "Filename":
				props.Filename.Invalid = true
				props.Filename.InvalidMsg = "All filenames must end with .md"
			}
		}

		return dashboard.CreateArticleFormContent(
			props,
			keywords,
			unusedFileNames,
		).Render(views.ExtractRenderDeps(ctx))
	}

	ctx.Response().Writer.Header().Add("HX-Redirect", "/dashboard/articles")
	return nil
}
