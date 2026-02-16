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

type FastLogger struct {
	logger *slog.Logger
}

type LoggerBuilderOptions struct {
	Name string
	Env  string
}

type LoggerBuilderOption func(*LoggerBuilderOptions)

func WithApplicationName(name string) LoggerBuilderOption {
	return func(opts *LoggerBuilderOptions) {
		opts.Name = name
	}
}

func WithEnvironment() LoggerBuilderOption {
	return func(opts *LoggerBuilderOptions) {
		opts.Env = Environment.Get()
	}
}

func LoggerBuilder(opts ...LoggerBuilderOption) Builder[*FastLogger] {
	cfg := NewConfig(LoggerConfig{}).Value()
	opt := &LoggerBuilderOptions{}
	for _, fn := range opts {
		fn(opt)
	}
	return func(ctx *gofast.BuilderContext) *FastLogger {
		return NewLogger(cfg.Level, opt.Name, opt.Env)
	}
}

func NewLogger(lvl, name, env string) *FastLogger {
	level := slog.LevelInfo
	switch strings.ToLower(lvl) {
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
	return &FastLogger{logger: logger}
}

func (l *FastLogger) Dbg(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *FastLogger) Inf(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *FastLogger) Wrn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *FastLogger) Err(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *FastLogger) With(args ...any) Logger {
	return &FastLogger{logger: l.logger.With(args...)}
}

func (l *FastLogger) WithGroup(name string) Logger {
	return &FastLogger{logger: l.logger.WithGroup(name)}
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

/*
** LoggerConfig
 */

type LoggerConfig struct {
	Level string `json:"Level"`
}

func (c LoggerConfig) Path() []string {
	return []string{"Logging"}
}
