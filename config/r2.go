package config

import "github.com/caarlos0/env/v11"

type r2 struct {
	AccountID       string `env:"R2_ACCOUNT_ID"        envDefault:"1e47e69178115fb8a1b60e498ac7f0a1"`
	Bucket          string `env:"R2_BUCKET"            envDefault:"mbv-blog"`
	PublicBaseURL   string `env:"R2_PUBLIC_BASE_URL"   envDefault:"https://media.mortenvistisen.com"`
	CoverPrefix     string `env:"R2_COVER_PREFIX"      envDefault:"covers"`
	AccessKeyID     string `env:"R2_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"R2_SECRET_ACCESS_KEY"`
	MaxUploadBytes  int64  `env:"R2_MAX_UPLOAD_BYTES"  envDefault:"10485760"`
}

func newR2Config() r2 {
	cfg := r2{}

	if err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	return cfg
}
