package controllers

import (
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/labstack/echo/v4"
)

func (c *Controller) Projects(ctx echo.Context) error {
	return views.ProjectsPage(views.Head{
		Title:       "Projects",
		Description: "A collection of on-going and retired projects I've build",
		Slug:        c.buildURLFromSlug("projects"),
		MetaType:    "website",
	}).Render(views.ExtractRenderDeps(ctx))
}
