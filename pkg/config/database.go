package config

type configDatabase struct {
	Port         string `env:"PORT"`
	Host         string `env:"HOST"`
	Name         string `env:"NAME"`
	User         string `env:"USER"`
	Password     string `env:"PASSWORD"`
	DatabaseKind string `env:"DATABASE_KIND"`
}
