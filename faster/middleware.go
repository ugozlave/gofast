package faster

import (
	"net/http"

	"github.com/ugozlave/gofast"
)

/*
** LogMiddleware
 */

type LogMiddleware struct {
	logger gofast.Logger
}

func NewLogMiddleware(ctx *gofast.BuilderContext) *LogMiddleware {
	return &LogMiddleware{
		logger: gofast.MustGetTypedLogger[LogMiddleware](ctx, gofast.Scoped),
	}
}

func (m *LogMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &writer{ResponseWriter: w}
		group := m.logger.
			WithGroup("http").
			With(
				gofast.LogMethod, r.Method,
				gofast.LogHost, r.Host,
				gofast.LogUrl, r.URL.String(),
				gofast.LogRemote, r.RemoteAddr,
				gofast.LogAgent, r.UserAgent(),
			)
		group.Inf("request received")
		defer func() {
			group.Inf("request finished",
				gofast.LogStatus, writer.status,
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
