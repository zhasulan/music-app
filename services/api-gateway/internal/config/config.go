package config

import "github.com/freedom-music/shared/config"

// Config holds environment-driven settings for the API Gateway.
type Config struct {
	AppName            string `env:"APP_NAME" envDefault:"api-gateway"`
	HTTPPort           int    `env:"HTTP_PORT" envDefault:"8080"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
	AuthServiceURL     string `env:"AUTH_SERVICE_URL" envDefault:"http://auth-service:8001"`
	UserServiceURL     string `env:"USER_SERVICE_URL" envDefault:"http://user-service:8002"`
	CatalogServiceURL  string `env:"CATALOG_SERVICE_URL" envDefault:"http://catalog-service:8003"`
	PlaylistServiceURL string `env:"PLAYLIST_SERVICE_URL" envDefault:"http://playlist-service:8004"`
	LibraryServiceURL  string `env:"LIBRARY_SERVICE_URL" envDefault:"http://library-service:8005"`
	PlaybackServiceURL string `env:"PLAYBACK_SERVICE_URL" envDefault:"http://playback-service:8006"`
	MediaServiceURL    string `env:"MEDIA_SERVICE_URL" envDefault:"http://media-service:8007"`
	EventsServiceURL   string `env:"EVENTS_SERVICE_URL" envDefault:"http://events-service:8010"`
	SearchServiceURL   string `env:"SEARCH_SERVICE_URL" envDefault:"http://search-service:8011"`
	RecommendationURL  string `env:"RECOMMENDATION_SERVICE_URL" envDefault:"http://recommendation-service:8012"`
	ProviderServiceURL string `env:"PROVIDER_SERVICE_URL" envDefault:"http://provider-service:8013"`
}

func Load() *Config {
	return config.MustLoad[Config]()
}
