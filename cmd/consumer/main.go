package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sms-dispatcher/api/handler/consumer"
	"sms-dispatcher/app"
	"sms-dispatcher/config"
	"sms-dispatcher/pkg/logger"
)

var configPath = flag.String("config", "config.json", "service configuration file")

func main() {
	flag.Parse()
	logger := logger.GetTracedLogger()
	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}
	cfg := config.MustReadConfig(*configPath)

	// Create application
	a, err := app.NewApp(cfg)
	if err != nil {
		logger.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	defer cleanup(a, logger)
	h := consumer.NewHandler(a)

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	logger.Info("starting SMS consumer")

	// Start consumer in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- h.Start(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			logger.Error("consumer stopped with error", "error", err)
			os.Exit(1)
		}
		logger.Info("consumer completed successfully")

	case <-ctx.Done():
		logger.Info("shutdown signal received, stopping consumer")

		if err := h.Stop(); err != nil {
			logger.Error("failed to stop consumer gracefully", "error", err)
		}

		// Wait for consumer to finish with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		select {
		case err := <-done:
			if err != nil && err != context.Canceled {
				logger.Error("consumer stopped with error during shutdown", "error", err)
			} else {
				logger.Info("consumer stopped gracefully")
			}
		case <-shutdownCtx.Done():
			logger.Warn("consumer shutdown timed out")
		}
	}
}

func cleanup(a app.App, logger interface {
	Info(string, ...any)
	Error(string, ...any)
}) {
	logger.Info("cleaning up resources")

	if a != nil && a.Rabbit() != nil {
		logger.Info("closing rabbit connection")
		if err := a.Rabbit().Close(); err != nil {
			logger.Error("failed to close rabbit connection", "error", err)
		} else {
			logger.Info("rabbit connection closed successfully")
		}
	}

	logger.Info("cleanup completed")
}
