package faster

import "github.com/ugozlave/gofast"

type NullLogger struct {
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

func (l *NullLogger) With(args ...any) gofast.Logger {
	return &NullLogger{}
}

func (l *NullLogger) WithGroup(name string) gofast.Logger {
	return &NullLogger{}
}
