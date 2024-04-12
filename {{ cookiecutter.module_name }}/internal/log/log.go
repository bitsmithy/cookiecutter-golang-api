package log

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"{{ cookiecutter.module_path }}/internal/version"

	stdLog "log"
)

const (
	DebugLevel log.Level = log.DebugLevel
	ErrorLevel log.Level = log.ErrorLevel
	FatalLevel log.Level = log.FatalLevel
	InfoLevel  log.Level = log.InfoLevel
	WarnLevel  log.Level = log.WarnLevel
)

type Logger struct {
	logger *log.Logger
}

func New() *Logger {
	return &Logger{
		logger: log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.DateTime,
		}).With("build.version", version.Get()),
	}
}

func (l *Logger) StdLogger(level log.Level) *stdLog.Logger {
	return l.logger.StandardLog(log.StandardLogOptions{
		ForceLevel: level,
	})
}

func (l *Logger) WithSpan(span trace.Span) *Logger {
	l.logger = l.logger.
		With("telemetry.traceID", span.SpanContext().TraceID().String()).
		With("telemetry.spanID", span.SpanContext().SpanID().String())

	return l
}

func WithSpan(span trace.Span) *Logger {
	return New().WithSpan(span)
}

func (l *Logger) WithSpanAttrs(span trace.Span, attrs ...attribute.KeyValue) *Logger {
	span.SetAttributes(attrs...)

	keyvals := []any{}
	for _, kv := range attrs {
		keyvals = append(keyvals, kv.Key, kv.Value.Emit())
	}

	return l.WithSpan(span).With(keyvals...)
}

func WithSpanAttrs(span trace.Span, attrs ...attribute.KeyValue) *Logger {
	return New().WithSpanAttrs(span, attrs...)
}

func (l *Logger) With(keyvals ...any) *Logger {
	l.logger = l.logger.With(keyvals...)
	return l
}

func With(keyvals ...any) *Logger {
	return New().With(keyvals...)
}

func (l *Logger) Debug(msg any, keyvals ...any) {
	l.logger.Helper()
	l.logger.Debug(msg, keyvals...)
}

func (l *Logger) Error(msg any, keyvals ...any) {
	l.logger.Helper()
	l.logger.Error(msg, keyvals...)
}

func (l *Logger) Info(msg any, keyvals ...any) {
	l.logger.Helper()
	l.logger.Info(msg, keyvals...)
}

func (l *Logger) Warn(msg any, keyvals ...any) {
	l.logger.Helper()
	l.logger.Warn(msg, keyvals...)
}

func Debug(msg any, keyvals ...any) {
	New().logger.Debug(msg, keyvals...)
}

func Error(msg any, keyvals ...any) {
	New().logger.Error(msg, keyvals...)
}

func Info(msg any, keyvals ...any) {
	New().logger.Info(msg, keyvals...)
}

func Warn(msg any, keyvals ...any) {
	New().logger.Warn(msg, keyvals...)
}
