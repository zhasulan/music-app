package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName  string `env:"APP_NAME" envDefault:"catalog-service"`
	HTTPPort int    `env:"HTTP_PORT" envDefault:"8003"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
