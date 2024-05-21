package misc

import (
	"fmt"

	"github.com/MBvisti/mortenvistisen/views"
	"github.com/labstack/echo/v4"
)

func InternalError(ctx echo.Context) error {
	from := "/"

	return views.InternalServerErr(ctx, views.InternalServerErrData{
		FromLocation: from,
	})
}

func Redirect(ctx echo.Context) error {
	toLocation := ctx.QueryParam("to")
	if toLocation == "" {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		return InternalError(ctx)
	}

	ctx.Response().Writer.Header().Add("HX-Redirect", fmt.Sprintf("/%s", toLocation))

	return nil
}
