package config

type configDatabase struct {
	Port         string `env:"PORT"`
	Host         string `env:"HOST"`
	Name         string `env:"NAME"`
	User         string `env:"DB_USER"`
	Password     string `env:"PASSWORD"`
	DatabaseKind string `env:"DATABASE_KIND"`
	SSL_MODE     string `env:"SSL_MODE"`
}
