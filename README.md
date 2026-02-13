# gofast

`gofast` is a lightweight, high-performance web framework for Go, designed to provide a clean and flexible structure with built-in dependency injection.

It supports three instance lifetimes:
- **Singleton**: a single instance shared throughout the application's lifetime.
- **Scoped**: a new instance created for each HTTP request.
- **Transient**: a fresh instance on every resolution.

`gofast` is built around four core interfaces that you implement according to your needs:
- `Controller`: for handling routes and endpoints.
- `Middleware`: for intercepting and processing requests.
- `Config`: for loading and managing application configuration.
- `Logger`: for customizable, extensible logging.

This design keeps the framework highly generic, giving you full control over implementations (databases, logging, configuration, etc.) without locking you into specific choices.

Perfect for building modern REST APIs, scalable web services, or microservices, `gofast` combines simplicity, performance, and modularity.

### Prerequisites

`gofast` requires [Go](https://go.dev/) version [1.24](https://go.dev/doc/devel/release#go1.24) or above.

## Installation

```shell
go get github.com/ugozlave/gofast
```

## Getting Started

```go
package main

import (
  "context"
  "os"
	"os/signal"

  "github.com/ugozlave/gofast/faster"
)

func main() {
  ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := faster.New()
	app.Run(ctx)
}
```

This minimal example starts a default web API application using the built-in HTTP server.

## Faster: Built-in Utilities for `gofast`

The `faster` package provides a ready-to-use set of common utilities that integrate seamlessly with `gofast`. It is designed to help you get started quickly while keeping full compatibility with the core framework's extensibility.

#### Features

- **Config Provider**  
  A flexible configuration system that loads values from json configuration files and overrides them with environment-specific files and environment variables.  
  Implements the core `Config` interface out of the box.

- **Structured Logger**  
  A simple structured logger based on `log/slog` package.  
  Implements the core `Logger` interface.

- **Health Controller**  
  A pre-built health check controller exposing a `/health` endpoint that returns `OK` when the application is running.  
  Ready to register with a single line.

- **Logging Middleware**  
  Automatically logs key request information (method, path, duration, status code, client IP, etc.) in structured format.

- **Recovery Middleware**  
  Catches panics in handlers or middlewares, logs the stack trace, and returns a clean `500 Internal Server Error` response without crashing the server.

## Examples

https://github.com/ugozlave/gofast-examples