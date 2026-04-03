package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freedom-music/playback-service/internal/client"
	"github.com/freedom-music/playback-service/internal/config"
	"github.com/freedom-music/playback-service/internal/handler"
	"github.com/freedom-music/playback-service/internal/service"
	"github.com/freedom-music/playback-service/internal/store"
	sharedLogger "github.com/freedom-music/shared/logger"
	"github.com/freedom-music/shared/middleware"
	"github.com/freedom-music/shared/observability"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	logger *zap.Logger
	router *gin.Engine
	redis  *redis.Client
}

func New() (*Application, error) {
	cfg := config.Load()
	log, err := sharedLogger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	observability.Startup(log)

	rds := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, DB: cfg.RedisDB})
	if err := rds.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	catalogClient := client.NewCatalogClient(cfg.CatalogServiceURL)
	st := store.NewRedisStore(rds, time.Duration(cfg.SessionTTLMinutes)*time.Minute)
	svc := service.NewPlaybackService(st, catalogClient)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(log))

	h := handler.NewPlaybackHandler(svc)
	h.RegisterRoutes(r, cfg.JWTSecret)

	return &Application{cfg: cfg, logger: log, router: r, redis: rds}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", a.cfg.HTTPPort), Handler: a.router}

	go func() {
		a.logger.Info("playback-service listening", zap.Int("port", a.cfg.HTTPPort))
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
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	if err := a.redis.Close(); err != nil {
		return err
	}
	return nil
}
