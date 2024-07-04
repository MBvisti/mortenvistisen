package router

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/controllers/api"
	"github.com/MBvisti/mortenvistisen/controllers/app"
	"github.com/MBvisti/mortenvistisen/controllers/authentication"
	"github.com/MBvisti/mortenvistisen/controllers/dashboard"
	"github.com/MBvisti/mortenvistisen/http/middleware"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo-contrib/echoprometheus"
	echomw "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	router     *echo.Echo
	middleware middleware.Middleware
	cfg        config.Cfg
	ctrlDeps   controllers.Dependencies
}

func NewRouter(
	ctrlDeps controllers.Dependencies,
	mw middleware.Middleware,
	cfg config.Cfg,
	logger *slog.Logger,
) *Router {
	router := echo.New()

	if cfg.App.Environment == "development" {
		router.Debug = true
	}

	router.Static("/static", "static")
	router.Use(mw.RegisterUserContext)
	router.Use(echomw.Recover())
	router.Use(echomw.RequestLoggerWithConfig(echomw.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v echomw.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	router.Use(
		echoprometheus.NewMiddleware("mortenvistisen_blog"),
	) // adds middleware to gather metrics
	router.GET("/metrics", echoprometheus.NewHandler())

	router.GET("/robots.txt", func(c echo.Context) error {
		return c.File("./resources/seo/robots.txt")
	})
	router.GET("/sitemap.xml", func(c echo.Context) error {
		return c.File("./resources/seo/sitemap.xml")
	})
	router.GET("/static/css/output.css", func(c echo.Context) error {
		if os.Getenv("ENVIRONMENT") == "production" {
			// Set cache headers for one year (adjust as needed)
			cacheTime := time.Now().AddDate(0, 0, 1)

			c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
			c.Response().
				Header().
				Set(echo.HeaderLastModified, cacheTime.UTC().Format(http.TimeFormat))
		}

		return c.File("./static/css/output.css")
	})
	router.GET("/static/js/:filename", func(c echo.Context) error {
		fm := c.Param("filename")

		if os.Getenv("ENVIRONMENT") == "production" {
			// Set cache headers for one year (adjust as needed)
			cacheTime := time.Now().AddDate(0, 1, 0)

			c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
			c.Response().
				Header().
				Set(echo.HeaderLastModified, cacheTime.UTC().Format(http.TimeFormat))
		}

		return c.File(fmt.Sprintf("./static/js/%s", fm))
	})
	router.GET("/static/images/:filename", func(c echo.Context) error {
		fm := c.Param("filename")

		if os.Getenv("ENVIRONMENT") == "production" {
			// Set cache headers for one year (adjust as needed)
			cacheTime := time.Now().AddDate(0, 1, 0)

			c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
			c.Response().
				Header().
				Set(echo.HeaderLastModified, cacheTime.UTC().Format(http.TimeFormat))
		}

		return c.File(fmt.Sprintf("./static/images/%s", fm))
	})

	return &Router{
		router:     router,
		middleware: mw,
		cfg:        cfg,
		ctrlDeps:   ctrlDeps,
	}
}

func (r *Router) GetInstance() *echo.Echo {
	return r.router
}

func (r *Router) LoadInRoutes() {
	r.loadApiV1Routes()
	r.loadDashboardRoutes()
	r.loadAppRoutes()
	r.loadAuthRoutes()
}

func (r *Router) loadDashboardRoutes() {
	router := r.router.Group("/dashboard", r.middleware.AuthOnly)

	router.GET("", func(c echo.Context) error {
		return dashboard.Index(c)
	})

	router.GET("/subscribers", func(c echo.Context) error {
		return dashboard.SubscribersIndex(c, r.ctrlDeps.DB, r.ctrlDeps.SubscriberModel)
	})

	router.POST("/subscribers/:id/send-verification-mail", func(c echo.Context) error {
		return dashboard.ResendVerificationMail(
			c,
			r.ctrlDeps.SubscriberModel,
			r.ctrlDeps.TokenService,
			r.ctrlDeps.EmailService,
		)
	})

	router.GET("/articles", func(c echo.Context) error {
		return dashboard.ArticlesIndex(c, r.ctrlDeps.ArticleModel)
	})
	router.GET("/articles/:slug/edit", func(c echo.Context) error {
		return dashboard.ArticleEdit(c, r.ctrlDeps.ArticleModel, r.ctrlDeps.PostManager)
	})
	router.GET("/articles/create", func(c echo.Context) error {
		return dashboard.ArticleCreate(c, r.ctrlDeps.ArticleModel, r.ctrlDeps.TagModel)
	})
	router.POST("/articles/store", func(c echo.Context) error {
		return dashboard.ArticleStore(c, r.ctrlDeps.ArticleModel)
	})
	router.PUT("/articles/:id/update", func(c echo.Context) error {
		return dashboard.ArticleUpdate(c, r.ctrlDeps.ArticleModel)
	})
	router.POST("/tags/store", func(c echo.Context) error {
		return dashboard.TagStore(c, r.ctrlDeps.TagModel)
	})

	router.GET("/newsletters", func(c echo.Context) error {
		return dashboard.NewslettersIndex(c, r.ctrlDeps.NewsletterModel, r.ctrlDeps.CookieStore)
	})
	router.GET("/newsletters/create", func(c echo.Context) error {
		return dashboard.NewsletterCreate(c, r.ctrlDeps.NewsletterModel, r.ctrlDeps.ArticleModel)
	})
	router.POST("/newsletters/store", func(c echo.Context) error {
		return dashboard.NewsletterStore(
			c,
			r.ctrlDeps.CookieStore,
			r.ctrlDeps.NewsletterModel,
			r.ctrlDeps.ArticleModel,
		)
	})
	router.GET("/newsletters/:id/edit", func(c echo.Context) error {
		return dashboard.NewslettersEdit(
			c,
			r.ctrlDeps.NewsletterModel,
			r.ctrlDeps.ArticleModel,
			r.ctrlDeps.CookieStore,
		)
	})
	router.PUT("/newsletters/:id/update", func(c echo.Context) error {
		return dashboard.NewsletterUpdate(
			c,
			r.ctrlDeps.NewsletterModel,
			r.ctrlDeps.ArticleModel,
			r.ctrlDeps.CookieStore,
		)
	})

	router.DELETE("/subscribers/:ID", func(c echo.Context) error {
		return dashboard.DeleteSubscriber(c, r.ctrlDeps.SubscriberModel)
	})
}

func (r *Router) loadAppRoutes() {
	router := r.router.Group("")

	router.GET("/", func(c echo.Context) error {
		return app.Index(c, r.ctrlDeps.ArticleModel)
	})
	router.GET("", func(c echo.Context) error {
		return app.Index(c, r.ctrlDeps.ArticleModel)
	})

	router.GET("/about", func(c echo.Context) error {
		return app.About(c)
	})

	router.GET("/newsletter", func(c echo.Context) error {
		return app.Newsletter(c)
	})

	router.GET("/projects", func(c echo.Context) error {
		return app.Projects(c)
	})

	router.GET("/posts/:postSlug", func(c echo.Context) error {
		return app.Article(c, r.ctrlDeps.ArticleModel, r.ctrlDeps.PostManager)
	})

	router.GET("/modal", func(c echo.Context) error {
		return app.RenderModal(c)
	})

	router.GET("/verify-subscriber", func(c echo.Context) error {
		return app.SubscriberEmailVerification(
			c,
			r.ctrlDeps.SubscriberModel,
			r.ctrlDeps.TokenService,
		)
	})

	router.POST("/subscribe", func(c echo.Context) error {
		return app.SubscriptionEvent(
			c,
			r.ctrlDeps.SubscriberModel,
		)
	})

	router.GET("/redirect", func(c echo.Context) error {
		to := c.QueryParam("to")
		return controllers.RedirectHx(c.Response(), fmt.Sprintf("/%s", to))
	})
}

func (r *Router) loadAuthRoutes() {
	router := r.router.Group("")
	router.GET("/register", func(c echo.Context) error {
		return authentication.CreateUser(c)
	})
	router.POST("/register", func(c echo.Context) error {
		return authentication.StoreUser(
			c,
			r.ctrlDeps.UserModel,
			r.ctrlDeps.TokenService,
			r.ctrlDeps.EmailService,
		)
	})

	router.GET("/login", func(c echo.Context) error {
		return authentication.CreateAuthenticatedSession(c)
	})
	router.POST("/login", func(c echo.Context) error {
		return authentication.StoreAuthenticatedSession(
			c,
			r.ctrlDeps.AuthService,
		)
	})

	router.GET("/verify-email", func(c echo.Context) error {
		return authentication.UserEmailVerification(
			c,
			r.ctrlDeps.UserModel,
			r.ctrlDeps.TokenService,
			r.ctrlDeps.AuthService,
		)
	})

	router.GET("/forgot-password", func(c echo.Context) error {
		return authentication.CreatePasswordReset(c)
	})
	router.POST("/forgot-password", func(c echo.Context) error {
		return authentication.StorePasswordReset(
			c,
			r.ctrlDeps.UserModel,
			r.ctrlDeps.EmailService,
			r.ctrlDeps.TokenService,
			r.cfg,
		)
	})
	router.GET("/reset-password", func(c echo.Context) error {
		return authentication.CreateResetPassword(c)
	})
	router.POST("/reset-password", func(c echo.Context) error {
		return authentication.StoreResetPassword(
			c,
			r.ctrlDeps.UserModel,
			r.ctrlDeps.AuthService,
			r.ctrlDeps.TokenService,
		)
	})
}

func (r *Router) loadApiV1Routes() {
	router := r.router.Group("/api/v1")

	router.GET("/health", func(c echo.Context) error {
		return api.AppHealth(c)
	})
}
