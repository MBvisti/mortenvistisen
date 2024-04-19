package mail

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	sender = "nopreply@mortenvistisen.com"

	// The character encoding for the email.
	charSet = "UTF-8"
)

type AwsSimpleEmailService struct {
	client *ses.SES
}

// SendMail implements mailClient.
func (a *AwsSimpleEmailService) SendMail(ctx context.Context, payload MailPayload) error {
	from := payload.From
	if payload.From == "" {
		from = sender
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

	result, err := a.client.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
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

	log.Print(result)

	return nil
}

func NewAwsSimpleEmailService() AwsSimpleEmailService {
	creds := credentials.NewEnvCredentials()
	conf := &aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: creds,
	}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	// Create an SES session.
	svc := ses.New(sess)
	return AwsSimpleEmailService{
		svc,
	}
}

var _ mailClient = (*AwsSimpleEmailService)(nil)
