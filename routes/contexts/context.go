package contexts

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AppKey struct{}

func (AppKey) String() string {
	return ""
}

type App struct {
	echo.Context
	UserID          uuid.UUID
	Email           string
	IsAuthenticated bool
	IsAdmin         bool
	CurrentPath     string
}

func ExtractApp(ctx context.Context) *App {
	appCtx, ok := ctx.Value(AppKey{}).(*App)
	if !ok {
		return &App{}
	}

	return appCtx
}

type FlashKey struct{}

func (FlashKey) String() string {
	return ""
}

// FlashType represents the type of flash message
type FlashType string

const (
	FlashSuccess FlashType = "success"
	FlashError   FlashType = "error"
	FlashWarning FlashType = "warning"
	FlashInfo    FlashType = "info"
)

// FlashMessage represents a single flash message
type FlashMessage struct {
	echo.Context
	ID        uuid.UUID
	Type      FlashType
	CreatedAt time.Time
	Message   string
}
