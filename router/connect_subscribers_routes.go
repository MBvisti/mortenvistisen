package router

import (
	"errors"
	"net/http"
	"time"

	"mortenvistisen/controllers"
	"mortenvistisen/router/middleware"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

func (r Router) RegisterSubscriberRoutes(subscriber controllers.Subscribers) error {
	errs := []error{}
	adminOnly := []echo.MiddlewareFunc{middleware.AdminOnly}

	_, err := r.e.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.SubscriberSignup.Path(),
		Name:    routes.SubscriberSignup.Name(),
		Handler: subscriber.Signup,
		Middlewares: []echo.MiddlewareFunc{
			middleware.IPRateLimiterWithBan(3, time.Hour, 24*time.Hour, routes.HomePage),
		},
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.SubscriberVerificationNew.Path(),
		Name:    routes.SubscriberVerificationNew.Name(),
		Handler: subscriber.VerificationNew,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.SubscriberVerificationCreate.Path(),
		Name:    routes.SubscriberVerificationCreate.Name(),
		Handler: subscriber.VerificationCreate,
		Middlewares: []echo.MiddlewareFunc{
			middleware.IPRateLimiterWithBan(
				5,
				time.Hour,
				24*time.Hour,
				routes.SubscriberVerificationNew,
			),
		},
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.SubscriberIndex.Path(),
		Name:        routes.SubscriberIndex.Name(),
		Handler:     subscriber.Index,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.SubscriberShow.Path(),
		Name:        routes.SubscriberShow.Name(),
		Handler:     subscriber.Show,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.SubscriberNew.Path(),
		Name:        routes.SubscriberNew.Name(),
		Handler:     subscriber.New,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        routes.SubscriberCreate.Path(),
		Name:        routes.SubscriberCreate.Name(),
		Handler:     subscriber.Create,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.SubscriberEdit.Path(),
		Name:        routes.SubscriberEdit.Name(),
		Handler:     subscriber.Edit,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        routes.SubscriberUpdate.Path(),
		Name:        routes.SubscriberUpdate.Name(),
		Handler:     subscriber.Update,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.e.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        routes.SubscriberDestroy.Path(),
		Name:        routes.SubscriberDestroy.Name(),
		Handler:     subscriber.Destroy,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
