package dashboard

import (
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/labstack/echo/v4"
)

func Index(ctx echo.Context, emailService services.Email) error {
	return dashboard.Index().Render(views.ExtractRenderDeps(ctx))
}
