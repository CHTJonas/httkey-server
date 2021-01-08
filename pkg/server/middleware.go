package server

import (
	"log"
	"net/http"
	"net/url"
)

func logMiddleware(next http.Handler) http.Handler {
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
		log.Println(r.RemoteAddr, r.Method, u.String())
		next.ServeHTTP(w, r)
	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "https://github.com/CHTJonas/httkey-server")
		next.ServeHTTP(w, r)
	})
}
