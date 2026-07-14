package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"mortenvistisen/internal/hypermedia"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/router"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

	"github.com/labstack/echo/v5"
)

type Pages struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
}

func NewPages(
	db storage.Pool,
	insertOnly queue.InsertOnly,
) Pages {
	return Pages{db, insertOnly}
}

func (p Pages) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.HomePage.Path(),
		Name:    routes.HomePage.Name(),
		Handler: p.Home,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method: http.MethodHead,
		Path:   routes.HomePage.Path(),
		Name:   routes.HomePage.Name() + ".head",
		Handler: func(etx *echo.Context) error {
			etx.NoContent(http.StatusOK)
			return nil
		},
	})
	if err != nil {
		errs = append(errs, err)
	}

	_ = r.AddRouteNotFound(p.NotFound)

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AboutPage.Path(),
		Name:    routes.AboutPage.Name(),
		Handler: p.About,
	})
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (p Pages) Home(etx *echo.Context) error {
	ctx := etx.Request().Context()
	db := p.db.Executor()

	articles, err := models.Article.LatestPublished(ctx, db, 3)
	if err != nil {
		return fmt.Errorf("load landing articles: %w", err)
	}

	projects, err := models.Project.LatestPublished(ctx, db, 2)
	if err != nil {
		return fmt.Errorf("load landing projects: %w", err)
	}

	newsletters, err := models.Newsletter.LatestPublished(ctx, db, 1)
	if err != nil {
		return fmt.Errorf("load landing newsletters: %w", err)
	}

	return hypermedia.RenderPage(etx, views.Landing{
		Articles:    articles,
		Projects:    projects,
		Newsletters: newsletters,
	}.Page())
}

func publicNotFound(etx *echo.Context) error {
	return hypermedia.RenderPage(etx, views.NotFound(), hypermedia.WithStatus(http.StatusNotFound))
}

func (p Pages) NotFound(etx *echo.Context) error {
	return publicNotFound(etx)
}

func (p Pages) About(etx *echo.Context) error {
	return hypermedia.RenderPage(etx, views.PageAbout())
}
