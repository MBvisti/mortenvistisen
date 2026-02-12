package router

import (
	"errors"
	"net/http"

	"mortenvistisen/controllers"
	"mortenvistisen/router/middleware"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

func (r Router) RegisterNewsletterRoutes(newsletter controllers.Newsletters) error {
	errs := []error{}
	adminOnly := []echo.MiddlewareFunc{middleware.AdminOnly}

	_, err := r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.NewsletterIndex.Path(),
		Name:        routes.NewsletterIndex.Name(),
		Handler:     newsletter.Index,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.NewsletterShow.Path(),
		Name:        routes.NewsletterShow.Name(),
		Handler:     newsletter.Show,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.NewsletterNew.Path(),
		Name:        routes.NewsletterNew.Name(),
		Handler:     newsletter.New,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        routes.NewsletterCreate.Path(),
		Name:        routes.NewsletterCreate.Name(),
		Handler:     newsletter.Create,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.NewsletterEdit.Path(),
		Name:        routes.NewsletterEdit.Name(),
		Handler:     newsletter.Edit,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        routes.NewsletterUpdate.Path(),
		Name:        routes.NewsletterUpdate.Name(),
		Handler:     newsletter.Update,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        routes.NewsletterDestroy.Path(),
		Name:        routes.NewsletterDestroy.Name(),
		Handler:     newsletter.Destroy,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
