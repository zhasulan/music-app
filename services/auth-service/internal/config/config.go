package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName            string `env:"APP_NAME" envDefault:"auth-service"`
	HTTPPort           int    `env:"HTTP_PORT" envDefault:"8001"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL        string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@postgres:5432/auth_service?sslmode=disable"`
	JWTSecret          string `env:"JWT_SECRET" envDefault:"dev-secret"`
	AccessTokenMinutes int    `env:"ACCESS_TOKEN_MINUTES" envDefault:"15"`
	RefreshTokenHours  int    `env:"REFRESH_TOKEN_HOURS" envDefault:"720"` // 30 days
}

func Load() *Config {
	return config.MustLoad[Config]()
}
