package gofast

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	LogApplication string = "application"
	LogEnvironment string = "environment"
	LogHostname    string = "hostname"
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
	DbgCtx(ctx context.Context, msg string, args ...any)
	Inf(msg string, args ...any)
	InfCtx(ctx context.Context, msg string, args ...any)
	Wrn(msg string, args ...any)
	WrnCtx(ctx context.Context, msg string, args ...any)
	Err(msg string, args ...any)
	ErrCtx(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
	WithGroup(name string) Logger
}

/*
** Logger
 */

type FastLogger struct {
	logger *slog.Logger
}

type LoggerBuilderOptions struct {
	Name     string
	Env      string
	Hostname string
}

func LoggerBuilder() Builder[*FastLogger] {
	return func(ctx *BuilderContext) *FastLogger {
		cfg := MustGetConfig[LoggerConfig](ctx, Singleton).Value()
		env := Environment.Get()
		hostname, _ := os.Hostname()
		name := ctx.Name()
		level := slog.LevelInfo
		switch strings.ToLower(cfg.Level) {
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
		handlerOpts := slog.HandlerOptions{
			Level: level,
		}
		var writer io.Writer
		switch cfg.Discard {
		case true:
			writer = io.Discard
		default:
			writer = os.Stdout
		}
		var handler slog.Handler
		switch cfg.Human {
		case true:
			handler = slog.NewTextHandler(writer, &handlerOpts)
		default:
			handler = slog.NewJSONHandler(writer, &handlerOpts)
		}
		logger := slog.New(handler)
		attrs := make([]any, 0, 3)
		if name != "" {
			attrs = append(attrs, slog.String(LogApplication, name))
		}
		if env != "" {
			attrs = append(attrs, slog.String(LogEnvironment, env))
		}
		if hostname != "" {
			attrs = append(attrs, slog.String(LogHostname, hostname))
		}
		return &FastLogger{logger: logger.With(attrs...)}
	}
}

func (l *FastLogger) Dbg(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *FastLogger) DbgCtx(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *FastLogger) Inf(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *FastLogger) InfCtx(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *FastLogger) Wrn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *FastLogger) WrnCtx(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *FastLogger) Err(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *FastLogger) ErrCtx(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *FastLogger) With(args ...any) Logger {
	return &FastLogger{logger: l.logger.With(args...)}
}

func (l *FastLogger) WithGroup(name string) Logger {
	return &FastLogger{logger: l.logger.WithGroup(name)}
}

/*
** LoggerConfig
 */

type LoggerConfig struct {
	Level   string `json:"Level"`
	Human   bool   `json:"Human"`
	Discard bool   `json:"Discard"`
}

func (c LoggerConfig) Path() []string {
	return []string{"Logging"}
}
