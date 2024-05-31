package dashboard

import (
	"fmt"
	"strconv"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/usecases"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func NewslettersIndex(
	ctx echo.Context,
	db database.Queries,
	sess *sessions.CookieStore,
) error {
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

	s, err := sess.Get(ctx.Request(), "flashMsg")
	if err != nil {
		return err
	}

	var showFlash bool
	for _, flash := range s.Flashes() {
		f, ok := flash.(string)
		if !ok {
			return err
		}

		if f == "newsletter-released" {
			showFlash = true
		}
	}

	if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return dashboard.Newsletter(viewData, pagination, csrf.Token(ctx.Request()), showFlash).
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

func NewslettersEdit(
	ctx echo.Context,
	db database.Queries,
	newsletterUsecase usecases.Newsletter,
	sess *sessions.CookieStore,
) error {
	newsletterIDParam := ctx.Param("id")
	newsletterID, err := uuid.Parse(newsletterIDParam)
	if err != nil {
		return err
	}

	newsletter, err := newsletterUsecase.Get(ctx.Request().Context(), newsletterID)
	if err != nil {
		return err
	}

	articles, err := db.QueryPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	var associatedArticleID uuid.UUID
	for _, article := range articles {
		if article.Slug == newsletter.ArticleSlug {
			associatedArticleID = article.ID
		}
	}

	releasedNewslettersCount, err := db.QueryReleasedNewslettersCount(
		ctx.Request().Context(),
	)
	if err != nil {
		return err
	}
	edition := strconv.Itoa(int(releasedNewslettersCount) + 1)

	s, err := sess.Get(ctx.Request(), "flashMsg")
	if err != nil {
		return err
	}

	var showFlash bool
	for _, flash := range s.Flashes() {
		f, ok := flash.(string)
		if !ok {
			return err
		}

		if f == "newsletter-draft-saved" {
			showFlash = true
		}
	}

	if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return dashboard.NewsletterEdit(dashboard.NewsletterEditViewData{
		Title:        newsletter.Title,
		Edition:      edition,
		NewsletterID: newsletterID,
		ArticleID:    associatedArticleID,
		MailPreview: templates.NewsletterMail{
			Title:      newsletter.Title,
			Edition:    edition,
			Paragraphs: newsletter.Paragraphs,
			ArticleLink: usecases.BuildURLFromSlug(
				usecases.FormatArticleSlug(newsletter.ArticleSlug),
			),
		},
		Articles: articles,
	}, csrf.Token(ctx.Request()), showFlash).
		Render(views.ExtractRenderDeps(ctx))
}

// TODO: implement
func NewsletterUpdate(
	ctx echo.Context,
	db database.Queries,
	newsletterUsecase usecases.Newsletter,
	sess *sessions.CookieStore,
) error {
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

	if preview == "true" {
		return previewNewsletter(
			ctx,
			db,
			newsletterUsecase,
			updateNewsletterPayload,
			"put",
			fmt.Sprintf("newsletters/%s/update", newsletterID),
		)
	}

	newsletter, validationErrs, err := newsletterUsecase.Update(
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
	if len(validationErrs) > 0 {
		articles, err := db.QueryPosts(ctx.Request().Context())
		if err != nil {
			return err
		}

		errors := make(map[string]components.InputError, len(validationErrs))

		for field, validationErr := range validationErrs {
			switch field {
			case "Title":
				errors["title"] = components.InputError{
					Msg:      validationErr,
					OldValue: updateNewsletterPayload.Title,
				}
			case "Paragraphs":
				errors["paragraph-elements"] = components.InputError{
					Msg: validationErr,
				}
			case "ArticleSlug":
				errors["article-id"] = components.InputError{
					Msg: validationErr,
				}
			case "Edition":
				errors["edition"] = components.InputError{
					Msg: validationErr,
				}
			}
		}

		return dashboard.NewsletterPreview(articles, templates.NewsletterMail{
			Title:      updateNewsletterPayload.Title,
			Edition:    updateNewsletterPayload.Edition,
			Paragraphs: updateNewsletterPayload.ParagraphElements,
		}, articleID, errors, csrf.Token(ctx.Request()), "put", fmt.Sprintf("newsletters/%s/update", articleID)).
			Render(views.ExtractRenderDeps(ctx))
	}

	articles, err := db.QueryPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	if updateNewsletterPayload.ReleaseOnCreate == "on" {
		return releaseNewsletter(
			ctx,
			newsletter,
			db,
			newsletterUsecase,
			updateNewsletterPayload,
			sess,
			"put",
			fmt.Sprintf("newsletters/%s/update", newsletterID),
		)
	}

	return dashboard.NewsletterEdit(dashboard.NewsletterEditViewData{
		NewsletterID: newsletter.ID,
		Title:        newsletter.Title,
		Edition:      strconv.Itoa(int(newsletter.Edition)),
		ArticleID:    articleID,
		MailPreview: templates.NewsletterMail{
			Title:       newsletter.Title,
			Edition:     strconv.Itoa(int(newsletter.Edition)),
			Paragraphs:  newsletter.Paragraphs,
			ArticleLink: newsletter.ArticleSlug,
		},
		Articles: articles,
	}, csrf.Token(ctx.Request()), false).
		Render(views.ExtractRenderDeps(ctx))
}

func previewNewsletter(ctx echo.Context,
	db database.Queries,
	newsletterUsecase usecases.Newsletter,
	storeNewsletterPayload newsletterPayload,
	hxAction string,
	endpoint string,
) error {
	paragraphIndex := ctx.QueryParam("paragraph-index")
	action := ctx.QueryParam("action")

	newsletterPreview, err := newsletterUsecase.Preview(
		ctx.Request().Context(),
		paragraphIndex,
		action,
		storeNewsletterPayload.Title,
		storeNewsletterPayload.ParagraphElements,
		storeNewsletterPayload.NewParagraphElement,
		storeNewsletterPayload.ArticleID,
	)
	if err != nil {
		return err
	}

	articles, err := db.QueryPosts(ctx.Request().Context())
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

	return dashboard.NewsletterPreview(articles, templates.NewsletterMail{
		Title:      newsletterPreview.Title,
		Edition:    strconv.Itoa(int(newsletterPreview.Edition)),
		Paragraphs: newsletterPreview.Paragraphs,
	}, parsedArticleID, make(map[string]components.InputError), csrf.Token(ctx.Request()), hxAction, endpoint).Render(views.ExtractRenderDeps(ctx))
}

func releaseNewsletter(ctx echo.Context,
	newsletter domain.Newsletter,
	db database.Queries,
	newsletterUsecase usecases.Newsletter,
	storeNewsletterPayload newsletterPayload,
	sess *sessions.CookieStore,
	hxAction string,
	endpoint string,
) error {
	validationErrs, err := newsletterUsecase.ReleaseNewsletter(
		ctx.Request().Context(),
		newsletter,
	)
	if err != nil {
		return err
	}

	if len(validationErrs) > 0 {
		articles, err := db.QueryPosts(ctx.Request().Context())
		if err != nil {
			return err
		}

		articleID, err := uuid.Parse(storeNewsletterPayload.ArticleID)
		if err != nil {
			return err
		}

		errors := make(map[string]components.InputError, len(validationErrs))

		for field, validationErr := range validationErrs {
			switch field {
			case "Title":
				errors["title"] = components.InputError{
					Msg:      validationErr,
					OldValue: storeNewsletterPayload.Title,
				}
			case "Paragraphs":
				errors["paragraph-elements"] = components.InputError{
					Msg: validationErr,
				}
			case "ArticleSlug":
				errors["article-id"] = components.InputError{
					Msg: validationErr,
				}
			case "Edition":
				errors["edition"] = components.InputError{
					Msg: validationErr,
				}
			}
		}

		return dashboard.NewsletterPreview(articles, templates.NewsletterMail{
			Title:      storeNewsletterPayload.Title,
			Edition:    storeNewsletterPayload.Edition,
			Paragraphs: storeNewsletterPayload.ParagraphElements,
		}, articleID, errors, csrf.Token(ctx.Request()), hxAction, endpoint).
			Render(views.ExtractRenderDeps(ctx))
	}

	s, err := sess.Get(ctx.Request(), "flashMsg")
	if err != nil {
		return err
	}
	s.AddFlash("newsletter-released")
	if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return controllers.RedirectHx(
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

func NewsletterStore(
	ctx echo.Context,
	db database.Queries,
	sess *sessions.CookieStore,
	newsletterUsecase usecases.Newsletter,
) error {
	preview := ctx.QueryParam("preview")
	var storeNewsletterPayload newsletterPayload
	if err := ctx.Bind(&storeNewsletterPayload); err != nil {
		return err
	}

	if preview == "true" {
		return previewNewsletter(
			ctx,
			db,
			newsletterUsecase,
			storeNewsletterPayload,
			"post",
			"newsletters/store",
		)
	}

	newsletter, err := newsletterUsecase.Create(ctx.Request().Context(),
		storeNewsletterPayload.Title,
		storeNewsletterPayload.Edition,
		storeNewsletterPayload.ParagraphElements,
		storeNewsletterPayload.ArticleID,
	)
	if err != nil {
		return err
	}

	if storeNewsletterPayload.ReleaseOnCreate == "" {
		s, err := sess.Get(ctx.Request(), "flashMsg")
		if err != nil {
			return err
		}

		s.AddFlash("newsletter-draft-saved")
		if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
			return err
		}
	}

	if storeNewsletterPayload.ReleaseOnCreate == "on" {
		return releaseNewsletter(
			ctx,
			newsletter,
			db,
			newsletterUsecase,
			storeNewsletterPayload,
			sess,
			"post",
			"newsletters/store",
		)
	}

	return controllers.RedirectHx(
		ctx.Response().Writer,
		fmt.Sprintf("/dashboard/newsletters/%v/edit", newsletter.ID),
	)
}
