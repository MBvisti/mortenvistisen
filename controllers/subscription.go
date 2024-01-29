package controllers

import (
	"log"

	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/labstack/echo/v4"
)

type SubscriptionEventForm struct {
	Email string `form:"hero-input"`
	Title string `form:"article-title"`
}

func (c *Controller) SubscriptionEvent(ctx echo.Context) error {
	var form SubscriptionEventForm
	if err := ctx.Bind(&form); err != nil {
		if err := c.mail.Send(ctx.Request().Context(), "hi@mortenvistisen.com", "sub-blog@mortenvistisen.com", "Failed to subscribe", "sub_report", err.Error()); err != nil {
			telemetry.Logger.Error("Failed to send email", "error", err)
		}
		return ctx.String(200, "You're now subscribed!")
	}
	log.Println(form.Email)
	if err := c.mail.Send(ctx.Request().Context(), "hi@mortenvistisen.com", "sub-blog@mortenvistisen.com", "New subscriber", "sub_report", form); err != nil {
		telemetry.Logger.Error("Failed to send email", "error", err)
	}

	return ctx.String(200, "You're now subscribed!")
}

func (c *Controller) RemoveSubscriptionEvent(ctx echo.Context) error {
	emailToRemove := ctx.QueryParam("email")
	form := SubscriptionEventForm{
		Email: emailToRemove,
		Title: "Remove this email was removed",
	}
	if err := c.mail.Send(ctx.Request().Context(), "hi@mortenvistisen.com", "sub-blog@mortenvistisen.com", "Subscriber asked to be removed", "sub_report", form); err != nil {
		telemetry.Logger.Error("Failed to send email", "error", err)
	}

	return ctx.String(200, "Your email has been removed from the list!")
}
