package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httphandler "tsctl/internal/delivery/http"
	"tsctl/internal/delivery/http/websocket"
	"tsctl/internal/repository"
	"tsctl/internal/usecase"
	"tsctl/pkg/config"
	"tsctl/pkg/logger"

	_ "tsctl/docs"
)

// @title Tailscale Controller API
// @version 2.0
// @description Self-contained HTTP API for Tailscale with embedded daemon
// @host localhost:8080
// @BasePath /

func main() {
	configPath := flag.String("c", "config/app.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	os.Remove(cfg.Logging.AppLogPath)

	if err := logger.Init(cfg.Logging.AppLogPath, cfg.Logging.Level, cfg.Logging.Format); err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("tsctl starting (embedded mode)")

	tsnetMgr := repository.NewTsnetManager(cfg)

	ctx := context.Background()

	tailscaleUC := usecase.NewTailscaleUseCase(tsnetMgr)

	wsHub := websocket.NewWebSocketHub()
	go wsHub.Run()

	handler := httphandler.NewHandler(tailscaleUC, wsHub)
	router := httphandler.SetupRouter(handler, cfg.Server.Mode)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Log.WithField("addr", addr).Info("HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.WithError(err).Fatal("server failed")
		}
	}()

	go func() {
		if err := tsnetMgr.Start(ctx); err != nil {
			logger.Log.WithError(err).Error("failed to start embedded daemon")
		}
	}()

	watcher := websocket.NewStatusWatcher(wsHub, tsnetMgr)
	go watcher.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.WithError(err).Error("server forced to shutdown")
	}

	logger.Log.Info("server stopped")
}
