package mail_client

import (
	"context"
	"log"
	"os"

	"github.com/MBvisti/mortenvistisen/services"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	sender = "nopreply@mortenvistisen.com"

	// The character encoding for the email.
	charSet = "UTF-8"
)

type AwsSimpleEmailService struct {
	client *ses.Client
}

func (a *AwsSimpleEmailService) GetStatistics(ctx context.Context) error {
	input := &ses.GetSendStatisticsInput{}
	stats, err := a.client.GetSendStatistics(ctx, input)
	if err != nil {
		return err
	}

	for _, dp := range stats.SendDataPoints {
		log.Println(dp)
	}

	return nil
}

// SendMail implements mailClient.
func (a *AwsSimpleEmailService) SendMail(ctx context.Context, payload services.MailPayload) error {
	from := payload.From
	if payload.From == "" {
		from = sender
	}
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				payload.To,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(payload.HtmlBody),
				},
				Text: &types.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(payload.TextBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(payload.Subject),
			},
		},
		Source: aws.String(from),
	}

	_, err := a.client.SendEmail(ctx, input)
	if err != nil {
		// if aerr, ok := err.(awserr.Error); ok {
		// 	switch aerr.Code() {
		// 	case ses.ErrCodeMessageRejected:
		// 		fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
		// 	case ses.ErrCodeMailFromDomainNotVerifiedException:
		// 		fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		// 	case ses.ErrCodeConfigurationSetDoesNotExistException:
		// 		fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		// 	default:
		// 		fmt.Println(aerr.Error())
		// 	}
		// } else {
		// 	// Print the error, cast err to awserr.Error to get the Code and
		// 	// Message from an error.
		// 	fmt.Println(err.Error())
		// }

		return err
	}

	return nil
}

func NewAwsSimpleEmailService() AwsSimpleEmailService {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	amazonConfiguration, createAmazonConfigurationError := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey, secretKey, "",
			),
		),
	)

	if createAmazonConfigurationError != nil {
		panic(createAmazonConfigurationError)
	}

	sesClient := ses.NewFromConfig(amazonConfiguration)

	// creds := credentials.NewEnvCredentials()
	// conf := &aws.Config{
	// 	Region:      aws.String("eu-central-1"),
	// 	Credentials: creds,
	// }
	// sess, err := session.NewSession(conf)
	// if err != nil {
	// 	panic(err)
	// }

	// Create an SES session.
	// svc := ses.New(sess)
	return AwsSimpleEmailService{
		sesClient,
	}
}
