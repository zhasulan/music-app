package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freedom-music/playlist-service/internal/client"
	"github.com/freedom-music/playlist-service/internal/config"
	"github.com/freedom-music/playlist-service/internal/handler"
	"github.com/freedom-music/playlist-service/internal/repository"
	"github.com/freedom-music/playlist-service/internal/service"
	sharedLogger "github.com/freedom-music/shared/logger"
	"github.com/freedom-music/shared/middleware"
	"github.com/freedom-music/shared/observability"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	logger *zap.Logger
	router *gin.Engine
	db     *pgxpool.Pool
}

func New() (*Application, error) {
	cfg := config.Load()
	log, err := sharedLogger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	observability.Startup(log)

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	catalogClient := client.NewCatalogClient(cfg.CatalogServiceURL)
	repo := repository.NewPlaylistRepository(db)
	svc := service.NewPlaylistService(repo, catalogClient)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(log))

	h := handler.NewPlaylistHandler(svc)
	h.RegisterRoutes(r, cfg.JWTSecret)

	return &Application{cfg: cfg, logger: log, router: r, db: db}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", a.cfg.HTTPPort), Handler: a.router}

	go func() {
		a.logger.Info("playlist-service listening", zap.Int("port", a.cfg.HTTPPort))
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
	a.db.Close()
	return nil
}
