package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/fragments"
	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
)

const (
	landingPageCacheKey     = "landingPage"
	articlesPageCacheKey    = "articlesPage"
	newslettersPageCacheKey = "newslettersPage"
)

var cacheDuration = time.Hour * time.Duration(168)

type App struct {
	db          psql.Postgres
	cache       otter.CacheWithVariableTTL[string, templ.Component]
	email       services.Mail
	postManager posts.Manager
}

func newApp(
	db psql.Postgres,
	cache otter.CacheWithVariableTTL[string, templ.Component],
	email services.Mail,
	postManager posts.Manager,
) App {
	return App{db, cache, email, postManager}
}

func (a *App) LandingPage(c echo.Context) error {
	if value, ok := a.cache.Get(landingPageCacheKey); ok {
		return views.HomePage(value).Render(renderArgs(c))
	}

	articlesPage, err := models.GetArticlesPage(
		c.Request().Context(),
		1,
		5,
		a.db.Pool,
	)
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}

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
			ReleaseDate: article.ReleaseDate,
			Tags:        tags,
		}
	}

	cachedComponent := views.Home(posts)
	if ok := a.cache.Set(landingPageCacheKey, cachedComponent, cacheDuration); !ok {
		return views.HomePage(cachedComponent).Render(renderArgs(c))
	}

	return views.HomePage(cachedComponent).
		Render(renderArgs(c))
}

func (a *App) AboutPage(c echo.Context) error {
	return views.AboutPage().Render(renderArgs(c))
}

func (a *App) ArticlePage(c echo.Context) error {
	slug := c.Param("postSlug")
	article, err := models.GetArticleBySlug(
		c.Request().Context(),
		slug,
		a.db.Pool,
	)
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}

	postContent, err := a.postManager.Parse(article.Filename)
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}

	return views.ArticlePage(views.ArticlePageData{
		Title:       article.Title,
		HeaderTitle: article.HeaderTitle,
		Content:     postContent,
		ReleaseDate: article.ReleaseDate,
	}).Render(renderArgs(c))
}

func (a *App) ArticlesPage(c echo.Context) error {
	if value, ok := a.cache.Get(articlesPageCacheKey); ok {
		return views.ArticlesPage(value).Render(renderArgs(c))
	}

	articles, err := models.GetAllArticles(
		c.Request().Context(),
		a.db.Pool,
	)
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}

	postsByYear := map[int][]views.Post{}
	var years []int

	for _, article := range articles {
		year := article.ReleaseDate.Year()

		if _, exists := postsByYear[year]; !exists {
			years = append(years, year)
		}

		tags := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tags[i] = tag.Name
		}

		postsByYear[year] = append(
			postsByYear[year],
			views.Post{
				Title:       article.Title,
				Slug:        article.Slug,
				Excerpt:     article.Excerpt,
				ReleaseDate: article.ReleaseDate,
				Tags:        tags,
			},
		)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	orderedPosts := make([]views.YearlyPosts, len(years))
	for i, year := range years {
		orderedPosts[i] = views.YearlyPosts{
			Year:     strconv.Itoa(year),
			Articles: postsByYear[year],
		}
	}

	cachedComponent := views.Articles(orderedPosts)
	if ok := a.cache.Set(articlesPageCacheKey, cachedComponent, cacheDuration); !ok {
		return views.ArticlesPage(cachedComponent).Render(renderArgs(c))
	}

	return views.ArticlesPage(cachedComponent).Render(renderArgs(c))
}

func (a *App) ProjectsPage(c echo.Context) error {
	return views.ProjectsPage().Render(renderArgs(c))
}

func (a *App) NewsletterPage(c echo.Context) error {
	newsletter, err := models.GetNewsletterBySlug(
		c.Request().Context(),
		a.db.Pool,
		c.Param("slug"),
	)
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}

	return views.NewsletterPage(newsletter).Render(renderArgs(c))
}

func (a *App) NewslettersPage(c echo.Context) error {
	if value, ok := a.cache.Get(newslettersPageCacheKey); ok {
		return views.NewslettersPage(value).Render(renderArgs(c))
	}

	newsletters, err := models.GetAllNewsletters(
		c.Request().Context(),
		a.db.Pool,
	)
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}

	newslettersByYear := map[int][]models.Newsletter{}
	var years []int

	for _, newsletter := range newsletters {
		year := newsletter.ReleasedAt.Year()

		if _, exists := newslettersByYear[year]; !exists {
			years = append(years, year)
		}

		newslettersByYear[year] = append(
			newslettersByYear[year],
			models.Newsletter{
				ID:         newsletter.ID,
				CreatedAt:  newsletter.CreatedAt,
				UpdatedAt:  newsletter.UpdatedAt,
				Title:      newsletter.Title,
				Content:    newsletter.Content,
				ReleasedAt: newsletter.ReleasedAt,
				Released:   newsletter.Released,
			},
		)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	orderedNewsletters := make([]views.YearlyNewsletters, len(years))
	for i, year := range years {
		orderedNewsletters[i] = views.YearlyNewsletters{
			Year:        strconv.Itoa(year),
			Newsletters: newslettersByYear[year],
		}
	}

	cachedComponent := views.Newsletters(orderedNewsletters)
	if ok := a.cache.Set(newslettersPageCacheKey, cachedComponent, cacheDuration); !ok {
		return views.NewslettersPage(cachedComponent).Render(renderArgs(c))
	}

	return views.NewslettersPage(cachedComponent).Render(renderArgs(c))
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
		return errorPage(c, views.ErrorPage())
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

		return errorPage(c, views.ErrorPage())
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

	if !firstOutcome.Success {
		return fragments.NewsletterForm(true).
			Render(renderArgs(c))
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
		return errorPage(c, views.ErrorPage())
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
		return errorPage(c, views.ErrorPage())
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
		return errorPage(c, views.ErrorPage())
	}

	return fragments.SubscribeResponse(false).Render(renderArgs(c))
}

func (a App) SubscriberEmailVerification(c echo.Context) error {
	type verificationRequest struct {
		Token string `query:"token"`
	}

	var payload verificationRequest
	if err := c.Bind(&payload); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	tx, err := a.db.BeginTx(c.Request().Context())
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}
	defer tx.Rollback(c.Request().Context())

	tkn, err := models.GetToken(c.Request().Context(), payload.Token, tx)
	if err != nil {
		return errorPage(
			c,
			views.ErrorPage(
				views.WithErrPageTitle("That doesn't seem right"),
				views.WithErrPageMsg(
					"Your token is either not valid anymore or have been deleted. Please request a new one or contact support.",
				),
			),
		)
	}

	if !tkn.IsValid() || tkn.Meta.Scope != models.ScopeEmailVerification {
		return views.SubscriberVerification(false).Render(renderArgs(c))
	}

	if _, err := models.UpdateSubscriberVerification(c.Request().Context(), tx, tkn.Meta.ResourceID, true); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	if err := models.DeleteToken(c.Request().Context(), tkn.ID, tx); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	if err := tx.Commit(c.Request().Context()); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	return views.SubscriberVerification(true).Render(renderArgs(c))
}

func (a App) UnsubscriptionEvent(c echo.Context) error {
	type unsubscribeRequest struct {
		Token string `query:"token"`
	}

	var payload unsubscribeRequest
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"UnsubscriptionEvent",
			"error",
			err,
			"token",
			payload.Token,
		)
		return errorPage(c, views.ErrorPage())
	}

	tx, err := a.db.BeginTx(c.Request().Context())
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}
	defer tx.Rollback(c.Request().Context())

	tkn, err := models.GetToken(c.Request().Context(), payload.Token, tx)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"UnsubscriptionEvent",
			"error",
			err,
			"token",
			payload.Token,
		)
		return errorPage(
			c,
			views.ErrorPage(
				views.WithErrPageTitle("That doesn't seem right"),
				views.WithErrPageMsg(
					"Your token is either not valid anymore or have been deleted. Please request a new one or contact support.",
				),
			),
		)
	}
	if !tkn.IsValid() || tkn.Meta.Scope != models.ScopeUnsubscribe {
		return views.SubscriptionDeletePage(false).Render(renderArgs(c))
	}

	if err := models.DeleteSubscriber(c.Request().Context(), tx, tkn.Meta.ResourceID); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	if err := models.DeleteToken(c.Request().Context(), tkn.ID, tx); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	if err := tx.Commit(c.Request().Context()); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	return views.SubscriptionDeletePage(true).Render(renderArgs(c))
}
