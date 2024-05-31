package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func appHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, []byte("app is healthy and running"))
}
