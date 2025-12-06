package faster

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ugozlave/gofast"
)

type FastLogger struct {
	*slog.Logger
}

func NewFastLogger(ctx *gofast.BuilderContext) *FastLogger {
	config := gofast.MustGetConfig[gofast.AppConfig](ctx, gofast.Singleton)
	application := config.Value().Name
	if application == "" {
		application, _ = os.Executable()
	}
	level := slog.LevelInfo
	switch strings.ToLower(config.Value().Log.Level) {
	case "debug", "dbg", "d":
		level = slog.LevelDebug
	case "info", "inf", "i":
		level = slog.LevelInfo
	case "warning", "warn", "wrn", "w":
		level = slog.LevelWarn
	case "error", "err", "e":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)
	logger := slog.New(handler)
	return &FastLogger{
		Logger: logger.
			With(slog.String(gofast.LogApplication, application)).
			With(slog.String(gofast.LogEnvironment, config.Value().Env)),
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

func (l *FastLogger) With(args ...any) gofast.Logger {
	return &FastLogger{
		Logger: l.Logger.With(args...),
	}
}

func (l *FastLogger) WithGroup(name string) gofast.Logger {
	return &FastLogger{
		Logger: l.Logger.WithGroup(name),
	}
}

func (l *FastLogger) WithApplication(v string) *FastLogger {
	return &FastLogger{
		Logger: l.Logger.With(slog.String(gofast.LogApplication, v)),
	}
}

func (l *FastLogger) WithEnvironment(v string) *FastLogger {
	return &FastLogger{
		Logger: l.Logger.With(slog.String(gofast.LogEnvironment, v)),
	}
}

func (l *FastLogger) WithService(v string) *FastLogger {
	return &FastLogger{
		Logger: l.Logger.With(slog.String(gofast.LogService, v)),
	}
}

/*
** NullLogger
 */

type NullLogger struct {
}

func NewNullLogger(_ *gofast.BuilderContext) *NullLogger {
	return &NullLogger{}
}

func (l *NullLogger) Dbg(msg string, args ...any) {
}

func (l *NullLogger) Inf(msg string, args ...any) {
}

func (l *NullLogger) Wrn(msg string, args ...any) {
}

func (l *NullLogger) Err(msg string, args ...any) {
}

func (l *NullLogger) With(args ...any) gofast.Logger {
	return &NullLogger{}
}

func (l *NullLogger) WithGroup(name string) gofast.Logger {
	return &NullLogger{}
}
