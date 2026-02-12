// Package controllers provides HTTP handlers for the web application.
package controllers

import (
	"context"
	"io"
	"mortenvistisen/internal/renderer"
	"mortenvistisen/router/cookies"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

func render(etx *echo.Context, t templ.Component) error {
	pathAwareComponent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		withPathCtx := renderer.WithRequestPath(ctx, etx.Request().URL.Path)
		return t.Render(withPathCtx, w)
	})

	return renderer.Render(
		etx,
		pathAwareComponent,
		[]renderer.CookieKey{
			cookies.AppKey,
			cookies.FlashKey,
		},
	)
}
