package config

import (
	"net"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Database struct {
	Kind     string `env:"DB_KIND"     envDefault:"postgres"`
	Port     string `env:"DB_PORT"`
	Host     string `env:"DB_HOST"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

func (d Database) URL() string {
	return PostgresURL(d.Host, d.Port, d.Name, d.User, d.Password, d.SSLMode)
}

// PostgresURL builds a postgres connection string from individual parts.
func PostgresURL(host, port, name, user, password, sslMode string) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + name,
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func newDatabaseConfig() Database {
	dataCfg := Database{}

	if err := env.ParseWithOptions(&dataCfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	return dataCfg
}
