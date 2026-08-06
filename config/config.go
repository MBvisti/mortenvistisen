// Package config provides application-wide configuration settings.
package config

import (
	"fmt"
	"os"
	"strings"

	"mortenvistisen/internal/server"

	"github.com/gosimple/slug"

	"go.uber.org/fx"
)

// Global application settings that can be used throughout the codebase with defaults.
var (
	Env = func() string {
		if os.Getenv("ENVIRONMENT") != "" {
			return os.Getenv("ENVIRONMENT")
		}

		return server.DevEnvironment
	}()
	ProjectName = func() string {
		if os.Getenv("PROJECT_NAME") != "" {
			return os.Getenv("PROJECT_NAME")
		}

		return "andurel"
	}()
	ServiceName = func() string {
		if os.Getenv("TELEMETRY_SERVICE_NAME") != "" {
			return os.Getenv("TELEMETRY_SERVICE_NAME")
		}

		return slug.Make(ProjectName)
	}()
	Domain = func() string {
		if os.Getenv("DOMAIN") != "" {
			return os.Getenv("DOMAIN")
		}

		return "localhost:8080"
	}()
	BaseURL = func() string {
		var protocol string

		if os.Getenv("PROTOCOL") != "" {
			protocol = os.Getenv("PROTOCOL")
		} else {
			protocol = "http"
		}

		return fmt.Sprintf("%s://%s", protocol, Domain)
	}()
	AppCookieSessionName = func() string {
		return "app sess " + slug.Make(strings.ToLower(ProjectName)) + " " + Env
	}()
	DefaultSenderSignature = func() string {
		if os.Getenv("DEFAULT_SENDER_SIGNATURE") != "" {
			return os.Getenv("DEFAULT_SENDER_SIGNATURE")
		}

		return "noreply@" + Domain
	}()
)

type Config struct {
	App       app
	DB        Database
	Telemetry telemetry
	Email     email
	AwsSes    awsSes
	R2        r2
	Auth      auth
}

func NewConfig() Config {
	return Config{
		App:       newAppConfig(),
		DB:        newDatabaseConfig(),
		Telemetry: newTelemetryConfig(),
		Email:     newEmailConfig(),
		AwsSes:    newAwsSesConfig(),
		R2:        newR2Config(),
		Auth:      newAuthConfig(),
	}
}

var Module = fx.Module("config", fx.Provide(NewConfig))
