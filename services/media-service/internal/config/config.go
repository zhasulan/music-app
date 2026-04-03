package config

import "github.com/freedom-music/shared/config"

type Config struct {
	AppName           string `env:"APP_NAME" envDefault:"media-service"`
	HTTPPort          int    `env:"HTTP_PORT" envDefault:"8007"`
	LogLevel          string `env:"LOG_LEVEL" envDefault:"info"`
	MinioEndpoint     string `env:"MINIO_ENDPOINT" envDefault:"minio:9000"`
	MinioAccess       string `env:"MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	MinioSecret       string `env:"MINIO_SECRET_KEY" envDefault:"minioadmin"`
	MinioBucket       string `env:"MINIO_BUCKET" envDefault:"music"`
	UseSSL            bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	PublicURL         string `env:"MINIO_PUBLIC_URL" envDefault:"http://localhost:9000"`
	CatalogServiceURL string `env:"CATALOG_SERVICE_URL" envDefault:"http://catalog-service:8003"`
	URLExpiryMinutes  int    `env:"URL_EXPIRY_MINUTES" envDefault:"60"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
