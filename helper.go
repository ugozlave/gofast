package gofast

import (
	"fmt"

	"github.com/ugozlave/cargo"
)

type Lifetime int

const (
	Singleton Lifetime = iota
	Scoped
	Transient
)

const (
	ScopeApplicationKey   = "gofast-scope-application"
	ScopeRequestKeyFormat = "gofast-scope-request-%s"
)

func Register[K any, V any](app *App, builder func(*BuilderContext) V) {
	cargo.RegisterKV[K](app.container, func(ctx cargo.BuilderContext) V {
		return builder(NewBuilderContext(ctx, app.container))
	})
}

func Get[T any](ctx *BuilderContext, lt Lifetime) T {
	switch lt {
	case Singleton:
		return cargo.MustGet[T](ctx.C(), ScopeApplicationKey, ctx)
	case Scoped:
		scope := ctx.RequestId()
		return cargo.MustGet[T](ctx.C(), fmt.Sprintf(ScopeRequestKeyFormat, scope), ctx)
	case Transient:
		return cargo.Build[T](ctx.C(), ctx)
	default:
		return cargo.MustGet[T](ctx.C(), ScopeApplicationKey, ctx)
	}
}

func Add[C Controller](app *App, builder func(ctx *BuilderContext) C) {
	Register[C](app, builder)
}

func Use[M Middleware](app *App, builder func(ctx *BuilderContext) M) {
	Register[M](app, builder)
}

func TypedLogger[S any](ctx *BuilderContext, lt Lifetime) Logger {
	logger := Get[Logger](ctx, Singleton).With(LogService, cargo.From[S]())
	switch lt {
	case Scoped:
		return logger.With(LogRequestId, ctx.RequestId())
	default:
		return logger
	}
}
