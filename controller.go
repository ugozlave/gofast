package gofast

import (
	"net/http"
)

type Controller interface {
	Prefix() string
	Routes() http.Handler
}
