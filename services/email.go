package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/views/emails"
)

type MailPayload struct {
	To       string
	From     string
	Subject  string
	HtmlBody string
	TextBody string
}

type MailClient interface {
	SendMail(ctx context.Context, payload MailPayload) error
}

type Email struct {
	cfg    config.Cfg
	client MailClient
}

func NewEmailSvc(
	cfg config.Cfg,
	client MailClient,
) Email {
	return Email{
		cfg,
		client,
	}
}

func (e *Email) SendNewSubscriberEmail(
	ctx context.Context,
	subscriberEmail string,
	activationToken, unsubscribeToken string,
) error {
	newsletterMail := emails.NewsletterWelcomeMail{
		ConfirmationLink: fmt.Sprintf(
			"%s://%s/verify-subscriber?token=%s",
			e.cfg.App.AppScheme,
			e.cfg.App.AppHost,
			activationToken,
		),
		UnsubscribeLink: fmt.Sprintf(
			"%s://%s/unsubscriber?token=%s",
			e.cfg.App.AppScheme,
			e.cfg.App.AppHost,
			unsubscribeToken,
		),
	}

	textVersion, err := newsletterMail.GenerateTextVersion()
	if err != nil {
		return err
	}

	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
	if err != nil {
		return err
	}

	return e.client.SendMail(ctx, MailPayload{
		To:       subscriberEmail,
		From:     "newsletter@mortenvistisen.com",
		Subject:  "MBV Blog - action required",
		HtmlBody: htmlVersion,
		TextBody: textVersion,
	})
}

func (e *Email) SendUserSignup(
	ctx context.Context,
	email string,
	activationTkn string,
) error {
	newsletterMail := emails.UserSignupWelcomeMail{
		ConfirmationLink: fmt.Sprintf(
			"%s://%s/verify-email?token=%s",
			e.cfg.App.AppScheme,
			e.cfg.App.AppHost,
			activationTkn,
		),
	}

	textVersion, err := newsletterMail.GenerateTextVersion()
	if err != nil {
		slog.ErrorContext(
			ctx,
			"could not generate text version of UserSignupWelcomeMail",
			"error",
			err,
		)
		return err
	}

	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
	if err != nil {
		slog.ErrorContext(
			ctx,
			"could not generate html version of UserSignupWelcomeMail",
			"error",
			err,
		)
		return err
	}

	return e.client.SendMail(ctx, MailPayload{
		To:       email,
		From:     "newsletter@mortenvistisen.com",
		Subject:  "MBV Blog - action required",
		HtmlBody: htmlVersion,
		TextBody: textVersion,
	})
}

func (e *Email) SendPasswordReset(
	ctx context.Context,
	email string,
	resetLink string,
) error {
	newsletterMail := emails.PasswordReset{
		ResetPasswordLink: resetLink,
	}

	textVersion, err := newsletterMail.GenerateTextVersion()
	if err != nil {
		slog.ErrorContext(
			ctx,
			"could not generate text version of PasswordReset",
			"error",
			err,
		)
		return err
	}

	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
	if err != nil {
		slog.ErrorContext(
			ctx,
			"could not generate html version of PasswordReset",
			"error",
			err,
		)
		return err
	}

	return e.client.SendMail(ctx, MailPayload{
		To:       email,
		From:     "newsletter@mortenvistisen.com",
		Subject:  "MBV Blog - action required",
		HtmlBody: htmlVersion,
		TextBody: textVersion,
	})
}

func (e *Email) Send(
	ctx context.Context,
	to,
	from,
	subject,
	textVersion,
	htmlVersion string,
) error {
	return e.client.SendMail(ctx, MailPayload{
		To:       to,
		From:     from,
		Subject:  subject,
		HtmlBody: htmlVersion,
		TextBody: textVersion,
	})
}
