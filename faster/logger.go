package faster

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ugozlave/gofast"
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
	LogDuration    string = "duration"
)

type Logger interface {
	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)
	Wrn(msg string, args ...any)
	Err(msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}

/*
** Logger
 */

type StdLogger struct {
	logger *slog.Logger
}

func StdLoggerBuilder() Builder[*StdLogger] {
	return func(ctx *gofast.BuilderContext) *StdLogger {
		return NewStdLogger(ctx)
	}
}

func NewStdLogger(ctx *gofast.BuilderContext) *StdLogger {
	name, env := ctx.Name(), ctx.Environment()
	level := slog.LevelInfo
	switch strings.ToLower("debug") {
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
	if name != "" {
		logger = logger.With(slog.String(LogApplication, name))
	}
	if env != "" {
		logger = logger.With(slog.String(LogEnvironment, env))
	}
	return &StdLogger{logger: logger}
}

func (l *StdLogger) Dbg(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *StdLogger) Inf(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *StdLogger) Wrn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *StdLogger) Err(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *StdLogger) With(args ...any) Logger {
	return &StdLogger{l.logger.With(args...)}
}

func (l *StdLogger) WithGroup(name string) Logger {
	return &StdLogger{logger: l.logger.WithGroup(name)}
}

/*
** NullLogger
 */

type NullLogger struct {
}

func NullLoggerBuilder() Builder[*NullLogger] {
	return func(ctx *gofast.BuilderContext) *NullLogger {
		return NewNullLogger(ctx)
	}
}

func NewNullLogger(ctx *gofast.BuilderContext) *NullLogger {
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

func (l *NullLogger) With(args ...any) Logger {
	return &NullLogger{}
}

func (l *NullLogger) WithGroup(name string) Logger {
	return &NullLogger{}
}
