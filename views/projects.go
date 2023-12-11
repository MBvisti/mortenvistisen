package views

import (
	"github.com/MBvisti/grafto/views/internal/components"
	"github.com/MBvisti/grafto/views/internal/layouts"
	"github.com/MBvisti/grafto/views/internal/pages"
	"github.com/labstack/echo/v4"
)

func Projects(ctx echo.Context) error {
	header := components.Head{
		Title:       "Projects made by mortenvistisen.com",
		Description: "Collection of all my projects; both deprecated and on-going.",
		Slug:        "https://mortenvistisen/projects",
		MetaType:    "website",
	}
	return layouts.Base(pages.Projects(), header).Render(extractRenderDeps(ctx))
}
