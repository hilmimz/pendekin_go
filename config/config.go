package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME" envDefault:"mydb"`
}

type AppConfig struct {
	AppName      string `env:"APP_NAME" envDefault:"localhost:8080"`
	AppEnv       string `env:"APP_ENV" envDefault:"development"`
	AliasLength  int    `env:"ALIAS_LENGTH" envDefault:"6"`
	ExpiresIn    int    `env:"EXPIRES_IN" envDefault:"24"`
	JWTSecret    string `env:"JWT_SECRET" envDefault:"secret"`
	JWTExpiresIn int    `env:"JWT_EXPIRES_IN" envDefault:"86400"`
	FrontendURL  string `env:"FRONTEND_URL" envDefault:"localhost:5173"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
