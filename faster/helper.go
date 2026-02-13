package faster

import (
	"github.com/ugozlave/gofast"
)

func New() *gofast.App {
	app := gofast.New(NewAppConfig())
	Add(app, HealthControllerBuilder())
	Use(app, LogMiddlewareBuilder())
	Use(app, RecoverMiddlewareBuilder())
	Use(app, TimeoutMiddlewareBuilder())
	Register[Logger](app, StdLoggerBuilder())
	Register[Cache](app, MemoryCacheBuilder())
	return app
}

func Register[K any, V any](app *gofast.App, builder func(*gofast.BuilderContext) V) {
	gofast.Register[K](app, builder)
}

func Add[C gofast.Controller](app *gofast.App, builder func(*gofast.BuilderContext) C) {
	gofast.Add(app, builder)
}

func Use[M gofast.Middleware](app *gofast.App, builder func(*gofast.BuilderContext) M) {
	gofast.Use(app, builder)
}

func Cfg[C ConfigProvider[T], T any](app *gofast.App, builder func(*gofast.BuilderContext) C) {
	gofast.Register[ConfigProvider[T]](app, builder)
}

func Get[T any](ctx *gofast.BuilderContext, lt gofast.Lifetime) T {
	return gofast.Get[T](ctx, lt)
}

func MustGet[T any](ctx *gofast.BuilderContext, lt gofast.Lifetime) T {
	return gofast.MustGet[T](ctx, lt)
}

func GetLogger[S any](ctx *gofast.BuilderContext, lt gofast.Lifetime) Logger {
	logger := Get[Logger](ctx, lt).With(LogService, gofast.From[S]())
	switch lt {
	case gofast.Scoped:
		return logger.With(LogRequestId, ctx.RequestID())
	default:
		return logger
	}
}

func MustGetLogger[S any](ctx *gofast.BuilderContext, lt gofast.Lifetime) Logger {
	logger := MustGet[Logger](ctx, lt).With(LogService, gofast.From[S]())
	switch lt {
	case gofast.Scoped:
		return logger.With(LogRequestId, ctx.RequestID())
	default:
		return logger
	}
}

func GetConfig[C any](ctx *gofast.BuilderContext, lt gofast.Lifetime) ConfigProvider[C] {
	return Get[ConfigProvider[C]](ctx, lt)
}

func MustGetConfig[C any](ctx *gofast.BuilderContext, lt gofast.Lifetime) ConfigProvider[C] {
	return MustGet[ConfigProvider[C]](ctx, lt)
}

func All[T any](ctx *gofast.BuilderContext, lt gofast.Lifetime) []T {
	return gofast.All[T](ctx, lt)
}

type Builder[T any] func(*gofast.BuilderContext) T
