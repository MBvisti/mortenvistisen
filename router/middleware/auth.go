package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/internal/routing"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
	echomw "github.com/labstack/echo/v5/middleware"
	"github.com/maypok86/otter/v2"
)

func AuthOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if cookies.ExtractFromCookieApp(c).IsAuthenticated {
			return next(c)
		}

		return c.Redirect(http.StatusSeeOther, routes.SessionNew.URL())
	}
}

func APIBasicAuth(cfg config.Config) echo.MiddlewareFunc {
	return echomw.BasicAuthWithConfig(echomw.BasicAuthConfig{
		Realm: "mortenvistisen API",
		Validator: func(_ *echo.Context, username string, password string) (bool, error) {
			expectedUsername := cfg.App.APIBasicAuthUsername
			expectedPassword := cfg.App.APIBasicAuthPassword
			if expectedUsername == "" || expectedPassword == "" {
				return false, nil
			}

			return constantTimeStringEqual(username, expectedUsername) &&
				constantTimeStringEqual(password, expectedPassword), nil
		},
	})
}

func constantTimeStringEqual(actual string, expected string) bool {
	actualHash := sha256.Sum256([]byte(actual))
	expectedHash := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(actualHash[:], expectedHash[:]) == 1
}

func IPRateLimiter(
	limit int32,
	redirectURL routing.Route,
) func(next echo.HandlerFunc) echo.HandlerFunc {
	cache := otter.Must(&otter.Options[string, int32]{
		MaximumSize:      1000,
		ExpiryCalculator: otter.ExpiryCreating[string, int32](10 * time.Minute),
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()
			allowed := false
			cache.Compute(ip, func(hits int32, found bool) (int32, otter.ComputeOp) {
				if !found {
					hits = 0
				}
				if hits >= limit {
					return hits, otter.CancelOp
				}

				allowed = true
				return hits + 1, otter.WriteOp
			})

			if !allowed {
				return c.String(http.StatusTooManyRequests, redirectURL.URL())
			}

			return next(c)
		}
	}
}
