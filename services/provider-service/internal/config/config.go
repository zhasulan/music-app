package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName       string `env:"APP_NAME" envDefault:"provider-service"`
	HTTPPort      int    `env:"HTTP_PORT" envDefault:"8013"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"info"`
	AudiusBaseURL string `env:"AUDIOUS_BASE_URL" envDefault:"https://api.audius.co/v1"`
	AudiusAppName string `env:"AUDIOUS_APP_NAME" envDefault:"FreedomMusic"`
	// Optional: bearer token for Audius if provided. Leave empty for public/read-only mode.
	AudiusBearer string `env:"AUDIOUS_BEARER_TOKEN" envDefault:""`
	CacheTTL     int    `env:"CACHE_TTL_SECONDS" envDefault:"300"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
