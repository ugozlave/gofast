package gofast

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ugozlave/cargo"
)

type App struct {
	server    *http.Server
	config    Config[AppConfig]
	container *cargo.Container
	generator UniqueIDGenerator
}

func New(cfg Config[AppConfig]) *App {
	ctn := cargo.New()
	cargo.RegisterKV[Config[AppConfig]](ctn, func(cargo.BuilderContext) Config[AppConfig] { return cfg })
	gen := NewSequenceIDGenerator()
	return &App{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Value().Server.Host, cfg.Value().Server.Port),
			Handler: nil,
		},
		config:    cfg,
		container: ctn,
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

func (app *App) Run(ctx context.Context) {
	ctx = context.WithValue(ctx, CtxApplicationName, app.config.Value().Name)
	ctx = context.WithValue(ctx, CtxEnvironment, app.config.Value().Env)

	app.container.CreateScope(fmt.Sprintf(ScopeApplicationKeyFormat, app.config.Value().Name))

	server := app.server
	server.Handler = NewHttpInjector(app.generator, app.container, ctx).Handler()

	if SETTINGS.DEBUG {
		app.Inspect()
	}

	go func() {
		fmt.Println("server start")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()
	fmt.Println()
}

func (app *App) Shutdown() {
	container := app.container
	defer container.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server := app.server
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

type Settings struct {
	DEBUG                  bool
	CONFIG_FILE_NAME       string
	CONFIG_FILE_EXT        string
	CONFIG_APPLICATION_KEY string
	ENV_PREFIX             string
}

var SETTINGS = &Settings{
	DEBUG:                  true,
	CONFIG_FILE_NAME:       "config",
	CONFIG_FILE_EXT:        "json",
	CONFIG_APPLICATION_KEY: "Application",
	ENV_PREFIX:             "GOFAST",
}
