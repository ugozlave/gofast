# gofast

**gofast** is a minimalistic and high-performance framework for building web APIs in Go with ease and speed.

Designed to be clean, fast, and dependency-free, it provides a built-in dependency injection system and sensible defaults so you can focus on writing your application logic without boilerplate or third-party overhead.

## Features

- ⚡ **Zero external dependencies**
- 🧩 **Built-in dependency injection**
  - Application-wide singletons
  - Request-scoped services
  - Transient lifecycle
- 🔧 **Config system**
- 📄 **Structured logger**
- 📑 **Request logging middleware**
- ❤️ **Healthcheck endpoint**
- 🧼 **Clean and intuitive API**

### Prerequisites

**gofast** requires [Go](https://go.dev/) version [1.24.3](https://go.dev/doc/devel/release#go1.24.3) or above.

## Installation

Use Go modules to install **gofast** in your application.

```shell
go get github.com/ugozlave/gofast
```

## Getting Started

```go
package main

import fast "github.com/ugozlave/gofast"

func main() {
	app := fast.New()
	app.Run()
}
```

This minimal example starts a default web API application using the built-in HTTP server.

By default, it runs an http.Server on port 8080, with preconfigured services such as config and logging. You can extend it by registering your own controller, services, and configuration.

```go
package main

import fast "github.com/ugozlave/gofast"

func main() {
	app := fast.New()

    // controllers
	fast.Add(app, fast.NewHealthController)

    // middleware
	fast.Use(app, fast.NewLogMiddleware)

	app.Run()
}
```

```shell
map[*gofast.LogMiddleware:0x400006a188 gofast.ConfigProvider[github.com/ugozlave/gofast.AppConfig]:0x400006a168 *gofast.HealthController:0x400006a180 gofast.Logger:0x400006a170]
server start
time=2025-07-02T16:42:04.673Z level=INFO msg="request received" application=example environment=development service=gofast.LogMiddleware requestId=1 http.method=GET http.host=localhost:8080 http.url=/health/ http.remote=127.0.0.1:49672 http.agent=curl/8.7.1
time=2025-07-02T16:42:04.673Z level=INFO msg="request finished" application=example environment=development service=gofast.LogMiddleware requestId=1 http.method=GET http.host=localhost:8080 http.url=/health/ http.remote=127.0.0.1:49672 http.agent=curl/8.7.1 http.status=200
^C
server stop
```

This example explicitly enables the built-in health check controller and the request logging middleware. These components are optional and can be added or removed based on your application needs.

```go
package main

import (
	"net/http"

	fast "github.com/ugozlave/gofast"
)

func main() {
	app := fast.New()

	// controllers
	fast.Add(app, fast.NewHealthController)
	fast.Add(app, NewMyController)

	// middleware
	fast.Use(app, fast.NewLogMiddleware)

	// services
	fast.Register[Service](app, NewMyService)

	app.Run()
}

type MyController struct {
	service Service
}

func NewMyController(ctx *fast.BuilderContext) *MyController {
	return &MyController{
		service: fast.Get[Service](ctx, fast.Scoped),
	}
}

func (c *MyController) Prefix() string {
	return "/my"
}

func (c *MyController) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", c.Get)
	return mux
}

func (c *MyController) Get(w http.ResponseWriter, r *http.Request) {
	c.service.DoSomething()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type Service interface {
	DoSomething()
}

type MyService struct {
	logger fast.Logger
}

func NewMyService(ctx *fast.BuilderContext) *MyService {
	return &MyService{
		logger: fast.TypedLogger[MyService](ctx, fast.Scoped),
	}
}

func (s *MyService) DoSomething() {
	s.logger.Inf("Hello World")
}
```

```shell
map[*gofast.LogMiddleware:0x4000196188 main.Service:0x4000196190 gofast.ConfigProvider[github.com/ugozlave/gofast.AppConfig]:0x4000196160 *main.MyController:0x4000196180 *gofast.HealthController:0x4000196178 gofast.Logger:0x4000196168]
server start
time=2025-07-02T16:56:21.677Z level=INFO msg="request received" application=example environment=development service=gofast.LogMiddleware requestId=1 http.method=GET http.host=localhost:8080 http.url=/my/ http.remote=127.0.0.1:42884 http.agent=curl/8.7.1
time=2025-07-02T16:56:21.677Z level=INFO msg="Hello World" application=example environment=development service=main.MyService requestId=1
time=2025-07-02T16:56:21.677Z level=INFO msg="request finished" application=example environment=development service=gofast.LogMiddleware requestId=1 http.method=GET http.host=localhost:8080 http.url=/my/ http.remote=127.0.0.1:42884 http.agent=curl/8.7.1 http.status=200
time=2025-07-02T16:56:23.864Z level=INFO msg="request received" application=example environment=development service=gofast.LogMiddleware requestId=2 http.method=GET http.host=localhost:8080 http.url=/my/ http.remote=127.0.0.1:42890 http.agent=curl/8.7.1
time=2025-07-02T16:56:23.864Z level=INFO msg="Hello World" application=example environment=development service=main.MyService requestId=2
time=2025-07-02T16:56:23.864Z level=INFO msg="request finished" application=example environment=development service=gofast.LogMiddleware requestId=2 http.method=GET http.host=localhost:8080 http.url=/my/ http.remote=127.0.0.1:42890 http.agent=curl/8.7.1 http.status=200
^C
server stop
```

This is a more complete example demonstrating how to define a custom controller and a custom service with gofast.

A new controller (MyController) is registered with a /my route prefix.

A scoped service (MyService) is registered and injected into the controller.

The service uses the built-in structured logger, automatically scoped and typed.

The controller exposes a route (GET /my/) that calls the service and responds with "OK".

This setup shows how to cleanly separate concerns using dependency injection, and how to leverage gofast’s internal DI system and logging without any third-party packages.

## Documentation

WIP