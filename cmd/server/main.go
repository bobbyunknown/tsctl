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
	"tsctl/internal/repository"
	"tsctl/internal/usecase"
	"tsctl/pkg/config"
	"tsctl/pkg/logger"

	_ "tsctl/docs"
)

// @title Tailscale Controller API
// @version 1.0
// @description HTTP API for controlling Tailscale daemon and operations
// @host localhost:8080
// @BasePath /
func clearLogFiles(paths ...string) {
	for _, path := range paths {
		if path == "" {
			continue
		}
		os.Remove(path)
	}
}

func main() {
	configPath := flag.String("c", "config/app.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	clearLogFiles(cfg.Logging.AppLogPath, cfg.Logging.DaemonLogPath)

	if err := logger.Init(cfg.Logging.AppLogPath, cfg.Logging.Level, cfg.Logging.Format); err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("tsctl starting")

	daemonRepo := repository.NewDaemonManager()
	tailscaleRepo := repository.NewTailscaleWrapper()

	daemonUC := usecase.NewDaemonUseCase(daemonRepo)
	tailscaleUC := usecase.NewTailscaleUseCase(tailscaleRepo)

	if cfg.Tailscale.AutoStart {
		logger.Log.Info("auto_start enabled, starting daemon")
		if err := daemonUC.Start(); err != nil {
			logger.Log.WithError(err).Error("failed to auto-start daemon")
		}
	}

	handler := httphandler.NewHandler(daemonUC, tailscaleUC)
	router := httphandler.SetupRouter(handler, cfg.Server.Mode)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Log.WithField("addr", addr).Info("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.WithError(err).Fatal("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.WithError(err).Error("server forced to shutdown")
	}

	status, _ := daemonUC.Status()
	if status != nil && status.Running {
		logger.Log.Info("stopping daemon")
		daemonUC.Stop()
	}

	logger.Log.Info("server stopped")
}
