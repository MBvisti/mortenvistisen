package controllers

import (
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/labstack/echo/v4"
)

func (c *Controller) Projects(ctx echo.Context) error {
	return views.ProjectsPage(views.Head{
		Title:       "Projects",
		Description: "A collection of on-going and retired projects I've build over the years. This includes business projects, open source projects and personal projects.",
		Slug:        c.buildURLFromSlug("projects"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}).Render(views.ExtractRenderDeps(ctx))
}
