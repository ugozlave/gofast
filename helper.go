package gofast

import (
	"fmt"
	"reflect"

	"github.com/ugozlave/cargo"
)

type Lifetime int

const (
	Singleton Lifetime = iota
	Scoped
	Transient
)

const (
	ScopeApplicationKeyFormat = "gofast-scope-application-%s"
	ScopeRequestKeyFormat     = "gofast-scope-request-%s"
)

func Register[K any, V any](app *App, builder func(*BuilderContext) V) {
	ctn := app.container
	key := From[K]()
	value := From[V]()
	if !value.AssignableTo(key) {
		panic(fmt.Sprintf("type %v is not assignable to %v", value, key))
	}
	ctn.Register(key.String(), value.String(), func(ctx cargo.BuilderContext) any {
		return builder(NewBuilderContext(ctx, ctn))
	})
}

func Add[C Controller](app *App, builder func(*BuilderContext) C) {
	Register[Controller](app, builder)
}

func Use[M Middleware](app *App, builder func(*BuilderContext) M) {
	Register[Middleware](app, builder)
}

func Log[L Logger](app *App, builder func(*BuilderContext) L) {
	Register[Logger](app, builder)
}

func Cfg[C Config[T], T any](app *App, builder func(*BuilderContext) C) {
	Register[Config[T]](app, builder)
}

func Get[T any](ctx *BuilderContext, lt Lifetime) T {
	ctn := ctx.container
	key := From[T]()
	var v T
	switch lt {
	case Singleton:
		name := ctx.ApplicationName()
		v = ctn.MustGet(key.String(), fmt.Sprintf(ScopeApplicationKeyFormat, name), ctx).(T)
	case Scoped:
		scope := ctx.RequestId()
		v = ctn.MustGet(key.String(), fmt.Sprintf(ScopeRequestKeyFormat, scope), ctx).(T)
	case Transient:
		v = ctn.MustBuild(key.String(), ctx).(T)
	}
	return v
}

func MustGet[T any](ctx *BuilderContext, lt Lifetime) T {
	ctn := ctx.container
	key := From[T]()
	var v T
	switch lt {
	case Singleton:
		name := ctx.ApplicationName()
		v = ctn.MustGet(key.String(), fmt.Sprintf(ScopeApplicationKeyFormat, name), ctx).(T)
	case Scoped:
		scope := ctx.RequestId()
		v = ctn.MustGet(key.String(), fmt.Sprintf(ScopeRequestKeyFormat, scope), ctx).(T)
	case Transient:
		v = ctn.MustBuild(key.String(), ctx).(T)
	}
	if any(v) == nil {
		panic(fmt.Sprintf("type %T is nil", new(T)))
	}
	return v
}

func GetLogger[S any](ctx *BuilderContext, lt Lifetime) Logger {
	logger := Get[Logger](ctx, lt).With(LogService, From[S]())
	switch lt {
	case Scoped:
		return logger.With(LogRequestId, ctx.RequestId())
	default:
		return logger
	}
}

func MustGetLogger[S any](ctx *BuilderContext, lt Lifetime) Logger {
	logger := MustGet[Logger](ctx, lt).With(LogService, From[S]())
	switch lt {
	case Scoped:
		return logger.With(LogRequestId, ctx.RequestId())
	default:
		return logger
	}
}

func GetConfig[C any](ctx *BuilderContext, lt Lifetime) Config[C] {
	return Get[Config[C]](ctx, lt)
}

func MustGetConfig[C any](ctx *BuilderContext, lt Lifetime) Config[C] {
	return MustGet[Config[C]](ctx, lt)
}

func All[T any](ctx *BuilderContext, lt Lifetime) []T {
	ctn := ctx.container
	key := From[T]()
	var instances []any
	switch lt {
	case Singleton:
		name := ctx.ApplicationName()
		instances = ctn.Gets(key.String(), fmt.Sprintf(ScopeApplicationKeyFormat, name), ctx)
	case Scoped:
		scope := ctx.RequestId()
		instances = ctn.Gets(key.String(), fmt.Sprintf(ScopeRequestKeyFormat, scope), ctx)
	case Transient:
		instances = ctn.Builds(key.String(), ctx)
	}
	result := make([]T, 0, len(instances))
	for _, instance := range instances {
		result = append(result, instance.(T))
	}
	return result
}

func From[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
