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
	gen UniqueIDGenerator
	ctn *cargo.Container
	ctx context.Context
}

func NewHttpInjector(gen UniqueIDGenerator, ctn *cargo.Container, ctx context.Context) *HttpInjector {
	return &HttpInjector{
		gen: gen,
		ctn: ctn,
		ctx: ctx,
	}
}

func (inj *HttpInjector) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create unique request ID
		id := inj.gen.Next()

		// create unique scope for the request
		scope := fmt.Sprintf(ScopeRequestKeyFormat, id)
		inj.ctn.CreateScope(scope)
		defer func() {
			inj.ctn.DeleteScope(scope)
		}()

		// create a new builder context
		ctx := NewBuilderContext(context.WithValue(inj.ctx, CtxRequestId, id), inj.ctn)

		// build controllers
		handler := inj.Controllers(ctx, scope)

		// build middlewares
		use := inj.Middlewares(ctx, scope)

		mux := use(handler)

		mux.ServeHTTP(w, r)
	})
}

func (inj *HttpInjector) Controllers(ctx cargo.BuilderContext, scope string) http.Handler {
	mux := http.NewServeMux()
	for _, ctrl := range cargo.All[Controller](inj.ctn, scope, ctx) {
		prefix := strings.TrimSuffix(ctrl.Prefix(), "/")
		mux.Handle(prefix+"/", http.StripPrefix(prefix, ctrl.Routes()))
	}
	return mux
}

func (inj *HttpInjector) Middlewares(ctx cargo.BuilderContext, scope string) func(http.Handler) http.Handler {
	return func(mux http.Handler) http.Handler {
		for _, mid := range slices.Backward(cargo.All[Middleware](inj.ctn, scope, ctx)) {
			mux = mid.Handle(mux)
		}
		return mux
	}
}
