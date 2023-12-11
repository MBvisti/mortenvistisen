package controllers

import (
	"github.com/MBvisti/grafto/views"
	"github.com/labstack/echo/v4"
)

func (c *Controller) Projects(ctx echo.Context) error {
	return views.Projects(ctx)
}
