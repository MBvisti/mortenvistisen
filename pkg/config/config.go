package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	database       configDatabase
	authentication configAuthentication
}

var Cfg Config = setupConfiguration()

func setupConfiguration() Config {
	databaseCfg := configDatabase{}
	if err := env.Parse(&databaseCfg); err != nil {
		panic(err)
	}

	authCfg := configAuthentication{}
	if err := env.Parse(&authCfg); err != nil {
		panic(err)
	}
	return Config{
		databaseCfg,
		authCfg,
	}
}

func (c Config) GetDatabaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		c.database.DatabaseKind, c.database.User, c.database.Password, c.database.Host, c.database.Port,
		c.database.Name,
	)
}

func (c Config) GetPwdPepper() string {
	return c.authentication.pwdPepper
}
