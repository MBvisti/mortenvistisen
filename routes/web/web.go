package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/labstack/echo/v4"
)

type Web struct {
	controllers controllers.Controller
	router      *echo.Echo
}

func NewWeb(router *echo.Echo, controllers controllers.Controller) Web {
	return Web{
		controllers,
		router,
	}
}

func (w *Web) miscRoutes() {
	w.router.GET("/robots.txt", func(c echo.Context) error {
		// Send the sitemap file
		return c.File("./resources/seo/robots.txt")
	})
	w.router.GET("/sitemap.xml", func(c echo.Context) error {
		// Send the sitemap file
		return c.File("./resources/seo/sitemap.xml")
	})
	w.router.GET("/static/css/output.css", func(c echo.Context) error {

		// Set cache headers for one year (adjust as needed)
		cacheTime := time.Now().AddDate(0, 0, 1)

		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		c.Response().Header().Set(echo.HeaderLastModified, cacheTime.UTC().Format(http.TimeFormat))

		return c.File("./static/css/output.css")
	})
	w.router.GET("/static/js/:filename", func(c echo.Context) error {
		fm := c.Param("filename")

		// Set cache headers for one year (adjust as needed)
		cacheTime := time.Now().AddDate(0, 1, 0)

		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		c.Response().Header().Set(echo.HeaderLastModified, cacheTime.UTC().Format(http.TimeFormat))

		return c.File(fmt.Sprintf("./static/js/%s", fm))
	})
	w.router.GET("/static/images/:filename", func(c echo.Context) error {
		fm := c.Param("filename")

		// Set cache headers for one year (adjust as needed)
		cacheTime := time.Now().AddDate(0, 1, 0)

		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		c.Response().Header().Set(echo.HeaderLastModified, cacheTime.UTC().Format(http.TimeFormat))

		return c.File(fmt.Sprintf("./static/images/%s", fm))
	})
}

func (w *Web) SetupWebRoutes() {
	w.UtilityRoutes()
	w.miscRoutes()

	w.ArticleRoutes()
	w.SiteRoutes()
	w.SubscribeRoutes()
	// w.UserRoutes()
	w.DashboardRoutes()
}
