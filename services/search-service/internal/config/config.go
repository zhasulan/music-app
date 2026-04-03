package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName           string `env:"APP_NAME" envDefault:"search-service"`
	HTTPPort          int    `env:"HTTP_PORT" envDefault:"8011"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"info"`
	CatalogServiceURL string `env:"CATALOG_SERVICE_URL" envDefault:"http://catalog-service:8003"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
