package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Api struct{}

func newApi() Api {
	return Api{}
}

func (a *Api) AppHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "app is healthy and running")
}
