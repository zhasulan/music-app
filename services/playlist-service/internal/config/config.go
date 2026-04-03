package config

import "github.com/freedom-music/shared/config"

// Config holds environment settings for playlist-service.
type Config struct {
	AppName           string `env:"APP_NAME" envDefault:"playlist-service"`
	HTTPPort          int    `env:"HTTP_PORT" envDefault:"8004"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL       string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@postgres:5432/playlist_service?sslmode=disable"`
	JWTSecret         string `env:"JWT_SECRET" envDefault:"supersecret"`
	CatalogServiceURL string `env:"CATALOG_SERVICE_URL" envDefault:"http://catalog-service:8003"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
