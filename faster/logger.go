package faster

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ugozlave/gofast"
)

/*
** Logger
 */

type Logger struct {
	*slog.Logger
}

func NewLogger(ctx *gofast.BuilderContext) *Logger {
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
	return &Logger{
		Logger: logger.
			With(slog.String(gofast.LogApplication, application)).
			With(slog.String(gofast.LogEnvironment, config.Value().Env)),
	}
}

func (l *Logger) Dbg(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *Logger) Inf(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *Logger) Wrn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *Logger) Err(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *Logger) With(args ...any) gofast.Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

func (l *Logger) WithGroup(name string) gofast.Logger {
	return &Logger{
		Logger: l.Logger.WithGroup(name),
	}
}

func (l *Logger) WithApplication(v string) *Logger {
	return &Logger{
		Logger: l.Logger.With(slog.String(gofast.LogApplication, v)),
	}
}

func (l *Logger) WithEnvironment(v string) *Logger {
	return &Logger{
		Logger: l.Logger.With(slog.String(gofast.LogEnvironment, v)),
	}
}

func (l *Logger) WithService(v string) *Logger {
	return &Logger{
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
