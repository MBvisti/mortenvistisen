package middleware

import (
	"net/http"
	"time"

	"mortenvistisen/internal/routing"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
	"github.com/maypok86/otter/v2"
)

type ipRateLimitState struct {
	Hits   int32
	Banned bool
}

func AuthOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if cookies.GetApp(c).IsAuthenticated {
			return next(c)
		}

		return c.Redirect(http.StatusSeeOther, routes.SessionNew.URL())
	}
}

func IPRateLimiter(
	limit int32,
	redirectURL routing.Route,
) func(next echo.HandlerFunc) echo.HandlerFunc {
	return ipRateLimiter(limit, 10*time.Minute, 10*time.Minute, redirectURL)
}

func IPRateLimiterWithBan(
	limit int32,
	window time.Duration,
	banDuration time.Duration,
	redirectURL routing.Route,
) func(next echo.HandlerFunc) echo.HandlerFunc {
	return ipRateLimiter(limit, window, banDuration, redirectURL)
}

func ipRateLimiter(
	limit int32,
	window time.Duration,
	banDuration time.Duration,
	redirectURL routing.Route,
) func(next echo.HandlerFunc) echo.HandlerFunc {
	cache := otter.Must(&otter.Options[string, ipRateLimitState]{
		MaximumSize: 1000,
		ExpiryCalculator: otter.ExpiryCreatingFunc(
			func(entry otter.Entry[string, ipRateLimitState]) time.Duration {
				if entry.Value.Banned {
					return banDuration
				}

				return window
			},
		),
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()

			hits, ok := cache.GetIfPresent(ip)
			if !ok {
				cache.Set(ip, ipRateLimitState{
					Hits: 1,
				})
				return next(c)
			}

			if hits.Banned {
				return c.Redirect(http.StatusTooManyRequests, redirectURL.URL())
			}

			if hits.Hits >= limit {
				cache.Set(ip, ipRateLimitState{
					Hits:   hits.Hits + 1,
					Banned: true,
				})
				return c.Redirect(http.StatusTooManyRequests, redirectURL.URL())
			}

			hits.Hits++
			cache.Set(ip, hits)

			return next(c)
		}
	}
}
