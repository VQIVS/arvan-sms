package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sms/config"
	"sms/internal/api/handlers/messaging"
	"sms/internal/app"
	"sms/pkg/logger"
	"syscall"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	// Load configuration
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	appLogger := logger.NewLogger(logger.LogLevel("info"))

	application, err := app.NewApp(cfg)
	if err != nil {
		appLogger.ErrorWithoutContext("Failed to initialize app", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logger.WithTraceID(ctx)

	smsService := application.SMSService(ctx)

	consumer := messaging.NewSMSConsumer(*smsService, appLogger, application.RabbitConn(), application.Config())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		appLogger.Logger.Info("Starting SMS consumer worker")
		if err := consumer.Run(ctx); err != nil && err != context.Canceled {
			errChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		appLogger.Logger.Info("Received shutdown signal", "signal", sig)
		cancel()
	case err := <-errChan:
		appLogger.Error(ctx, "Consumer error", "error", err)
		cancel()
	}

	appLogger.Logger.Info("SMS consumer worker shutdown complete")
}
