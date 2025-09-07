package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type LogLevel string

type contextKey string

const TraceIDKey contextKey = "trace_id"

type Logger struct {
	*slog.Logger
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo, // default log level
		})),
	}
}

func GenerateTraceID() string {
	return uuid.New().String()
}

func WithTraceID(ctx context.Context) context.Context {
	if ctx.Value(TraceIDKey) == nil {
		return context.WithValue(ctx, TraceIDKey, GenerateTraceID())
	}
	return ctx
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	traceID := GetTraceID(ctx)
	if traceID != "" {
		args = append(args, slog.String(string(TraceIDKey), traceID))
	}
	l.Logger.Info(msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	traceID := GetTraceID(ctx)
	if traceID != "" {
		args = append(args, slog.String(string(TraceIDKey), traceID))
	}
	l.Logger.Error(msg, args...)
}

func (l *Logger) ErrorWithoutContext(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}
