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
	container *cargo.Container
	context   context.Context
	generator UniqueIDGenerator
}

func New() *App {
	c := cargo.New()
	c.Scopes.Create(ScopeApplicationKey)
	cargo.RegisterKV[Logger](c, NewFastLoggerWithDefaults)
	return &App{
		container: c,
		server: &http.Server{
			Addr:    ":8080",
			Handler: nil,
		},
		context:   context.Background(),
		generator: NewSequenceIDGenerator(),
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
	cargo.Inspect(app.container)

	server := app.server
	server.Handler = Handler(app)
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

func Handler(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create unique request ID
		id := app.generator.Next()

		// create unique scope for the request
		scope := fmt.Sprintf(ScopeRequestKeyFormat, id)
		app.container.Scopes.Create(scope)
		defer func() {
			app.container.Scopes.Delete(scope)
		}()

		// create a new builder context
		ctx := NewBuilderContext(context.WithValue(app.context, CtxRequestId, id), app.container)

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
