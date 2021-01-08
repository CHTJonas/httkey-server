package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(path, addr string, readTimeout time.Duration) *http.Server {
	ks := NewKeyserver(path)
	r := mux.NewRouter()
	r.MatcherFunc(alwaysMatch).Handler(ks)
	r.Use(headerMiddleware)
	r.Use(logMiddleware)
	return &http.Server{
		Handler:     r,
		Addr:        addr,
		ReadTimeout: readTimeout,
	}
}

func alwaysMatch(_ *http.Request, _ *mux.RouteMatch) bool {
	return true
}
