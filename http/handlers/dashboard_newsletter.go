package handlers

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/MBvisti/mortenvistisen/views/emails"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (d Dashboard) NewslettersIndex(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	pageLimit := 7

	offset, currentPage, err := d.base.GetOffsetAndCurrPage(page, pageLimit)
	if err != nil {
		return err
	}

	newsletters, err := d.newsletterSvc.List(
		ctx.Request().Context(),
		int32(pageLimit),
		int32(offset),
	)
	if err != nil {
		return err
	}

	totalNewslettersCount, err := d.newsletterSvc.Count(ctx.Request().Context(), false)
	if err != nil {
		return err
	}

	viewData := make([]dashboard.NewsletterViewData, 0, len(newsletters))
	for _, newsletter := range newsletters {
		viewData = append(viewData, dashboard.NewsletterViewData{
			ID:         newsletter.ID.String(),
			Title:      newsletter.Title,
			Released:   newsletter.Released,
			ReleasedAt: newsletter.ReleasedAt.String(),
			Edition:    strconv.Itoa(int(newsletter.Edition)),
		})
	}

	pagination := components.PaginationProps{
		CurrentPage: currentPage,
		TotalPages:  d.base.CalculateNumberOfPages(int(totalNewslettersCount), 7),
	}

	msg, err := d.base.CookieStore.GetFlashMessages(
		ctx.Request(),
		ctx.Response(),
		"newsletter-released",
	)
	if err != nil {
		return err
	}

	return dashboard.Newsletter(viewData, pagination, csrf.Token(ctx.Request()), len(msg) > 0).
		Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) NewsletterCreate(ctx echo.Context) error {
	releasedNewslettersCount, err := d.newsletterSvc.Count(ctx.Request().Context(), true)
	if err != nil {
		return err
	}

	articles, err := d.articleSvc.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	edition := strconv.Itoa(int(releasedNewslettersCount) + 1)

	newsletter := emails.NewsletterMail{Edition: edition}
	var newsletterPreview bytes.Buffer
	if err := newsletter.Render(ctx.Request().Context(), &newsletterPreview); err != nil {
		return err
	}

	return dashboard.CreateNewsletter(articles, uuid.UUID{}, newsletter, newsletterPreview.String(), csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) NewslettersEdit(ctx echo.Context) error {
	newsletterIDParam := ctx.Param("id")
	newsletterID, err := uuid.Parse(newsletterIDParam)
	if err != nil {
		return err
	}

	newsletter, err := d.newsletterSvc.ByID(ctx.Request().Context(), newsletterID)
	if err != nil {
		return err
	}

	articles, err := d.articleSvc.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	var associatedArticleID uuid.UUID
	for _, article := range articles {
		if article.Slug == newsletter.ArticleSlug {
			associatedArticleID = article.ID
		}
	}

	releasedNewslettersCount, err := d.newsletterSvc.Count(ctx.Request().Context(), true)
	if err != nil {
		return err
	}

	edition := strconv.Itoa(int(releasedNewslettersCount) + 1)

	msg, err := d.base.CookieStore.GetFlashMessages(
		ctx.Request(),
		ctx.Response(),
		"newsletter-draft-saved",
	)
	if err != nil {
		return err
	}

	return dashboard.NewsletterEdit(dashboard.NewsletterEditViewData{
		Title:        newsletter.Title,
		Edition:      edition,
		NewsletterID: newsletterID,
		ArticleID:    associatedArticleID,
		MailPreview: emails.NewsletterMail{
			Title:      newsletter.Title,
			Edition:    edition,
			Paragraphs: newsletter.Paragraphs,
			ArticleLink: d.base.BuildURLFromSlug(
				d.base.FormatArticleSlug(newsletter.ArticleSlug),
			),
		},
		Articles: articles,
	}, csrf.Token(ctx.Request()), len(msg) > 0).
		Render(views.ExtractRenderDeps(ctx))
}

// TODO: implement
func (d Dashboard) NewsletterUpdate(ctx echo.Context) error {
	preview := ctx.QueryParam("preview")
	var updateNewsletterPayload newsletterPayload
	if err := ctx.Bind(&updateNewsletterPayload); err != nil {
		return err
	}
	newsletterIDParam := ctx.Param("id")
	newsletterID, err := uuid.Parse(newsletterIDParam)
	if err != nil {
		return err
	}

	articleID, err := uuid.Parse(updateNewsletterPayload.ArticleID)
	if err != nil {
		return err
	}

	edition, err := strconv.Atoi(updateNewsletterPayload.Edition)
	if err != nil {
		return err
	}

	if preview == "true" {
		return d.previewNewsletter(
			ctx,
			updateNewsletterPayload,
			"put",
			fmt.Sprintf("newsletters/%s/update", newsletterID),
			int32(edition),
		)
	}

	newsletter, err := d.newsletterSvc.Update(
		ctx.Request().Context(),
		updateNewsletterPayload.Title,
		updateNewsletterPayload.Edition,
		updateNewsletterPayload.ParagraphElements,
		updateNewsletterPayload.ArticleID,
		newsletterID,
	)
	if err != nil {
		return err
	}
	// if len(validationErrs) > 0 {
	// 	articles, err := db.QueryPosts(ctx.Request().Context())
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	errors := make(map[string]components.InputError, len(validationErrs))
	//
	// 	for field, validationErr := range validationErrs {
	// 		switch field {
	// 		case "Title":
	// 			errors["title"] = components.InputError{
	// 				Msg:      validationErr,
	// 				OldValue: updateNewsletterPayload.Title,
	// 			}
	// 		case "Paragraphs":
	// 			errors["paragraph-elements"] = components.InputError{
	// 				Msg: validationErr,
	// 			}
	// 		case "ArticleSlug":
	// 			errors["article-id"] = components.InputError{
	// 				Msg: validationErr,
	// 			}
	// 		case "Edition":
	// 			errors["edition"] = components.InputError{
	// 				Msg: validationErr,
	// 			}
	// 		}
	// 	}
	//
	// 	return dashboard.NewsletterPreview(articles, templates.NewsletterMail{
	// 		Title:      updateNewsletterPayload.Title,
	// 		Edition:    updateNewsletterPayload.Edition,
	// 		Paragraphs: updateNewsletterPayload.ParagraphElements,
	// 	}, articleID, errors, csrf.Token(ctx.Request()), "put", fmt.Sprintf("newsletters/%s/update", articleID)).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }

	articles, err := d.articleSvc.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	if updateNewsletterPayload.ReleaseOnCreate == "on" {
		return d.releaseNewsletter(
			ctx,
			newsletter,
			"put",
			fmt.Sprintf("newsletters/%s/update", newsletterID),
		)
	}

	return dashboard.NewsletterEdit(dashboard.NewsletterEditViewData{
		NewsletterID: newsletter.ID,
		Title:        newsletter.Title,
		Edition:      strconv.Itoa(int(newsletter.Edition)),
		ArticleID:    articleID,
		MailPreview: emails.NewsletterMail{
			Title:       newsletter.Title,
			Edition:     strconv.Itoa(int(newsletter.Edition)),
			Paragraphs:  newsletter.Paragraphs,
			ArticleLink: newsletter.ArticleSlug,
		},
		Articles: articles,
	}, csrf.Token(ctx.Request()), false).
		Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) previewNewsletter(ctx echo.Context,
	storeNewsletterPayload newsletterPayload,
	hxAction string,
	endpoint string,
	edition int32,
) error {
	paragraphIndex := ctx.QueryParam("paragraph-index")
	action := ctx.QueryParam("action")

	newsletterPreview, err := d.newsletterSvc.Preview(
		ctx.Request().Context(),
		paragraphIndex,
		action,
		storeNewsletterPayload.Title,
		storeNewsletterPayload.ParagraphElements,
		storeNewsletterPayload.NewParagraphElement,
		storeNewsletterPayload.ArticleID,
		edition,
	)
	if err != nil {
		return err
	}

	articles, err := d.articleSvc.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	var parsedArticleID uuid.UUID
	if storeNewsletterPayload.ArticleID != "" {
		id, err := uuid.Parse(storeNewsletterPayload.ArticleID)
		if err != nil {
			return err
		}

		parsedArticleID = id
	}

	newsletter := emails.NewsletterMail{
		Title:      newsletterPreview.Title,
		Edition:    strconv.Itoa(int(newsletterPreview.Edition)),
		Paragraphs: newsletterPreview.Paragraphs,
	}
	var preview bytes.Buffer
	if err := newsletter.Render(ctx.Request().Context(), &preview); err != nil {
		return err
	}

	return dashboard.NewsletterPreview(articles, newsletter, preview.String(), parsedArticleID, make(map[string]components.InputError), csrf.Token(ctx.Request()), hxAction, endpoint).
		Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) releaseNewsletter(ctx echo.Context,
	newsletter models.Newsletter,
	hxAction string,
	endpoint string,
) error {
	log.Print(hxAction, endpoint)
	_, err := d.newsletterSvc.Release(
		ctx.Request().Context(),
		newsletter,
	)
	if err != nil {
		return err
	}

	// if len(validationErrs) > 0 {
	// 	articles, err := db.QueryPosts(ctx.Request().Context())
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	articleID, err := uuid.Parse(storeNewsletterPayload.ArticleID)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	errors := make(map[string]components.InputError, len(validationErrs))
	//
	// 	for field, validationErr := range validationErrs {
	// 		switch field {
	// 		case "Title":
	// 			errors["title"] = components.InputError{
	// 				Msg:      validationErr,
	// 				OldValue: storeNewsletterPayload.Title,
	// 			}
	// 		case "Paragraphs":
	// 			errors["paragraph-elements"] = components.InputError{
	// 				Msg: validationErr,
	// 			}
	// 		case "ArticleSlug":
	// 			errors["article-id"] = components.InputError{
	// 				Msg: validationErr,
	// 			}
	// 		case "Edition":
	// 			errors["edition"] = components.InputError{
	// 				Msg: validationErr,
	// 			}
	// 		}
	// 	}
	//
	// 	return dashboard.NewsletterPreview(articles, templates.NewsletterMail{
	// 		Title:      storeNewsletterPayload.Title,
	// 		Edition:    storeNewsletterPayload.Edition,
	// 		Paragraphs: storeNewsletterPayload.ParagraphElements,
	// 	}, articleID, errors, csrf.Token(ctx.Request()), hxAction, endpoint).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }
	//

	if err := d.base.CookieStore.CreateFlashMsg(ctx.Request(), ctx.Response(), "newsletter-released"); err != nil {
		return err
	}

	return d.base.RedirectHx(
		ctx.Response().Writer,
		"/dashboard/newsletters",
	)
}

type newsletterPayload struct {
	Title               string   `form:"title"                 validate:"required,gte=4"`
	Edition             string   `form:"edition"               validate:"required"`
	ArticleID           string   `form:"article-id"            validate:"required"`
	NewParagraphElement string   `form:"new-paragraph-element"`
	ParagraphElements   []string `form:"paragraph-element"     validate:"required,gte=1"`
	ReleaseOnCreate     string   `form:"release-on-create"`
}

func (d Dashboard) NewsletterStore(ctx echo.Context) error {
	preview := ctx.QueryParam("preview")

	var storeNewsletterPayload newsletterPayload
	if err := ctx.Bind(&storeNewsletterPayload); err != nil {
		return err
	}

	edition, err := strconv.Atoi(storeNewsletterPayload.Edition)
	if err != nil {
		return err
	}

	if preview == "true" {
		return d.previewNewsletter(
			ctx,
			storeNewsletterPayload,
			"post",
			"newsletters/store",
			int32(edition),
		)
	}

	edtion, err := strconv.Atoi(storeNewsletterPayload.Edition)
	if err != nil {
		return err
	}

	newsletter, err := d.newsletterSvc.CreateDraft(ctx.Request().Context(),
		storeNewsletterPayload.Title,
		int32(edtion),
		storeNewsletterPayload.ParagraphElements,
		storeNewsletterPayload.ArticleID,
	)
	if err != nil {
		return err
	}

	if storeNewsletterPayload.ReleaseOnCreate == "" {
		if err := d.base.CookieStore.CreateFlashMsg(ctx.Request(), ctx.Response(), "newsletter-draft-saved"); err != nil {
			return err
		}
	}

	if storeNewsletterPayload.ReleaseOnCreate == "on" {
		return d.releaseNewsletter(
			ctx,
			newsletter,
			"post",
			"newsletters/store",
		)
	}

	return d.base.RedirectHx(
		ctx.Response().Writer,
		fmt.Sprintf("/dashboard/newsletters/%v/edit", newsletter.ID),
	)
}
