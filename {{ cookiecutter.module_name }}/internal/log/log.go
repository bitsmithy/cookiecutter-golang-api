package log

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"

	"{{ cookiecutter.module_path }}/internal/version"
)

func New() *slog.Logger {
	handler := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})
	return slog.New(handler)
}

func buildInfo() slog.Attr {
	return slog.Group("build", slog.String("version", version.Get()))
}

func baseLogger() *slog.Logger {
	return New().With(buildInfo())
}

func Debug(ctx context.Context, msg string, args ...any) {
	baseLogger().DebugContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	baseLogger().ErrorContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	baseLogger().InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	baseLogger().WarnContext(ctx, msg, args...)
}
