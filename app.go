package gofast

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/ugozlave/cargo"
)

type App struct {
	server    *http.Server
	config    ConfigProvider[AppConfig]
	container *cargo.Container
	context   context.Context
	generator UniqueIDGenerator
}

func New() *App {
	ctx := context.Background()
	cfg := NewAppConfig()
	ctn := cargo.New()
	ctn.Scopes.Create(ScopeApplicationKey)
	cargo.RegisterKV[ConfigProvider[AppConfig]](ctn, func(cargo.BuilderContext) *Config[AppConfig] { return cfg })
	cargo.RegisterKV[Logger](ctn, func(c cargo.BuilderContext) *FastLogger { return NewFastLoggerWithDefaults(NewBuilderContext(c, ctn)) })
	gen := NewSequenceIDGenerator()
	return &App{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Value().Server.Host, cfg.Value().Server.Port),
			Handler: Handler(gen, ctn, ctx),
		},
		config:    cfg,
		container: ctn,
		context:   ctx,
		generator: gen,
	}
}

func (app *App) WithContext(ctx context.Context) *App {
	if ctx == nil {
		panic("context cannot be nil")
	}
	app.context = ctx
	return app
}

func (app *App) WithIDGenerator(generator UniqueIDGenerator) *App {
	if generator == nil {
		panic("ID generator cannot be nil")
	}
	app.generator = generator
	return app
}

func (app *App) Run() {
	app.Inspect()

	server := app.server
	go func() {
		fmt.Println("server start")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("server shutdown failed:", err.Error())
	}
	fmt.Println("server stop")
}

func (app *App) Inspect() {
	cargo.Inspect(app.container)
	//fmt.Println(app.config.Value())
}

func Handler(gen UniqueIDGenerator, ctn *cargo.Container, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create unique request ID
		id := gen.Next()

		// create unique scope for the request
		scope := fmt.Sprintf(ScopeRequestKeyFormat, id)
		ctn.Scopes.Create(scope)
		defer func() {
			ctn.Scopes.Delete(scope)
		}()

		// create a new builder context
		ctx := NewBuilderContext(context.WithValue(ctx, CtxRequestId, id), ctn)

		// build the controllers
		mux := controllers(ctx, scope)

		// apply middlewares
		mux = middlewares(ctx, scope, mux)

		mux.ServeHTTP(w, r)
	})
}

func controllers(ctx cargo.BuilderContext, scope string) http.Handler {
	var mux *http.ServeMux = http.NewServeMux()
	container := ctx.C()
	for t := range container.Services {
		if t.Implements(reflect.TypeOf((*Controller)(nil)).Elem()) {
			ctrl := container.Get(t, scope, ctx).(Controller)
			mux.Handle(ctrl.Prefix()+"/", http.StripPrefix(ctrl.Prefix(), ctrl.Routes()))
		}
	}
	return mux
}

func middlewares(ctx cargo.BuilderContext, scope string, ctrl http.Handler) http.Handler {
	container := ctx.C()
	for t := range container.Services {
		if t.Implements(reflect.TypeOf((*Middleware)(nil)).Elem()) {
			mw := container.Get(t, scope, ctx).(Middleware)
			ctrl = mw.Handle(ctrl)
		}
	}
	return ctrl
}
