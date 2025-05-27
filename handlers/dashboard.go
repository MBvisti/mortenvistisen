package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/views/dashboard"
)

type Dashboard struct {
	db psql.Postgres
}

func newDashboard(db psql.Postgres) Dashboard {
	return Dashboard{
		db: db,
	}
}

func (d Dashboard) Index(ctx echo.Context) error {
	// Get page parameter from query string, default to 1
	pageStr := ctx.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get articles with pagination
	pageSize := 10
	result, err := models.GetArticlesPaginated(setAppCtx(ctx), d.db.Pool, page, pageSize)
	if err != nil {
		// Log error and show empty state
		result = models.PaginationResult{
			Articles:    []models.Article{},
			TotalCount:  0,
			Page:        1,
			PageSize:    pageSize,
			TotalPages:  0,
			HasNext:     false,
			HasPrevious: false,
		}
	}

	return dashboard.Home(result).Render(renderArgs(ctx))
}
