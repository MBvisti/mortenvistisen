package router

import (
	"errors"
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/middleware"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

func (r Router) RegisterArticleRoutes(article controllers.Articles) error {
	errs := []error{}

	adminOnly := []echo.MiddlewareFunc{middleware.AdminOnly}

	_, err := r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.ArticleIndex.Path(),
		Name:        routes.ArticleIndex.Name(),
		Handler:     article.Index,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.ArticleShow.Path(),
		Name:        routes.ArticleShow.Name(),
		Handler:     article.Show,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.ArticleNew.Path(),
		Name:        routes.ArticleNew.Name(),
		Handler:     article.New,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        routes.ArticleCreate.Path(),
		Name:        routes.ArticleCreate.Name(),
		Handler:     article.Create,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.ArticleEdit.Path(),
		Name:        routes.ArticleEdit.Name(),
		Handler:     article.Edit,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        routes.ArticleUpdate.Path(),
		Name:        routes.ArticleUpdate.Name(),
		Handler:     article.Update,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        routes.ArticleDestroy.Path(),
		Name:        routes.ArticleDestroy.Name(),
		Handler:     article.Destroy,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
