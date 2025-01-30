package handlers

import (
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/labstack/echo/v4"
)

type Dashboard struct{}

func newDashboard() Dashboard {
	return Dashboard{}
}

func (d *Dashboard) Index(c echo.Context) error {
	return dashboard.Home().Render(renderArgs(c))
}
