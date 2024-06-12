package services

import (
	"context"
	"fmt"

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

type mailClient interface {
	SendMail(ctx context.Context, payload MailPayload) error
}

type EmailSvc struct {
	cfg    config.Cfg
	client mailClient
}

func NewEmailSvc(
	cfg config.Cfg,
	client mailClient,
) EmailSvc {
	return EmailSvc{
		cfg,
		client,
	}
}

func (e *EmailSvc) SendNewSubscriberEmail(
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

func (e *EmailSvc) Send(
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
