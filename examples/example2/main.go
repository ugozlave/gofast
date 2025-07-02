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
