package cookies

import (
	"strings"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/internal/server"

	"github.com/gosimple/slug"
	"github.com/labstack/echo-contrib/v5/session"
	"github.com/labstack/echo/v5"
	"github.com/rs/xid"
)

type FlashMessage struct {
	ID        xid.ID
	Type      FlashType
	CreatedAt time.Time
	Message   string
}

var (
	flashSession = func() string {
		if config.Env == server.ProdEnvironment {
			return slug.Make(strings.ToLower(config.ProjectName) + " " + "flash_key")
		}

		return slug.Make(strings.ToLower(config.ProjectName) + " " + "dev_flash_key")
	}()
	flashSessionName = "flash_session"
)

type FlashType string

const (
	FlashSuccess FlashType = "success"
	FlashError   FlashType = "error"
	FlashWarning FlashType = "warning"
	FlashInfo    FlashType = "info"
)

func AddFlash(
	c *echo.Context, flashType FlashType, msg string,
) error {
	sess, err := session.Get(flashSession, c)
	if err != nil {
		return err
	}

	sess.AddFlash(FlashMessage{
		ID:        xid.New(),
		Type:      flashType,
		CreatedAt: time.Now(),
		Message:   msg,
	}, flashSessionName)

	return sess.Save(c.Request(), c.Response())
}

func ExtractFlashes(c *echo.Context) ([]FlashMessage, error) {
	sess, err := session.Get(flashSession, c)
	if err != nil {
		return nil, err
	}

	var flashMessages []FlashMessage
	for _, flash := range sess.Flashes(flashSessionName) {
		if msg, ok := flash.(FlashMessage); ok {
			flashMessages = append(flashMessages, msg)
		}
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return nil, err
	}

	return flashMessages, nil
}
