package config

import (
	"github.com/caarlos0/env/v10"
)

type Cfg struct {
	Db                Database
	Auth              Authentication
	App               App
	ExternalProviders ExternalProviders
	Telemetry         Telemetry
}

func New() Cfg {
	databaseCfg := Database{}
	if err := env.ParseWithOptions(&databaseCfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	authCfg := Authentication{}
	if err := env.ParseWithOptions(&authCfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	appCfg := App{}
	if err := env.ParseWithOptions(&appCfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	externalProviders := ExternalProviders{}
	if err := env.ParseWithOptions(&externalProviders, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	return Cfg{
		databaseCfg,
		authCfg,
		appCfg,
		externalProviders,
		newTelemetry(),
	}
}
