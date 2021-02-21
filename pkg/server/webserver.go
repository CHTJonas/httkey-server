package server

import (
	"context"
	"net/http"
	"time"

	"github.com/CHTJonas/httkey-server/assets"
	"github.com/gorilla/mux"
)

type Webserver struct {
	r   *mux.Router
	ks  *Keyserver
	srv *http.Server
}

func NewWebserver(path, addr string, readTimeout time.Duration) *Webserver {
	ks := NewKeyserver(path)

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(assets.Server())
	r.MatcherFunc(alwaysMatch).Handler(ks)
	r.Use(recoveryMiddleware)

	srv := &http.Server{
		Handler:     r,
		Addr:        addr,
		ReadTimeout: readTimeout,
	}

	return &Webserver{
		r:   r,
		ks:  ks,
		srv: srv,
	}
}

func (w *Webserver) RegisterMiddleware(mwf mux.MiddlewareFunc) {
	w.r.Use(mwf)
}

func (w *Webserver) ListenAndServe() error {
	return w.srv.ListenAndServe()
}

func (w *Webserver) Shutdown(ctx context.Context) error {
	return w.srv.Shutdown(ctx)
}

func alwaysMatch(_ *http.Request, _ *mux.RouteMatch) bool {
	return true
}
