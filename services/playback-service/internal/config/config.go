package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName           string `env:"APP_NAME" envDefault:"playback-service"`
	HTTPPort          int    `env:"HTTP_PORT" envDefault:"8006"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"info"`
	RedisAddr         string `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RedisDB           int    `env:"REDIS_DB" envDefault:"0"`
	JWTSecret         string `env:"JWT_SECRET" envDefault:"supersecret"`
	CatalogServiceURL string `env:"CATALOG_SERVICE_URL" envDefault:"http://catalog-service:8003"`
	SessionTTLMinutes int    `env:"SESSION_TTL_MINUTES" envDefault:"720"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
