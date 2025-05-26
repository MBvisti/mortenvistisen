package handlers

import (
	"log/slog"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/router/routes"
	"github.com/mbvisti/mortenvistisen/views"
)

const landingPageCacheKey = "landingPage"

type App struct {
	db    psql.Postgres
	cache otter.Cache[string, templ.Component]
}

func newApp(
	db psql.Postgres,
) App {
	cacheBuilder, err := otter.NewBuilder[string, templ.Component](20)
	if err != nil {
		panic(err)
	}

	pageCacher, err := cacheBuilder.WithTTL(48 * time.Hour).Build()
	if err != nil {
		panic(err)
	}

	return App{db, pageCacher}
}

func (a App) LandingPage(c echo.Context) error {
	if value, ok := a.cache.Get(landingPageCacheKey); ok {
		return views.HomePage(value).Render(renderArgs(c))
	}

	return views.HomePage(views.Home()).Render(renderArgs(c))
}

func (a App) AboutPage(c echo.Context) error {
	return views.AboutPage().Render(renderArgs(c))
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
