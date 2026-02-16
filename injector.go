package gofast

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/ugozlave/cargo"
)

type Injector interface {
	Handler() http.Handler
}

type HttpInjector struct {
	ctn *cargo.Container
}

func NewHttpInjector(ctn *cargo.Container) *HttpInjector {
	return &HttpInjector{
		ctn: ctn,
	}
}

func (inj *HttpInjector) Handler() http.Handler {
	gen := Get[UniqueIDGenerator](NewBuilderContext(context.Background(), inj.ctn), Transient)
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
		ctx := NewBuilderContext(context.WithValue(r.Context(), CtxRequestId, id), inj.ctn)

		// build controllers
		handler := inj.Controllers(ctx)

		// build middlewares
		use := inj.Middlewares(ctx)

		mux := use(handler)

		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (inj *HttpInjector) Controllers(ctx *BuilderContext) http.Handler {
	mux := http.NewServeMux()
	for _, ctrl := range All[Controller](ctx, Scoped) {
		prefix := strings.Trim(ctrl.Prefix(), "/")
		mux.Handle("/"+prefix, StripPrefix("/"+prefix, ctrl.Routes()))
		mux.Handle("/"+prefix+"/", StripPrefix("/"+prefix, ctrl.Routes()))
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

func StripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		if len(p) < len(r.URL.Path) && !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		if len(rp) < len(r.URL.RawPath) && !strings.HasPrefix(rp, "/") {
			rp = "/" + rp
		}
		if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			r2.URL.RawPath = rp
			h.ServeHTTP(w, r2)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
