package gofast

import (
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
}

func NewHealthController(_ *BuilderContext) *HealthController {
	return &HealthController{}
}

func (c *HealthController) Prefix() string {
	return "/health"
}

func (c *HealthController) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", c.handle)
	return mux
}

func (c *HealthController) handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
