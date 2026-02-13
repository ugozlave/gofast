package faster

import (
	"net/http"

	"github.com/ugozlave/gofast"
)

/*
** HealthController
 */

type HealthController struct {
}

func HealthControllerBuilder() Builder[*HealthController] {
	return func(ctx *gofast.BuilderContext) *HealthController {
		return NewHealthController(ctx)
	}
}

func NewHealthController(ctx *gofast.BuilderContext) *HealthController {
	return &HealthController{}
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
