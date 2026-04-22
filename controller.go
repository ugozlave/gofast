package gofast

import (
	"encoding/json"
	"net/http"
)

type Controller interface {
	Prefix() string
	Routes() http.Handler
}

/*
** HealthController
 */

type HealthController struct {
	Services []HealthChecker
}

func HealthControllerBuilder() Builder[*HealthController] {
	return func(ctx *BuilderContext) *HealthController {
		return &HealthController{
			Services: All[HealthChecker](ctx, Scoped),
		}
	}
}

func (c *HealthController) Prefix() string {
	return "health"
}

func (c *HealthController) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", c.handle)
	return mux
}

func (c *HealthController) handle(w http.ResponseWriter, r *http.Request) {
	status := map[string]bool{}
	code := http.StatusOK
	for _, service := range c.Services {
		name, healthy, err := service.HealthCheck()
		if err != nil {
			code = http.StatusServiceUnavailable
		}
		status[name] = healthy

	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(status)
}
