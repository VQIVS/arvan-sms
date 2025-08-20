package rabbit

import (
	"log/slog"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  2.0,
	}
}

type RetryStrategy struct {
	config RetryConfig
}

func NewRetryStrategy(config RetryConfig) *RetryStrategy {
	return &RetryStrategy{config: config}
}

func (rs *RetryStrategy) Execute(body []byte, handler func([]byte) error, logger *slog.Logger) error {
	var lastErr error

	for attempt := 1; attempt <= rs.config.MaxAttempts; attempt++ {
		if err := handler(body); err != nil {
			lastErr = err
			logger.Warn("handler execution failed",
				"attempt", attempt,
				"max_attempts", rs.config.MaxAttempts,
				"error", err)

			if attempt < rs.config.MaxAttempts {
				delay := rs.calculateDelay(attempt)
				logger.Debug("retrying after delay", "delay", delay)
				time.Sleep(delay)
				continue
			}
		} else {
			if attempt > 1 {
				logger.Info("handler executed successfully after retries", "attempts", attempt)
			} else {
				logger.Info("handler executed successfully")
			}
			return nil
		}
	}

	return lastErr
}

func (rs *RetryStrategy) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(rs.config.BaseDelay) *
		(rs.config.Multiplier * float64(attempt-1)))

	if delay > rs.config.MaxDelay {
		delay = rs.config.MaxDelay
	}

	return delay
}
