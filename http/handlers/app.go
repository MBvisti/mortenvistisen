package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/fragments"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
)

const landingPageCacheKey = "LandingPage"

type App struct {
	db          psql.Postgres
	cache       otter.CacheWithVariableTTL[string, string]
	email       services.Mail
	postManager posts.Manager
}

func newApp(
	db psql.Postgres,
	cache otter.CacheWithVariableTTL[string, string],
	email services.Mail,
	postManager posts.Manager,
) App {
	return App{db, cache, email, postManager}
}

func (a *App) LandingPage(ctx echo.Context) error {
	articlesPage, err := models.GetArticlesPage(
		ctx.Request().Context(),
		1,
		5,
		a.db.Pool,
	)
	if err != nil {
		return views.ErrorPage().Render(renderArgs(ctx))
	}

	slog.Info("articles", "num", len(articlesPage.Articles))

	posts := make([]views.Post, len(articlesPage.Articles))
	for i, article := range articlesPage.Articles {
		tags := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tags[i] = tag.Name
		}

		posts[i] = views.Post{
			Title:       article.Title,
			Slug:        article.Slug,
			Excerpt:     article.Excerpt,
			ReleaseDate: article.ReleaseDate.String(),
			Tags:        tags,
		}
	}

	return views.HomePage(posts, csrf.Token(ctx.Request())).
		Render(renderArgs(ctx))
}

func (a *App) AboutPage(ctx echo.Context) error {
	return views.AboutPage().Render(renderArgs(ctx))
}

func (a *App) ArticlePage(ctx echo.Context) error {
	slug := ctx.Param("postSlug")
	article, err := models.GetArticleBySlug(
		ctx.Request().Context(),
		slug,
		a.db.Pool,
	)
	if err != nil {
		return views.ErrorPage().Render(renderArgs(ctx))
	}

	postContent, err := a.postManager.Parse(article.Filename)
	if err != nil {
		return err
	}

	return views.ArticlePage(views.ArticlePageData{
		HeaderTitle: article.Title,
		Content:     postContent,
		ReleaseDate: article.ReleaseDate,
	}).Render(renderArgs(ctx))
}

func (a *App) ArticlesPage(ctx echo.Context) error {
	articles, err := models.GetAllArticles(
		ctx.Request().Context(),
		a.db.Pool,
	)
	if err != nil {
		return views.ErrorPage().Render(renderArgs(ctx))
	}

	posts := make([]views.Post, len(articles))
	for i, article := range articles {
		tags := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tags[i] = tag.Name
		}
		posts[i] = views.Post{
			Title:       article.Title,
			Slug:        article.Slug,
			Excerpt:     article.Excerpt,
			ReleaseDate: article.ReleaseDate.String(),
			Tags:        tags,
		}
	}

	return views.ArticlesOverview(posts).Render(renderArgs(ctx))
}

func (a *App) ProjectsPage(ctx echo.Context) error {
	return views.ProjectsPage().Render(renderArgs(ctx))
}

func (a *App) NewslettersPage(ctx echo.Context) error {
	return views.NewslettersPage().Render(renderArgs(ctx))
}

func (a *App) SubscriptionEvent(c echo.Context) error {
	type subscriptionEventForm struct {
		Email             string `form:"hero-input"`
		Title             string `form:"article-title"`
		TurnstileResponse string `form:"cf-turnstile-response"`
	}

	var form subscriptionEventForm
	if err := c.Bind(&form); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to bind subscription form",
			"error",
			err,
		)
		c.Response().Header().Add("HX-Retarget", "body")
		c.Response().Header().Add("HX-Reswap", "outerHTML")
		return views.ErrorPage().Render(renderArgs(c))
	}

	ip := c.RealIP()

	formData := url.Values{}
	formData.Set("secret", config.Cfg.TurnstileSecretKey)
	formData.Set("response", form.TurnstileResponse)
	formData.Set("remoteip", ip)

	cfVerifyReq, err := http.NewRequestWithContext(c.Request().Context(),
		http.MethodPost,
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to create cloudflare verification request",
			"error",
			err,
		)

		c.Response().Header().Add("HX-Retarget", "body")
		c.Response().Header().Add("HX-Reswap", "outerHTML")
		return views.ErrorPage().Render(renderArgs(c))
	}

	cfVerifyReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(cfVerifyReq)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to execute cloudflare verification request",
			"error",
			err,
		)
		return fragments.SubscribeResponse(false).Render(renderArgs(c))
	}

	defer res.Body.Close()

	type verificationResponse struct {
		Success bool `json:"success"`
	}

	var firstOutcome verificationResponse
	if err := json.NewDecoder(res.Body).Decode(&firstOutcome); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to decode cloudflare verification response",
			"error",
			err,
		)
		return fragments.SubscribeResponse(false).Render(renderArgs(c))
	}

	if firstOutcome.Success {
		return fragments.SubscribeResponse(false).Render(renderArgs(c))
	}

	subcriber, err := models.NewSubscriber(
		c.Request().Context(),
		a.db.Pool,
		models.NewSubscriberPayload{
			Email:        form.Email,
			SubscribedAt: time.Now(),
			Referer:      form.Title,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fragments.SubscribeResponse(true).Render(renderArgs(c))
		}

		slog.ErrorContext(
			c.Request().Context(),
			"failed to create new subscriber",
			"error",
			err,
			"email",
			form.Email,
		)
		return fragments.SubscribeResponse(false).Render(renderArgs(c))
	}

	now := time.Now()

	activationTkn, err := models.NewToken(
		c.Request().Context(),
		models.NewTokenPayload{
			Expiration: now.Add(48 * time.Hour),
			Meta: models.MetaInformation{
				Resource:   models.ResourceSubscriber,
				ResourceID: subcriber.ID,
				Scope:      models.ScopeEmailVerification,
			},
		},
		a.db.Pool,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to create activation token",
			"error",
			err,
			"subscriber_id",
			subcriber.ID,
		)
		return err
	}

	unsubTkn, err := models.NewToken(
		c.Request().Context(),
		models.NewTokenPayload{
			Expiration: now.Add(365 * (24 * time.Hour)),
			Meta: models.MetaInformation{
				Resource:   models.ResourceSubscriber,
				ResourceID: subcriber.ID,
				Scope:      models.ScopeUnsubscribe,
			},
		},
		a.db.Pool,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to create unsubscribe token",
			"error",
			err,
			"subscriber_id",
			subcriber.ID,
		)
		return err
	}

	if err := a.email.SendNewSubscriber(
		c.Request().Context(), subcriber.Email, activationTkn, unsubTkn,
	); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to send new subscriber email",
			"error",
			err,
			"subscriber_id",
			subcriber.ID,
			"email",
			subcriber.Email,
		)
		return err
	}

	return fragments.SubscribeResponse(false).Render(renderArgs(c))
}
