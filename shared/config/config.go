package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

// Load parses environment variables into the provided struct type.
// Example:
// type Config struct { HTTPPort int `env:"HTTP_PORT,default=8080"` }
// cfg, err := config.Load[Config]()
func Load[T any]() (*T, error) {
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}
	return &cfg, nil
}

// MustLoad panics on error. Use only for process boot where failure should stop startup.
func MustLoad[T any]() *T {
	cfg, err := Load[T]()
	if err != nil {
		panic(err)
	}
	return cfg
}
