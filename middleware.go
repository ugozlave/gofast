package gofast

import (
	"net/http"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

/*
** LogMiddleware
 */

type LogMiddleware struct {
	logger Logger
}

func NewLogMiddleware(ctx *BuilderContext) *LogMiddleware {
	return &LogMiddleware{
		logger: TypedLogger[LogMiddleware](ctx, Scoped),
	}
}

func (m *LogMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &writer{ResponseWriter: w}
		group := m.logger.
			WithGroup("http").
			With(
				LogMethod, r.Method,
				LogHost, r.Host,
				LogUrl, r.URL.String(),
				LogRemote, r.RemoteAddr,
				LogAgent, r.UserAgent(),
			)
		group.Inf("request received")
		defer func() {
			group.Inf("request finished",
				LogStatus, writer.status,
			)
		}()
		next.ServeHTTP(writer, r)
	})
}

type writer struct {
	http.ResponseWriter
	status int
}

func (w *writer) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
