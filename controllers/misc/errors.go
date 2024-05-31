package misc

import (
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/labstack/echo/v4"
)

func InternalError(ctx echo.Context) error {
	from := "/"

	return views.InternalServerErr(ctx, views.InternalServerErrData{
		FromLocation: from,
	})
}
