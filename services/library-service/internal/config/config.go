package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName           string `env:"APP_NAME" envDefault:"library-service"`
	HTTPPort          int    `env:"HTTP_PORT" envDefault:"8005"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL       string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@postgres:5432/library_service?sslmode=disable"`
	JWTSecret         string `env:"JWT_SECRET" envDefault:"supersecret"`
	CatalogServiceURL string `env:"CATALOG_SERVICE_URL" envDefault:"http://catalog-service:8003"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
