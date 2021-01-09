package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
)

func DefaultLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := r.URL.Scheme
		if len(scheme) == 0 {
			scheme = "http"
		}
		u := &url.URL{
			Scheme: scheme,
			Host:   r.Host,
			Path:   r.URL.Path,
		}
		httpAction := fmt.Sprintf("\"%s %s %s\"", r.Method, u.String(), r.Proto)
		log.Println(r.RemoteAddr, httpAction)
		next.ServeHTTP(w, r)
	})
}

func ApacheLogMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func ProxyMiddleware(next http.Handler) http.Handler {
	return handlers.ProxyHeaders(next)
}

func serverHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "https://github.com/CHTJonas/httkey-server")
		next.ServeHTTP(w, r)
	})
}
