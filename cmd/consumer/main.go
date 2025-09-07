package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sms/config"
	"sms/internal/api/handlers/messaging"
	"sms/internal/app"
	"sms/pkg/logger"
	"syscall"
)

var configPath = flag.String("config", "config.yaml", "service configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}

	c := config.MustReadConfig(*configPath)
	appLogger := logger.NewLogger(logger.LogLevel("info"))

	appContainer := app.NewMustApp(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logger.WithTraceID(ctx)

	smsService := appContainer.SMSService(ctx)
	consumer := messaging.NewSMSConsumer(*smsService, appLogger, appContainer.RabbitConn(), appContainer.Config())

	// Graceful shutdown handling
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
		appLogger.Info(ctx, "Consumer error", "error", err)
		cancel()
	}

	appLogger.Logger.Info("SMS consumer worker shutdown complete")
}
