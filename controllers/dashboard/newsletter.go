package dashboard

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/MBvisti/mortenvistisen/controllers/internal/utilities"
	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func NewslettersIndex(ctx echo.Context, db database.Queries) error {
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

	newsletters, err := db.QueryNewsletterInPages(ctx.Request().Context(), int32(offset))
	if err != nil {
		return err
	}

	releasedNewslettersCount, err := db.QueryReleasedNewslettersCount(
		ctx.Request().Context(),
	)
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

func NewsletterCreate(ctx echo.Context, db database.Queries) error {
	releasedNewslettersCount, err := db.QueryReleasedNewslettersCount(
		ctx.Request().Context(),
	)
	if err != nil {
		return err
	}

	articles, err := db.QueryPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	edition := strconv.Itoa(int(releasedNewslettersCount) + 1)

	return dashboard.CreateNewsletter(articles, uuid.UUID{}, templates.NewsletterMail{
		Edition: edition,
	}, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func NewsletterEdit(ctx echo.Context, db database.Queries) error {
	newsletterIDParam := ctx.Param("id")
	newsletterID, err := uuid.Parse(newsletterIDParam)
	if err != nil {
		return err
	}

	newsletter, err := db.QueryNewsletterByID(ctx.Request().Context(), newsletterID)
	if err != nil {
		return err
	}

	articles, err := db.QueryPosts(ctx.Request().Context())
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

// TODO: implement
func NewsletterUpdate(ctx echo.Context, db database.Queries) error {
	newsletterIDParam := ctx.Param("id")
	newsletterID, err := uuid.Parse(newsletterIDParam)
	if err != nil {
		return err
	}

	newsletter, err := db.QueryNewsletterByID(ctx.Request().Context(), newsletterID)
	if err != nil {
		return err
	}

	articles, err := db.QueryPosts(ctx.Request().Context())
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

type newsletterCreatePreviewPayload struct {
	Title               string    `form:"title"`
	ParagraphElements   []string  `form:"paragraph-element"`
	NewParagraphElement string    `form:"new-paragraph-element"`
	ArticleID           uuid.UUID `form:"article-id"`
}

func newsletterCreatePreview(ctx echo.Context, db database.Queries) error {
	paragraphIndex := ctx.QueryParam("paragraph-index")
	action := ctx.QueryParam("action")

	var previewPayload newsletterCreatePreviewPayload
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

	articles, err := db.QueryPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	releasedNewslettersCount, err := db.QueryReleasedNewslettersCount(
		ctx.Request().Context(),
	)
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

type newsletterStorePayload struct {
	Title             string    `form:"title"             validate:"required,gte=4"`
	Edition           string    `form:"edition"           validate:"required"`
	ArticleID         uuid.UUID `form:"article-id"        validate:"required"`
	ParagraphElements []string  `form:"paragraph-element" validate:"required,gte=1"`
	ReleaseOnCreate   string    `form:"release-on-create"`
}

func NewsletterStore(
	ctx echo.Context,
	v *validator.Validate,
	db database.Queries,
	mail mail.Mail,
) error {
	preview := ctx.QueryParam("preview")
	if preview == "true" {
		return newsletterCreatePreview(ctx, db)
	}

	var storeNewsletterPayload newsletterStorePayload
	if err := ctx.Bind(&storeNewsletterPayload); err != nil {
		return err
	}

	if err := v.Struct(storeNewsletterPayload); err != nil {
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

	newsletter, err := db.InsertNewsletter(
		ctx.Request().Context(), insertArgs,
	)
	if err != nil {
		return err
	}

	if newsletter.Released.Bool {
		verifiedSubs, err := db.QueryVerifiedSubscribers(ctx.Request().Context())
		if err != nil {
			return err
		}

		article, err := db.QueryPostByID(
			ctx.Request().Context(),
			newsletter.AssociatedArticleID,
		)
		if err != nil {
			return err
		}

		newsletterMail := templates.NewsletterMail{
			Title:       newsletter.Title,
			Edition:     strconv.Itoa(int(newsletter.Edition.Int32)),
			Paragraphs:  storeNewsletterPayload.ParagraphElements,
			ArticleLink: utilities.BuildURLFromSlug(utilities.FormatArticleSlug(article.Slug)),
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
			if err := mail.Send(
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
