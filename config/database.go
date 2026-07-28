package config

import "github.com/caarlos0/env/v11"

type Database struct {
	DatabaseURL string `env:"DATABASE_URL"`
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
