package gofast

import (
	"log/slog"
	"os"
)

const (
	LogApplication string = "application"
	LogEnvironment string = "environment"
	LogService     string = "service"
	LogRequestId   string = "requestId"
	LogHttp        string = "http"
	LogMethod      string = "method"
	LogHost        string = "host"
	LogUrl         string = "url"
	LogRemote      string = "remote"
	LogAgent       string = "agent"
	LogStatus      string = "status"
)

type Logger interface {
	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)
	Wrn(msg string, args ...any)
	Err(msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}

type FastLogger struct {
	*slog.Logger
}

func NewFastLogger(ctx *BuilderContext) *FastLogger {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)
	logger := slog.New(handler)
	return &FastLogger{
		Logger: logger,
	}
}

func (l *FastLogger) Dbg(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *FastLogger) Inf(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *FastLogger) Wrn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *FastLogger) Err(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *FastLogger) With(args ...any) Logger {
	return &FastLogger{
		Logger: l.Logger.With(args...),
	}
}

func (l *FastLogger) WithGroup(name string) Logger {
	return &FastLogger{
		Logger: l.Logger.WithGroup(name),
	}
}

func NewFastLoggerWithDefaults(ctx *BuilderContext) *FastLogger {
	config := Get[ConfigProvider[AppConfig]](ctx, Singleton)
	application := config.Value().App.Name
	if application == "" {
		application, _ = os.Executable()
	}
	return NewFastLogger(ctx).
		WithApplication(application).
		WithEnvironment(config.Value().Env)
}

func (l *FastLogger) WithApplication(v string) *FastLogger {
	return &FastLogger{
		Logger: l.Logger.With(slog.String("application", v)),
	}
}

func (l *FastLogger) WithEnvironment(v string) *FastLogger {
	return &FastLogger{
		Logger: l.Logger.With(slog.String("environment", v)),
	}
}
