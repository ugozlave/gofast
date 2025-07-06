package gofast

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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

func New(cfg ConfigProvider[AppConfig]) *App {
	ctx := context.WithValue(context.Background(), CtxEnvironment, cfg.Value().Env)
	ctn := cargo.New()
	ctn.CreateScope(ScopeApplicationKey)
	cargo.RegisterKV[ConfigProvider[AppConfig]](ctn, func(cargo.BuilderContext) ConfigProvider[AppConfig] { return cfg })
	gen := NewSequenceIDGenerator()
	return &App{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Value().Server.Host, cfg.Value().Server.Port),
			Handler: nil,
		},
		config:    cfg,
		container: ctn,
		context:   ctx,
		generator: gen,
	}
}

func (app *App) WithIDGenerator(generator UniqueIDGenerator) *App {
	if generator == nil {
		panic("ID generator cannot be nil")
	}
	app.generator = generator
	return app
}

func (app *App) Run() {
	if DEBUG {
		app.Inspect()
	}

	server := app.server
	server.Handler = NewHttpInjector(app.generator, app.container, app.context).Handler()
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
	fmt.Println()
	fmt.Println("Config:")
	fmt.Printf(".   %v\n", app.config.Value())
	fmt.Println()
}

var DEBUG bool = true
