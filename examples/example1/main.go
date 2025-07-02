package main

import (
	fast "github.com/ugozlave/gofast"
)

func main() {
	app := fast.New()

	// controllers
	fast.Add(app, fast.NewHealthController)

	// middleware
	fast.Use(app, fast.NewLogMiddleware)

	app.Run()
}
