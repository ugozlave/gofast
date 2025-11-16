package faster

import "github.com/ugozlave/gofast"

func New() *gofast.App {
	app := gofast.New(NewDefaultAppConfig())
	Log(app)
	Health(app)
	Recover(app)
	return app
}

func Log(app *gofast.App) {
	gofast.Log(app, NewFastLogger)
	gofast.Use(app, NewLogMiddleware)
}

func Health(app *gofast.App) {
	gofast.Add(app, NewHealthController)
}

func Recover(app *gofast.App) {
	gofast.Use(app, NewRecoverMiddleware)
}
