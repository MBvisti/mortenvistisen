package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
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
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/MBvisti/mortenvistisen/views/validation"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgtype"
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

	currentPage := 1
	pagination := components.PaginationPayload{
		CurrentPage:     currentPage,
		NextPage:        currentPage + 1,
		PrevPage:        currentPage - 1,
		HasNextNextPage: len(viewData)-7 >= 7,
	}

	if len(viewData) < 7 {
		pagination.NoNextPage = true
	}

	return dashboard.Subscribers(viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticles(ctx echo.Context) error {
	page := ctx.QueryParam("page")

	var currentPage int
	if page == "" {
		currentPage = 1
	}
	if page != "" {
		cp, err := strconv.Atoi(page)
		if err != nil {
			return err
		}

		currentPage = cp
	}

	offset := 0
	if currentPage == 2 {
		offset = 7
	}

	if currentPage > 2 {
		offset = 7 * (currentPage - 1)
	}

	articles, err := c.db.QueryAllPosts(ctx.Request().Context(), int32(offset))
	if err != nil {
		return err
	}

	var totalPostCount int

	viewData := make([]dashboard.ArticleViewData, 0, len(articles))
	for i, article := range articles {
		if i == 0 {
			totalPostCount = int(article.TotalPostsCount)
		}
		viewData = append(viewData, dashboard.ArticleViewData{
			ID:         article.ID.String(),
			Title:      article.Title,
			Draft:      article.Draft,
			ReleasedAt: article.ReleasedAt.Time.String(),
			Slug:       article.Slug,
		})
	}

	pagination := components.PaginationPayload{
		CurrentPage:     currentPage,
		NextPage:        currentPage + 1,
		PrevPage:        currentPage - 1,
		HasNextNextPage: totalPostCount-7 >= 7,
	}

	if len(articles) < 7 {
		pagination.NoNextPage = true
	}

	return dashboard.Articles(viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardNewsletter(ctx echo.Context) error {
	page := ctx.QueryParam("page")

	var currentPage int
	if page == "" {
		currentPage = 1
	}
	if page != "" {
		cp, err := strconv.Atoi(page)
		if err != nil {
			return err
		}

		currentPage = cp
	}

	offset := 0
	if currentPage == 2 {
		offset = 7
	}

	if currentPage > 2 {
		offset = 7 * (currentPage - 1)
	}

	newsletters, err := c.db.QueryNewsletterInPages(ctx.Request().Context(), int32(offset))
	if err != nil {
		return err
	}

	releasedNewslettersCount, err := c.db.QueryReleasedNewslettersCount(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.NewsletterViewData, 0, len(newsletters))
	for _, newsletter := range newsletters {
		viewData = append(viewData, dashboard.NewsletterViewData{
			ID:         newsletter.ID.String(),
			Title:      newsletter.Title,
			Released:   newsletter.Released.Bool,
			ReleasedAt: newsletter.ReleasedAt.Time.String(),
			Edition:    strconv.Itoa(int(newsletter.Edition.Int32)),
		})
	}

	pagination := components.PaginationPayload{
		CurrentPage:     currentPage,
		NextPage:        currentPage + 1,
		PrevPage:        currentPage - 1,
		HasNextNextPage: releasedNewslettersCount-7 >= 7,
	}

	if len(newsletters) < 7 {
		pagination.NoNextPage = true
	}

	return dashboard.Newsletter(viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardNewsletterCreate(ctx echo.Context) error {
	releasedNewslettersCount, err := c.db.QueryReleasedNewslettersCount(ctx.Request().Context())
	if err != nil {
		return err
	}

	articles, err := c.db.QueryPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	edition := strconv.Itoa(int(releasedNewslettersCount) + 1)

	return dashboard.CreateNewsletter(articles, uuid.UUID{}, templates.NewsletterMail{
		Edition: edition,
	}, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardNewsletterEdit(ctx echo.Context) error {
	newsletterIDParam := ctx.Param("id")
	newsletterID, err := uuid.Parse(newsletterIDParam)
	if err != nil {
		return err
	}

	newsletter, err := c.db.QueryNewsletterByID(ctx.Request().Context(), newsletterID)
	if err != nil {
		return err
	}

	articles, err := c.db.QueryPosts(ctx.Request().Context())
	if err != nil {
		return err
	}
	edition := strconv.Itoa(int(newsletter.Edition.Int32))
	return dashboard.NewsletterEdit(dashboard.NewsletterEditViewData{
		Title:     newsletter.Title,
		Edition:   edition,
		ArticleID: newsletter.AssociatedArticleID,
		MailPreview: templates.NewsletterMail{
			Title:       newsletter.Title,
			Edition:     edition,
			Paragraphs:  []string{""},
			ArticleLink: "",
		},
		Articles: articles,
	}, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

type NewsletterCreatePreview struct {
	Title               string    `form:"title"`
	ParagraphElements   []string  `form:"paragraph-element"`
	NewParagraphElement string    `form:"new-paragraph-element"`
	ArticleID           uuid.UUID `form:"article-id"`
}

func (c *Controller) dashboardNewsletterCreatePreview(ctx echo.Context) error {
	paragraphIndex := ctx.QueryParam("paragraph-index")
	action := ctx.QueryParam("action")

	var previewPayload NewsletterCreatePreview
	if err := ctx.Bind(&previewPayload); err != nil {
		return err
	}

	paras := previewPayload.ParagraphElements

	if paragraphIndex != "" && action == "del" {
		log.Print(paras)
		index, err := strconv.Atoi(paragraphIndex)
		if err != nil {
			return err
		}

		if action == "del" {
			paras = append(
				paras[:index],
				paras[index+1:]...)
		}
		log.Print(paras)
	}

	if previewPayload.NewParagraphElement != "" && action != "del" {
		paras = append(
			paras,
			previewPayload.NewParagraphElement,
		)
	}

	articles, err := c.db.QueryPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	releasedNewslettersCount, err := c.db.QueryReleasedNewslettersCount(ctx.Request().Context())
	(ctx.Request().Context())
	if err != nil {
		return err
	}

	edition := strconv.Itoa(int(releasedNewslettersCount) + 1)

	return dashboard.NewsletterPreview(articles, templates.NewsletterMail{
		Title:      previewPayload.Title,
		Edition:    edition,
		Paragraphs: paras,
	}, previewPayload.ArticleID, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

type StoreNewsletterPayload struct {
	Title             string    `form:"title"             validate:"required,gte=4"`
	Edition           string    `form:"edition"           validate:"required"`
	ArticleID         uuid.UUID `form:"article-id"        validate:"required"`
	ParagraphElements []string  `form:"paragraph-element" validate:"required,gte=1"`
	ReleaseOnCreate   string    `form:"release-on-create"`
}

func (c *Controller) DashboardNewsletterStore(ctx echo.Context) error {
	preview := ctx.QueryParam("preview")
	if preview == "true" {
		return c.dashboardNewsletterCreatePreview(ctx)
	}

	var storeNewsletterPayload StoreNewsletterPayload
	if err := ctx.Bind(&storeNewsletterPayload); err != nil {
		return err
	}

	if err := c.validate.Struct(storeNewsletterPayload); err != nil {
		log.Print(err)
		return err
	}

	edition, err := strconv.Atoi(storeNewsletterPayload.Edition)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(storeNewsletterPayload.ParagraphElements); err != nil {
		return err
	}

	now := database.ConvertToPGTimestamptz(time.Now())
	insertArgs := database.InsertNewsletterParams{
		ID:                  uuid.New(),
		CreatedAt:           now,
		UpdatedAt:           now,
		Title:               storeNewsletterPayload.Title,
		Edition:             sql.NullInt32{Int32: int32(edition), Valid: true},
		Body:                buf.Bytes(),
		AssociatedArticleID: storeNewsletterPayload.ArticleID,
	}

	if storeNewsletterPayload.ReleaseOnCreate == "on" {
		insertArgs.Released = pgtype.Bool{Bool: true, Valid: true}
		insertArgs.ReleasedAt = now
	}

	newsletter, err := c.db.InsertNewsletter(
		ctx.Request().Context(), insertArgs,
	)
	if err != nil {
		return err
	}

	if newsletter.Released.Bool {
		verifiedSubs, err := c.db.QueryVerifiedSubscribers(ctx.Request().Context())
		if err != nil {
			return err
		}

		article, err := c.db.QueryPostByID(ctx.Request().Context(), newsletter.AssociatedArticleID)
		if err != nil {
			return err
		}

		newsletterMail := templates.NewsletterMail{
			Title:       newsletter.Title,
			Edition:     strconv.Itoa(int(newsletter.Edition.Int32)),
			Paragraphs:  storeNewsletterPayload.ParagraphElements,
			ArticleLink: c.buildURLFromSlug(c.formatArticleSlug(article.Slug)),
		}

		htmlMail, err := newsletterMail.GenerateHtmlVersion()
		if err != nil {
			return err
		}

		textMail, err := newsletterMail.GenerateTextVersion()
		if err != nil {
			return err
		}

		for _, verifiedSub := range verifiedSubs {
			if err := c.mail.Send(
				ctx.Request().Context(),
				verifiedSub.Email.String,
				"newsletter@mortenvistisen.com",
				fmt.Sprintf("MBV newsletter edition: %v", newsletter.Edition.Int32),
				textMail,
				htmlMail,
			); err != nil {
				slog.Error("could not send email", "error", err)
				return err
			}
		}

		ctx.Response().Header().Add("HX-Refresh", "true")
		ctx.Response().Writer.Header().Add("HX-Redirect", "/dashboard/newsletter")
		return ctx.String(http.StatusCreated, "newsletter sent")
	}

	ctx.Response().Header().Add("HX-Refresh", "true")
	ctx.Response().
		Writer.Header().
		Add("HX-Redirect", fmt.Sprintf("/dashboard/newsletter/%v/edit", newsletter.ID))
	return ctx.String(http.StatusCreated, "newsletter sent")
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

type UpdateArticlePayload struct {
	Title       string `form:"title"`
	HeaderTitle string `form:"header-title"`
	Slug        string `form:"slug"`
	Excerpt     string `form:"excerpt"`
	EstReadTime int32  `form:"est-read-time"`
	UpdatedAt   string `form:"updated-at"`
	ReleasedAt  string `form:"released-at"`
	IsLive      string `form:"is-live"`
}

func (c *Controller) DashboardArticleUpdate(ctx echo.Context) error {
	id := ctx.Param("id")

	var updateArticlePayload UpdateArticlePayload
	if err := ctx.Bind(&updateArticlePayload); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	parsedReleasedAt := carbon.Parse(updateArticlePayload.ReleasedAt).ToStdTime()

	var release bool
	if updateArticlePayload.IsLive == "on" {
		release = true
	}

	post, err := services.UpdatePost(ctx.Request().Context(), &c.db, c.validate, entity.UpdatePost{
		ID:                parsedID,
		Title:             updateArticlePayload.Title,
		HeaderTitle:       updateArticlePayload.HeaderTitle,
		Excerpt:           updateArticlePayload.Excerpt,
		Slug:              updateArticlePayload.Slug,
		ReleaedAt:         parsedReleasedAt,
		ReleaseNow:        release,
		EstimatedReadTime: updateArticlePayload.EstReadTime,
	})
	if err != nil {
		return err
	}

	return dashboard.ArticleEditForm(dashboard.ArticleEditViewData{
		ID:          post.ID.String(),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		HeaderTitle: post.HeaderTitle,
		Filename:    post.Filename,
		Slug:        post.Slug,
		Excerpt:     post.Excerpt,
		Draft:       post.Draft,
		ReleasedAt:  post.ReleaseDate,
		ReadTime:    post.ReadTime,
	}, true).
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
