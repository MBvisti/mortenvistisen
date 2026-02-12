package router

import (
	"errors"
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/middleware"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

func (r Router) RegisterPagesRoutes(pages controllers.Pages) error {
	errs := []error{}

	_, err := r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.HomePage.Path(),
		Name:    routes.HomePage.Name(),
		Handler: pages.Home,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AboutPage.Path(),
		Name:    routes.AboutPage.Name(),
		Handler: pages.About,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.ArticleOverview.Path(),
		Name:    routes.ArticleOverview.Name(),
		Handler: pages.ArticlesOverview,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Article.Path(),
		Name:    routes.Article.Name(),
		Handler: pages.Article,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.NewsletterOverview.Path(),
		Name:    routes.NewsletterOverview.Name(),
		Handler: pages.NewslettersOverview,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Newsletter.Path(),
		Name:    routes.Newsletter.Name(),
		Handler: pages.Newsletter,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.ProjectOverview.Path(),
		Name:    routes.ProjectOverview.Name(),
		Handler: pages.ProjectsOverview,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Project.Path(),
		Name:    routes.Project.Name(),
		Handler: pages.Project,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminHome.Path(),
		Name:        routes.AdminHome.Name(),
		Handler: pages.AdminHome,
		Middlewares: []echo.MiddlewareFunc{
			middleware.AdminOnly,
		},
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
