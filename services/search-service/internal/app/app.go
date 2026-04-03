package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freedom-music/search-service/internal/client"
	"github.com/freedom-music/search-service/internal/config"
	"github.com/freedom-music/search-service/internal/handler"
	"github.com/freedom-music/search-service/internal/service"
	sharedLogger "github.com/freedom-music/shared/logger"
	"github.com/freedom-music/shared/middleware"
	"github.com/freedom-music/shared/observability"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	logger *zap.Logger
	router *gin.Engine
}

func New() (*Application, error) {
	cfg := config.Load()
	log, err := sharedLogger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	observability.Startup(log)

	catClient := client.NewCatalogClient(cfg.CatalogServiceURL)
	tracks, err := catClient.Tracks()
	if err != nil {
		return nil, fmt.Errorf("load tracks: %w", err)
	}
	albums, err := catClient.Albums()
	if err != nil {
		return nil, fmt.Errorf("load albums: %w", err)
	}
	artists, err := catClient.Artists()
	if err != nil {
		return nil, fmt.Errorf("load artists: %w", err)
	}

	svc := service.New(tracks, albums, artists)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(log))

	h := handler.New(svc)
	h.Register(r)

	return &Application{cfg: cfg, logger: log, router: r}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", a.cfg.HTTPPort), Handler: a.router}

	go func() {
		a.logger.Info("search-service listening", zap.Int("port", a.cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.logger.Info("shutdown signal received")

	shutdown := make(chan struct{})
	go func() {
		_ = srv.Close()
		close(shutdown)
	}()

	select {
	case <-shutdown:
	case <-time.After(10 * time.Second):
		return fmt.Errorf("shutdown timed out")
	}
	return nil
}
