package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freedom-music/media-service/internal/client"
	"github.com/freedom-music/media-service/internal/config"
	"github.com/freedom-music/media-service/internal/handler"
	"github.com/freedom-music/media-service/internal/service"
	"github.com/freedom-music/media-service/internal/storage"
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

	store, err := storage.NewMinioStorage(cfg.MinioEndpoint, cfg.MinioAccess, cfg.MinioSecret, cfg.MinioBucket, cfg.UseSSL, cfg.PublicURL, time.Duration(cfg.URLExpiryMinutes)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("storage init: %w", err)
	}
	if err := store.EnsureBucket(context.Background()); err != nil {
		return nil, fmt.Errorf("ensure bucket: %w", err)
	}

	catalogClient := client.NewCatalogClient(cfg.CatalogServiceURL)
	svc := service.NewMediaService(store, catalogClient)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(log))

	h := handler.NewMediaHandler(svc)
	h.RegisterRoutes(r)

	return &Application{cfg: cfg, logger: log, router: r}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", a.cfg.HTTPPort), Handler: a.router}

	go func() {
		a.logger.Info("media-service listening", zap.Int("port", a.cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
