package services

import (
	"context"
	"fmt"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/emails"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type MailPayload struct {
	To       string
	From     string
	Subject  string
	HtmlBody string
	TextBody string
}

const (
	charSet   = "UTF-8"
	awsRegion = "eu-central-1"
)

var defaultSender = config.Cfg.DefaultSenderSignature

type Mail struct {
	client *ses.SES
}

func NewMail(
// client MailClient,
) Mail {
	creds := credentials.NewEnvCredentials()
	conf := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: creds,
	}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	ses := ses.New(sess)

	return Mail{
		ses,
	}
}

func (m *Mail) SendNewSubscriber(
	ctx context.Context,
	subscriberEmail string,
	activationToken models.Token,
	unsubscribeToken models.Token,
) error {
	newsletterMail := emails.NewsletterWelcomeMail{
		ConfirmationLink: fmt.Sprintf(
			"%s/verify-subscriber?token=%s",
			config.Cfg.GetFullDomain(),
			activationToken.Hash,
		),
		UnsubscribeLink: fmt.Sprintf(
			"%s/unsubscriber?token=%s",
			config.Cfg.GetFullDomain(),
			unsubscribeToken.Hash,
		),
	}

	htmlVersion, textVersion, err := newsletterMail.Generate(ctx)
	if err != nil {
		return err
	}

	return m.send(ctx, MailPayload{
		To:       subscriberEmail,
		From:     "newsletter@mortenvistisen.com",
		Subject:  "Morten Vistisen Newsletter - action required",
		HtmlBody: htmlVersion.String(),
		TextBody: textVersion.String(),
	})
}

// func (e *Email) SendUserSignup(
// 	ctx context.Context,
// 	email string,
// 	activationTkn string,
// ) error {
// 	newsletterMail := emails.UserSignupWelcomeMail{
// 		ConfirmationLink: fmt.Sprintf(
// 			"%s://%s/verify-email?token=%s",
// 			e.cfg.App.AppScheme,
// 			e.cfg.App.AppHost,
// 			activationTkn,
// 		),
// 	}
//
// 	textVersion, err := newsletterMail.GenerateTextVersion()
// 	if err != nil {
// 		slog.ErrorContext(
// 			ctx,
// 			"could not generate text version of UserSignupWelcomeMail",
// 			"error",
// 			err,
// 		)
// 		return err
// 	}
//
// 	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
// 	if err != nil {
// 		slog.ErrorContext(
// 			ctx,
// 			"could not generate html version of UserSignupWelcomeMail",
// 			"error",
// 			err,
// 		)
// 		return err
// 	}
//
// 	return e.client.SendMail(ctx, MailPayload{
// 		To:       email,
// 		From:     "newsletter@mortenvistisen.com",
// 		Subject:  "MBV Blog - action required",
// 		HtmlBody: htmlVersion,
// 		TextBody: textVersion,
// 	})
// }
//
// func (e *Email) SendPasswordReset(
// 	ctx context.Context,
// 	email string,
// 	resetLink string,
// ) error {
// 	newsletterMail := emails.PasswordReset{
// 		ResetPasswordLink: resetLink,
// 	}
//
// 	textVersion, err := newsletterMail.GenerateTextVersion()
// 	if err != nil {
// 		slog.ErrorContext(
// 			ctx,
// 			"could not generate text version of PasswordReset",
// 			"error",
// 			err,
// 		)
// 		return err
// 	}
//
// 	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
// 	if err != nil {
// 		slog.ErrorContext(
// 			ctx,
// 			"could not generate html version of PasswordReset",
// 			"error",
// 			err,
// 		)
// 		return err
// 	}
//
// 	return e.client.SendMail(ctx, MailPayload{
// 		To:       email,
// 		From:     "newsletter@mortenvistisen.com",
// 		Subject:  "MBV Blog - action required",
// 		HtmlBody: htmlVersion,
// 		TextBody: textVersion,
// 	})
// }
//
// func (e *Email) SendNewsletter(
// 	ctx context.Context,
// 	title string,
// 	content string,
// 	email string,
// 	UnsubscribeTkn string,
// ) error {
// 	parsedContent, err := e.mdParser.ParseContent(content)
// 	if err != nil {
// 		return err
// 	}
//
// 	newsletterMail := emails.NewsletterMail{
// 		Title:   title,
// 		Content: parsedContent,
// 		UnsubscribeLink: fmt.Sprintf(
// 			"%s://%s/unsubscriber?token=%s",
// 			e.cfg.App.AppScheme,
// 			e.cfg.App.AppHost,
// 			UnsubscribeTkn,
// 		),
// 	}
//
// 	textVersion, err := newsletterMail.GenerateTextVersion()
// 	if err != nil {
// 		slog.ErrorContext(
// 			ctx,
// 			"could not generate text version of UserSignupWelcomeMail",
// 			"error",
// 			err,
// 		)
// 		return err
// 	}
//
// 	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
// 	if err != nil {
// 		slog.ErrorContext(
// 			ctx,
// 			"could not generate html version of UserSignupWelcomeMail",
// 			"error",
// 			err,
// 		)
// 		return err
// 	}
//
// 	return e.client.SendMail(ctx, MailPayload{
// 		To:       email,
// 		From:     "newsletter@mortenvistisen.com",
// 		Subject:  "MBV Blog - Newsletter",
// 		HtmlBody: htmlVersion,
// 		TextBody: textVersion,
// 	})
// }

func (m *Mail) send(
	ctx context.Context,
	payload MailPayload,
) error {
	from := payload.From
	if payload.From == "" {
		from = defaultSender
	}
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(payload.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(payload.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(payload.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(payload.Subject),
			},
		},
		Source: aws.String(from),
	}

	_, err := m.client.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(
					ses.ErrCodeMailFromDomainNotVerifiedException,
					aerr.Error(),
				)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(
					ses.ErrCodeConfigurationSetDoesNotExistException,
					aerr.Error(),
				)
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return err
	}

	return nil
}
