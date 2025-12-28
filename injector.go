package gofast

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/ugozlave/cargo"
)

type Injector interface {
	Handler() http.Handler
}

type HttpInjector struct {
	ctn *cargo.Container
	ctx context.Context
}

func NewHttpInjector(ctn *cargo.Container, ctx context.Context) *HttpInjector {
	return &HttpInjector{
		ctn: ctn,
		ctx: ctx,
	}
}

func (inj *HttpInjector) Handler() http.Handler {
	gen := Get[UniqueIDGenerator](NewBuilderContext(inj.ctx, inj.ctn), Transient)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create unique request ID
		id := gen.Next()

		// create unique scope for the request
		scope := fmt.Sprintf(ScopeRequestKeyFormat, id)
		inj.ctn.CreateScope(scope)
		defer func() {
			inj.ctn.DeleteScope(scope)
		}()

		// create a new builder context
		ctx := NewBuilderContext(context.WithValue(inj.ctx, CtxRequestId, id), inj.ctn)

		// build controllers
		handler := inj.Controllers(ctx)

		// build middlewares
		use := inj.Middlewares(ctx)

		mux := use(handler)

		mux.ServeHTTP(w, r)
	})
}

func (inj *HttpInjector) Controllers(ctx *BuilderContext) http.Handler {
	mux := http.NewServeMux()
	for _, ctrl := range All[Controller](ctx, Scoped) {
		prefix := strings.Trim(ctrl.Prefix(), "/")
		mux.Handle("/"+prefix+"/", http.StripPrefix("/"+prefix, ctrl.Routes()))
	}
	return mux
}

func (inj *HttpInjector) Middlewares(ctx *BuilderContext) func(http.Handler) http.Handler {
	return func(mux http.Handler) http.Handler {
		for _, mid := range slices.Backward(All[Middleware](ctx, Scoped)) {
			mux = mid.Handle(mux)
		}
		return mux
	}
}
