package config

type ExternalProviders struct {
	// PostmarkApiToken   string `env:"POSTMARK_API_TOKEN"`
	AwsAccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
}
