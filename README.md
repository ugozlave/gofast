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
	app := fast.New(faster.NewAppConfig())
	app.Run()
}
```

This minimal example starts a default web API application using the built-in HTTP server.

## Documentation

WIP

## Examples

https://github.com/ugozlave/gofast-examples