package gofast

import (
	"net/http"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}
