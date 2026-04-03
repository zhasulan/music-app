package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freedom-music/api-gateway/internal/config"
	"github.com/freedom-music/api-gateway/internal/handler"
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

	r := gin.New()
	// Avoid automatic 307 redirects that can break browser XHRs behind the proxy.
	r.RedirectTrailingSlash = false
	r.RemoveExtraSlash = true
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(log))

	rt, err := handler.NewRouter(
		cfg.AuthServiceURL,
		cfg.UserServiceURL,
		cfg.CatalogServiceURL,
		cfg.PlaylistServiceURL,
		cfg.LibraryServiceURL,
		cfg.PlaybackServiceURL,
		cfg.MediaServiceURL,
	)
	if err != nil {
		return nil, err
	}
	rt.RegisterRoutes(r)

	return &Application{cfg: cfg, logger: log, router: r}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", a.cfg.HTTPPort), Handler: a.router}

	go func() {
		a.logger.Info("api-gateway listening", zap.Int("port", a.cfg.HTTPPort))
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
