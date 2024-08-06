package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MBvisti/mortenvistisen/http/handlers"
	"github.com/MBvisti/mortenvistisen/http/middleware"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"

	"github.com/labstack/echo-contrib/echoprometheus"
	echomw "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	router                 *echo.Echo
	middleware             middleware.Middleware
	cfg                    config.Cfg
	apiHandlers            handlers.Api
	appHandlers            handlers.App
	authenticationHandlers handlers.Authentication
	dashboardHandlers      handlers.Dashboard
	registrationHandlers   handlers.Registration
	baseHandlers           handlers.Base
}

func NewRouter(
	mw middleware.Middleware,
	cfg config.Cfg,
	apiHandlers handlers.Api,
	appHandlers handlers.App,
	authenticationHandlers handlers.Authentication,
	dashboardHandlers handlers.Dashboard,
	registrationHandlers handlers.Registration,
	baseHandlers handlers.Base,
) *Router {
	router := echo.New()

	router.Debug = true
	if cfg.App.Environment == config.PROD_ENVIRONMENT {
		router.Debug = false
		router.Use(echomw.GzipWithConfig(echomw.GzipConfig{
			Level: 5,
			Skipper: func(c echo.Context) bool {
				return strings.Contains(c.Path(), "metrics")
			},
		}))
		router.Use(
			echoprometheus.NewMiddleware("mortenvistisen_blog"),
		)
		router.GET("/metrics", echoprometheus.NewHandler())
	}

	router.Static("/static", "static")
	router.Use(mw.RegisterUserContext)
	router.Use(slogecho.New(slog.Default()))
	router.Use(echomw.Recover())

	router.GET("/robots.txt", func(c echo.Context) error {
		return c.File("./resources/seo/robots.txt")
	})
	router.GET("/4zd8j69sf3ju2hnfxmebr3czub8uu63m.txt", func(c echo.Context) error {
		return c.File("./resources/seo/index_now.txt")
	})
	router.GET("/sitemap.xml", func(c echo.Context) error {
		return c.File("./resources/seo/sitemap.xml")
	})
	router.GET("/static/css/output.css", func(c echo.Context) error {
		if os.Getenv("ENVIRONMENT") == config.PROD_ENVIRONMENT {
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
		router,
		mw,
		cfg,
		apiHandlers,
		appHandlers,
		authenticationHandlers,
		dashboardHandlers,
		registrationHandlers,
		baseHandlers,
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
