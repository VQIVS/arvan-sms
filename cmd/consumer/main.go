package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sms-dispatcher/api/handler/consumer"
	"sms-dispatcher/app"
	"sms-dispatcher/config"
	"sms-dispatcher/pkg/logger"

	"github.com/google/uuid"
)

func main() {
	cfg := config.MustReadConfig("config.json")

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	h := consumer.New(a)
	logger.GetLogger().With("trace_id", uuid.NewString()).Info("consumer started")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	defer func() {
		if a.Rabbit() != nil {
			logger.GetLogger().With("trace_id", uuid.NewString()).Info("closing rabbit connection")
			a.Rabbit().Close()
		}
	}()

	if err := h.Start(ctx); err != nil {
		log.Printf("consumer stopped with error: %v", err)
		os.Exit(1)
	}
}
