package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/emails"
)

func main() {
	ctx := context.Background()
	e := echo.New()

	passwordReset := emails.PasswordReset{
		ResetLink: fmt.Sprintf(
			"%s/%s?token=%s",
			config.Cfg.GetFullDomain(),
			"reset-password",
			"wvSwI8Yq02o9cmJ6zVSTkP44lXGJZjmMF8v10vxAhrrV6UyzRr59ogUzdo3VKP7y",
		),
	}
	passwordResetHtml, passwordResetText, _ := passwordReset.Generate(ctx)

	signupWelcome := emails.SignupWelcome{
		VerificationCode: "43dd1w",
	}
	signupWelcomeHtml, signupWelcomeText, _ := signupWelcome.Generate(ctx)

	newsletterHtml, newsletterTxt, _ := emails.Newsletter{
		Subject:         "YoYo",
		Content:         "<h1>BUY</h1><p>Hello</p>",
		UnsubscribeLink: "https://mbvlabs.com",
	}.Generate(
		ctx,
	)

	textGroup := e.Group("/text")
	textGroup.GET("/password-reset", func(c echo.Context) error {
		return c.String(http.StatusOK, passwordResetText.String())
	})
	textGroup.GET("/signup-welcome", func(c echo.Context) error {
		return c.String(http.StatusOK, signupWelcomeText.String())
	})
	textGroup.GET("/newsletter", func(c echo.Context) error {
		return c.String(http.StatusOK, newsletterTxt.String())
	})

	htmlGroup := e.Group("/html")
	htmlGroup.GET("/password-reset", func(c echo.Context) error {
		return c.HTML(http.StatusOK, passwordResetHtml.String())
	})
	htmlGroup.GET("/signup-welcome", func(c echo.Context) error {
		return c.HTML(http.StatusOK, signupWelcomeHtml.String())
	})
	htmlGroup.GET("/newsletter", func(c echo.Context) error {
		return c.HTML(http.StatusOK, newsletterHtml.String())
	})

	slog.Info("starting the password server on port: 4444")
	log.Fatal(e.Start(":4444"))
}
