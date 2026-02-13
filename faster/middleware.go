package faster

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ugozlave/gofast"
)

/*
** LogMiddleware
 */

type LogMiddleware struct {
	logger Logger
}

func LogMiddlewareBuilder() Builder[*LogMiddleware] {
	return func(ctx *gofast.BuilderContext) *LogMiddleware {
		return NewLogMiddleware(ctx)
	}
}

func NewLogMiddleware(ctx *gofast.BuilderContext) *LogMiddleware {
	return &LogMiddleware{
		logger: MustGetLogger[LogMiddleware](ctx, gofast.Scoped),
	}
}

func (m *LogMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
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
				LogDuration, time.Since(t),
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

/*
** RecoverMiddleware
 */

type RecoverMiddleware struct {
	logger Logger
}

func RecoverMiddlewareBuilder() Builder[*RecoverMiddleware] {
	return func(ctx *gofast.BuilderContext) *RecoverMiddleware {
		return NewRecoverMiddleware(ctx)
	}
}

func NewRecoverMiddleware(ctx *gofast.BuilderContext) *RecoverMiddleware {
	return &RecoverMiddleware{
		logger: MustGetLogger[RecoverMiddleware](ctx, gofast.Scoped),
	}
}

func (m *RecoverMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				msg := fmt.Sprintf("panic: %s", rec)
				if gofast.SETTINGS.DEBUG {
					msg += fmt.Sprintf("\n\n%s", debug.Stack())
				}
				m.logger.Err(msg)
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

/*
** TimeoutMiddleware
 */

type TimeoutMiddleware struct {
	Timeout time.Duration
}

func TimeoutMiddlewareBuilder() Builder[*TimeoutMiddleware] {
	return func(ctx *gofast.BuilderContext) *TimeoutMiddleware {
		return NewTimeoutMiddleware(ctx)
	}
}

func NewTimeoutMiddleware(ctx *gofast.BuilderContext) *TimeoutMiddleware {
	return &TimeoutMiddleware{
		Timeout: 30 * time.Second,
	}
}

func (m *TimeoutMiddleware) Handle(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, m.Timeout, "Timeout")
}
