package gofast

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

/*
** BodyLimiterMiddleware
 */

type BodyLimiterMiddleware struct {
	Limit int64
}

func BodyLimiterMiddlewareBuilder() Builder[*BodyLimiterMiddleware] {
	return func(*BuilderContext) *BodyLimiterMiddleware {
		return &BodyLimiterMiddleware{
			Limit: 1 << 20, // 1MB
		}
	}
}

func (m *BodyLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, m.Limit)
		next.ServeHTTP(w, r)
	})
}

/*
** LogMiddleware
 */

type LogMiddleware struct {
	logger Logger
}

func LogMiddlewareBuilder() Builder[*LogMiddleware] {
	return func(ctx *BuilderContext) *LogMiddleware {
		return &LogMiddleware{
			logger: MustGetLogger[LogMiddleware](ctx, Scoped),
		}
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
		group.Dbg("request received")
		defer func() {
			group.Dbg("request finished",
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
	return func(ctx *BuilderContext) *RecoverMiddleware {
		return &RecoverMiddleware{
			logger: MustGetLogger[RecoverMiddleware](ctx, Scoped),
		}
	}
}

func (m *RecoverMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				msg := fmt.Sprintf("panic: %s", rec)
				if SETTINGS.DEBUG {
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
	return func(ctx *BuilderContext) *TimeoutMiddleware {
		return &TimeoutMiddleware{
			Timeout: 30 * time.Second,
		}
	}
}

func (m *TimeoutMiddleware) Handle(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, m.Timeout, "Timeout")
}
