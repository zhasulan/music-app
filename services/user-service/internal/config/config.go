package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName     string `env:"APP_NAME" envDefault:"user-service"`
	HTTPPort    int    `env:"HTTP_PORT" envDefault:"8002"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@postgres:5432/user_service?sslmode=disable"`
	JWTSecret   string `env:"JWT_SECRET" envDefault:"dev-secret"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
