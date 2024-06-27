package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AppHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "app is healthy and running")
}
