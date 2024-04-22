package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/server/middleware"
	"github.com/labstack/echo/v4"
)

type Routes struct {
	router      *echo.Echo
	controllers controllers.Controller
	middleware  middleware.Middleware
	cfg         config.Cfg
}

func NewRoutes(ctrl controllers.Controller, mw middleware.Middleware, cfg config.Cfg) *Routes {
	router := echo.New()

	if cfg.App.Environment == "development" {
		router.Debug = true
	}

	router.Static("/static", "static")
	router.Use(mw.RegisterUserContext)

	return &Routes{
		router:      router,
		controllers: ctrl,
		middleware:  mw,
		cfg:         cfg,
	}
}

func (r *Routes) web() {
	authRoutes(r.router, r.controllers)
	errorRoutes(r.router, r.controllers)
	dashboardRoutes(r.router, r.controllers, r.middleware)
	appRoutes(r.router, r.controllers)
	miscRoutes(r.router)
}

func (r *Routes) api() {
	apiRouter := r.router.Group("/api")
	apiRoutes(apiRouter, r.controllers)
}

func (r *Routes) SetupRoutes() *echo.Echo {
	r.web()
	r.api()

	return r.router
}

func appRoutes(router *echo.Echo, ctrl controllers.Controller) {
	router.GET("/", func(c echo.Context) error {
		return ctrl.HomeIndex(c)
	})

	router.GET("/about", func(c echo.Context) error {
		return ctrl.About(c)
	})

	router.GET("/newsletter", func(c echo.Context) error {
		return ctrl.Newsletter(c)
	})

	router.GET("/projects", func(c echo.Context) error {
		return ctrl.Projects(c)
	})

	router.GET("/posts/:postSlug", func(c echo.Context) error {
		return ctrl.Article(c)
	})

	router.GET("/modal", func(c echo.Context) error {
		return ctrl.RenderModal(c)
	})

	router.POST("/subscribe", func(c echo.Context) error {
		return ctrl.SubscriptionEvent(c)
	})
}

func miscRoutes(router *echo.Echo) {
	env := os.Getenv("ENVIRONMENT")

	router.GET("/robots.txt", func(c echo.Context) error {
		return c.File("./resources/seo/robots.txt")
	})

	router.GET("/sitemap.xml", func(c echo.Context) error {
		return c.File("./resources/seo/sitemap.xml")
	})

	router.GET("/static/css/output.css", func(c echo.Context) error {
		if env == "production" {
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

		if env == "production" {
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

		if env == "production" {
			// Set cache headers for one year (adjust as needed)
			cacheTime := time.Now().AddDate(0, 1, 0)

			c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
			c.Response().
				Header().
				Set(echo.HeaderLastModified, cacheTime.UTC().Format(http.TimeFormat))
		}

		return c.File(fmt.Sprintf("./static/images/%s", fm))
	})
}

func dashboardRoutes(router *echo.Echo, ctrl controllers.Controller, mw middleware.Middleware) {
	dashboardRouter := router.Group("/dashboard", mw.AuthOnly, mw.AdminOnly)

	dashboardRouter.GET("", func(c echo.Context) error {
		return ctrl.DashboardIndex(c)
	})

	dashboardRouter.GET("/subscribers", func(c echo.Context) error {
		return ctrl.DashboardSubscribers(c)
	})
	dashboardRouter.POST("/subscriber/:id/send-verification-mail", func(c echo.Context) error {
		return ctrl.DashboardResendVerificationMail(c)
	})

	dashboardRouter.GET("/articles", func(c echo.Context) error {
		return ctrl.DashboardArticles(c)
	})
	dashboardRouter.GET("/article/:slug/edit", func(c echo.Context) error {
		return ctrl.DashboardArticleEdit(c)
	})
	dashboardRouter.GET("/article/create", func(c echo.Context) error {
		return ctrl.DashboardArticleCreate(c)
	})
	dashboardRouter.POST("/article/store", func(c echo.Context) error {
		return ctrl.DashboadPostStore(c)
	})
	dashboardRouter.POST("/tag/store", func(c echo.Context) error {
		return ctrl.DashboadTagStore(c)
	})

	dashboardRouter.DELETE("/subscriber/:ID", func(c echo.Context) error {
		return ctrl.DeleteSubscriber(c)
	})
}

func errorRoutes(router *echo.Echo, ctrl controllers.Controller) {
	router.GET("/400", func(c echo.Context) error {
		return ctrl.InternalError(c)
	})

	router.GET("/404", func(c echo.Context) error {
		return ctrl.InternalError(c)
	})

	router.GET("/500", func(c echo.Context) error {
		return ctrl.InternalError(c)
	})

	router.GET("/redirect", func(c echo.Context) error {
		return ctrl.Redirect(c)
	})
}
