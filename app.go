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
	config    Config[AppConfig]
	container *cargo.Container
}

func New(cfg Config[AppConfig]) *App {
	ctn := cargo.New()
	Register[Config[AppConfig]](&App{container: ctn}, func(*BuilderContext) Config[AppConfig] { return cfg })
	Register[UniqueIDGenerator](&App{container: ctn}, func(*BuilderContext) UniqueIDGenerator { return &SequenceIDGenerator{} })
	return &App{
		config:    cfg,
		container: ctn,
	}
}

func (app *App) Run(ctx context.Context) {

	cfg := app.config.Value()

	ctx = context.WithValue(ctx, CtxApplicationName, cfg.Name)
	ctx = context.WithValue(ctx, CtxEnvironment, cfg.Env)

	ctn := app.container
	defer ctn.Close()
	ctn.CreateScope(fmt.Sprintf(ScopeApplicationKeyFormat, cfg.Name))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := http.Server{
		Addr:    addr,
		Handler: NewHttpInjector(ctn, ctx).Handler(),
	}

	if SETTINGS.DEBUG {
		app.Inspect()
	}

	fmt.Printf("server start [%v]\n", addr)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	fmt.Println()

	timeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(timeout); err != nil {
		fmt.Println("server shutdown failed:", err.Error())
	}

	fmt.Println("server stop")

}

func (app *App) Inspect() {
	ctn := app.container
	ctn.Inspect()
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
