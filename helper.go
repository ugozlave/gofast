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
	var v T
	switch lt {
	case Singleton:
		v = cargo.MustGet[T](ctx.C(), ScopeApplicationKey, ctx)
	case Scoped:
		scope := ctx.RequestId()
		v = cargo.MustGet[T](ctx.C(), fmt.Sprintf(ScopeRequestKeyFormat, scope), ctx)
	case Transient:
		v = cargo.MustBuild[T](ctx.C(), ctx)
	default:
		v = cargo.MustGet[T](ctx.C(), ScopeApplicationKey, ctx)
	}
	return v
}

func MustGet[T any](ctx *BuilderContext, lt Lifetime) T {
	var v T
	switch lt {
	case Singleton:
		v = cargo.MustGet[T](ctx.C(), ScopeApplicationKey, ctx)
	case Scoped:
		scope := ctx.RequestId()
		v = cargo.MustGet[T](ctx.C(), fmt.Sprintf(ScopeRequestKeyFormat, scope), ctx)
	case Transient:
		v = cargo.MustBuild[T](ctx.C(), ctx)
	default:
		v = cargo.MustGet[T](ctx.C(), ScopeApplicationKey, ctx)
	}
	if any(v) == nil {
		panic(fmt.Sprintf("type %T is nil", new(T)))
	}
	return v
}

func Add[C Controller](app *App, builder func(ctx *BuilderContext) C) {
	Register[Controller](app, builder)
}

func Use[M Middleware](app *App, builder func(ctx *BuilderContext) M) {
	Register[Middleware](app, builder)
}

func Log[L Logger](app *App, builder func(ctx *BuilderContext) L) {
	Register[Logger](app, builder)
}

func Cfg[C any](app *App, builder func(ctx *BuilderContext) *Config[C]) {
	Register[ConfigProvider[C]](app, builder)
}

func GetLogger[S any](ctx *BuilderContext, lt Lifetime) Logger {
	logger := Get[Logger](ctx, lt).With(LogService, cargo.From[S]())
	switch lt {
	case Scoped:
		return logger.With(LogRequestId, ctx.RequestId())
	default:
		return logger
	}
}

func MustGetLogger[S any](ctx *BuilderContext, lt Lifetime) Logger {
	logger := MustGet[Logger](ctx, lt).With(LogService, cargo.From[S]())
	switch lt {
	case Scoped:
		return logger.With(LogRequestId, ctx.RequestId())
	default:
		return logger
	}
}

func GetConfig[C any](ctx *BuilderContext, lt Lifetime) ConfigProvider[C] {
	return Get[ConfigProvider[C]](ctx, lt)
}

func MustGetConfig[C any](ctx *BuilderContext, lt Lifetime) ConfigProvider[C] {
	return MustGet[ConfigProvider[C]](ctx, lt)
}
