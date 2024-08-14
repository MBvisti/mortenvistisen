package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Api struct {
	base Base
}

func NewApi(base Base) Api {
	return Api{base}
}

func (a Api) Health(ctx echo.Context) error {
	return ctx.JSON(
		http.StatusOK,
		fmt.Sprintf("app is healthy and running on commit: %v", a.base.Config.App.Version),
	)
}
