package handlers

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/router/routes"
	"github.com/mbvisti/mortenvistisen/services"
	"github.com/mbvisti/mortenvistisen/views/fragments"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/mbvisti/mortenvistisen/views"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type App struct {
	db           psql.Postgres
	cacheManager *CacheManager
}

func newApp(
	db psql.Postgres,
	cacheManager *CacheManager,
) App {
	return App{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (a App) LandingPage(c echo.Context) error {
	if value, ok := a.cacheManager.GetPageCache().Get(landingPageCacheKey); ok {
		return views.HomePage(value).Render(renderArgs(c))
	}

	articles, err := models.GetPublishedArticles(extractCtx(c), a.db.Pool)
	if err != nil {
		return err
	}

	var payload []views.HomeArticle
	for _, ar := range articles {
		tags := make([]string, len(ar.Tags))

		for i, t := range ar.Tags {
			tags[i] = t.Title
		}

		payload = append(payload, views.HomeArticle{
			Title:       ar.Title,
			Description: ar.Excerpt,
			Slug:        ar.Slug,
			PublishedAt: ar.FirstPublishedAt,
			Tags:        tags,
		})
	}

	homeComponent := views.Home(payload)
	if ok := a.cacheManager.GetPageCache().Set(landingPageCacheKey, homeComponent); !ok {
		slog.ErrorContext(
			c.Request().Context(),
			"could not set landing page cache",
			"error",
			err,
		)
	}

	return views.HomePage(homeComponent).Render(renderArgs(c))
}

func (a App) NewslettersPage(c echo.Context) error {
	newsletters, err := models.GetPublishedNewsletters(extractCtx(c), a.db.Pool)
	if err != nil {
		return err
	}

	// Group newsletters by year
	newslettersByYear := make(map[int][]views.HomeNewsletter)
	for _, newsletter := range newsletters {
		year := newsletter.ReleasedAt.Year()
		newslettersByYear[year] = append(
			newslettersByYear[year],
			views.HomeNewsletter{
				Title:      newsletter.Title,
				Slug:       newsletter.Slug,
				ReleasedAt: newsletter.ReleasedAt,
				Excerpt:    "", // Placeholder for future excerpt field
			},
		)
	}

	return views.NewslettersPage(views.Newsletters(newslettersByYear)).
		Render(renderArgs(c))
}

func (a App) AboutPage(c echo.Context) error {
	return views.AboutPage().Render(renderArgs(c))
}

func (a App) NotFoundPage(c echo.Context) error {
	return views.NotFoundPage().Render(renderArgs(c))
}

func (a App) ProjectsPage(c echo.Context) error {
	return views.ProjectsPage().Render(renderArgs(c))
}

func (a App) ArticlePage(c echo.Context) error {
	slug := c.Param("articleSlug")

	if value, ok := a.cacheManager.GetArticleCache().Get(articlePageCacheKey + slug); ok {
		return views.ArticlePage(value).
			Render(renderArgs(c))
	}

	article, err := models.GetArticleBySlug(
		c.Request().Context(),
		a.db.Pool,
		slug,
	)
	if err != nil {
		return err
	}

	manager := NewManager()
	ar, e := manager.ParseContent(article.Content)
	if e != nil {
		return e
	}

	props := views.ArticlePageProps{
		Slug:      slug,
		MetaTitle: article.MetaTitle,
		MetaDesc:  article.MetaDescription,
		Content: views.Article(
			article.ImageLink,
			ar,
			article.Title,
			article.FirstPublishedAt,
			article.UpdatedAt,
		),
	}

	if ok := a.cacheManager.GetArticleCache().Set(articlePageCacheKey+slug, props); !ok {
		slog.ErrorContext(
			c.Request().Context(),
			"could not set article cache",
			"error",
			err,
		)
	}

	return views.ArticlePage(props).
		Render(renderArgs(c))
}

func (a App) NewsletterPage(c echo.Context) error {
	slug := c.Param("newsletterSlug")

	if value, ok := a.cacheManager.GetNewsletterCache().Get(newsletterPageCacheKey + slug); ok {
		return views.NewsletterPage(value).
			Render(renderArgs(c))
	}

	newsletter, err := models.GetNewsletterBySlug(
		c.Request().Context(),
		a.db.Pool,
		slug,
	)
	if err != nil {
		return err
	}

	manager := NewManager()
	content, e := manager.ParseContent(newsletter.Content)
	if e != nil {
		return e
	}

	props := views.NewsletterPageProps{
		Slug:      slug,
		MetaTitle: newsletter.Title,
		MetaDesc:  newsletter.Title, // Using title as meta description since no specific description field
		Content: views.Newsletter(
			content,
			newsletter.Title,
			newsletter.ReleasedAt,
			newsletter.UpdatedAt,
		),
	}

	if ok := a.cacheManager.GetNewsletterCache().Set(newsletterPageCacheKey+slug, props); !ok {
		slog.ErrorContext(
			c.Request().Context(),
			"could not set newsletter cache",
			"error",
			err,
		)
	}

	return views.NewsletterPage(props).
		Render(renderArgs(c))
}

func (a App) Redirect(c echo.Context) error {
	to := c.QueryParam("to")
	for _, r := range routes.AllRoutes {
		if to == r.Path {
			return redirectHx(c.Response(), to)
		}
	}

	slog.InfoContext(c.Request().Context(),
		"security warning: someone tried to missue open redirect",
		"to", to,
		"ip", c.RealIP(),
	)

	return redirect(c.Response(), c.Request(), "/")
}

type Manager struct {
	posts           embed.FS
	markdownHandler goldmark.Markdown
}

func NewManager() Manager {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.TabWidth(4),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return Manager{
		markdownHandler: md,
	}
}

func (pm *Manager) GetPost(name string) (string, error) {
	source, err := pm.posts.ReadFile(name)
	if err != nil {
		slog.Error("failed to read markdown file", "error", err)
		return "", err
	}

	return string(source), nil
}

func (pm *Manager) Parse(name string) (string, error) {
	source, err := pm.posts.ReadFile(name)
	if err != nil {
		slog.Error("failed to read markdown file", "error", err)
		return "", err
	}

	// Parse Markdown content
	var htmlOutput bytes.Buffer
	if err := pm.markdownHandler.Convert(source, &htmlOutput); err != nil {
		slog.Error("failed to parse markdown file", "error", err)
		return "", err
	}

	return htmlOutput.String(), nil
}

func (pm *Manager) ParseContent(content string) (string, error) {
	// Parse Markdown content
	var htmlOutput bytes.Buffer
	if err := pm.markdownHandler.Convert([]byte(content), &htmlOutput); err != nil {
		slog.Error("failed to parse markdown file", "error", err)
		return "", err
	}

	return htmlOutput.String(), nil
}

func (a App) SubscribeNewsletter(c echo.Context) error {
	email := c.FormValue("email")
	referer := c.FormValue("referer")
	turnstileToken := c.FormValue("turnstileToken")

	span := trace.SpanFromContext(c.Request().Context())
	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("referer", referer),
		attribute.String("turnstileToken", turnstileToken),
	)

	// Validate Turnstile
	isValid, err := verifyTurnstileToken(
		c.Request().Context(),
		turnstileToken,
		c.RealIP(),
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"turnstile verification error",
			"error",
			err,
		)

		return fragments.NewsletterSubscription("unknown", true).
			Render(renderArgs(c))
	}

	if !isValid {
		slog.ErrorContext(
			c.Request().Context(),
			"turnstile invalid",
		)
		return fragments.NewsletterSubscription("unknown", true).
			Render(renderArgs(c))
	}

	_, _, err = services.SubscribeToNewsletter(
		c.Request().Context(),
		a.db,
		a.db.Queue(),
		email,
		referer,
	)
	if err != nil {
		if errors.Is(err, services.ErrSubscriberExists) {
			return fragments.NewsletterError("You're already subscribed to our newsletter!").
				Render(renderArgs(c))
		}

		slog.ErrorContext(
			c.Request().Context(),
			"error subscribing to newsletter",
			"error",
			err,
			"email",
			email,
		)
		return fragments.NewsletterError("").Render(renderArgs(c))
	}

	return fragments.NewsletterVerificationForm(email).Render(renderArgs(c))
}

func (a App) VerifyNewsletterSubscription(c echo.Context) error {
	email := c.FormValue("email")
	code := c.FormValue("code")

	span := trace.SpanFromContext(c.Request().Context())
	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("code", code),
	)

	if email == "" || code == "" {
		return fragments.NewsletterError("Email and verification code are required.").
			Render(renderArgs(c))
	}

	err := services.VerifySubscriberEmail(
		c.Request().Context(),
		a.db,
		email,
		code,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"error verifying newsletter subscription",
			"error",
			err,
			"email",
			email,
		)

		if err.Error() == "subscriber already verified" {
			return fragments.NewsletterError("This email is already verified!").
				Render(renderArgs(c))
		}

		return fragments.NewsletterError("Invalid verification code. Please check your email and try again.").
			Render(renderArgs(c))
	}

	return fragments.NewsletterSuccess().Render(renderArgs(c))
}

func (a App) VerifyNewsletterPage(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return fragments.NewsletterError("Invalid verification link. Please try again.").
			Render(renderArgs(c))
	}

	return fragments.NewsletterVerificationForm(email).Render(renderArgs(c))
}

type turnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Action      string   `json:"action"`
	CData       string   `json:"cdata"`
}

func verifyTurnstileToken(
	ctx context.Context,
	token string,
	remoteIP string,
) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("turnstile token is required")
	}

	data := url.Values{}
	data.Set("secret", config.Cfg.TurnstileSecretKey)
	data.Set("response", token)
	if remoteIP != "" {
		data.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to verify turnstile token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf(
			"turnstile verification failed with status: %d with body: %v",
			resp.StatusCode,
			string(b),
		)
	}

	var turnstileResp turnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&turnstileResp); err != nil {
		return false, fmt.Errorf("failed to decode turnstile response: %w", err)
	}

	return turnstileResp.Success, nil
}

func (a App) HandleUnsubscribe(c echo.Context) error {
	tokenValue := c.Param("token")
	if tokenValue == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Token is required",
		})
	}

	ctx := c.Request().Context()

	// Get and validate the token
	token, err := models.GetHashedToken(ctx, a.db.Pool, tokenValue)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Invalid or expired unsubscribe link",
		})
	}

	// Validate token is for unsubscribe scope
	if token.Meta.Scope != models.ScopeUnsubscribe {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid token scope",
		})
	}

	// Validate token is for subscriber resource
	if token.Meta.Resource != models.ResourceSubscriber {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid token resource",
		})
	}

	// Check if token is still valid
	if !token.IsValid() {
		return c.JSON(http.StatusGone, map[string]string{
			"error": "Unsubscribe link has expired",
		})
	}

	// Get the subscriber to ensure they exist
	subscriber, err := models.GetSubscriber(ctx, a.db.Pool, token.Meta.ResourceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Subscriber not found",
		})
	}

	// Start database transaction
	tx, err := a.db.BeginTx(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process unsubscribe",
		})
	}
	defer a.db.RollBackTx(ctx, tx)

	// Delete the subscriber
	if err := models.DeleteSubscriber(ctx, tx, subscriber.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to unsubscribe",
		})
	}

	// Delete the token to prevent reuse
	if err := models.DeleteToken(ctx, tx, token.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process unsubscribe",
		})
	}

	// Commit the transaction
	if err := a.db.CommitTx(ctx, tx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to complete unsubscribe",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Successfully unsubscribed from newsletter",
		"email":   subscriber.Email,
	})
}
