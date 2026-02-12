package router

import (
	"errors"
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/middleware"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

func (r Router) RegisterTagRoutes(tag controllers.Tags) error {
	errs := []error{}
	adminOnly := []echo.MiddlewareFunc{middleware.AdminOnly}

	_, err := r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.TagIndex.Path(),
		Name:        routes.TagIndex.Name(),
		Handler:     tag.Index,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.TagShow.Path(),
		Name:        routes.TagShow.Name(),
		Handler:     tag.Show,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.TagNew.Path(),
		Name:        routes.TagNew.Name(),
		Handler:     tag.New,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        routes.TagCreate.Path(),
		Name:        routes.TagCreate.Name(),
		Handler:     tag.Create,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.TagEdit.Path(),
		Name:        routes.TagEdit.Name(),
		Handler:     tag.Edit,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        routes.TagUpdate.Path(),
		Name:        routes.TagUpdate.Name(),
		Handler:     tag.Update,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        routes.TagDestroy.Path(),
		Name:        routes.TagDestroy.Name(),
		Handler:     tag.Destroy,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
