package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/gorilla/handlers"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)
		scheme := r.URL.Scheme
		if len(scheme) == 0 {
			scheme = "http"
		}
		u := &url.URL{
			Scheme: scheme,
			Host:   r.Host,
			Path:   r.URL.Path,
		}
		addr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			addr = r.RemoteAddr
		}
		httpInfo := fmt.Sprintf("\"%s %s %s\"", r.Method, u.String(), r.Proto)
		refererInfo := fmt.Sprintf("\"%s\"", r.Referer())
		if refererInfo == "\"\"" {
			refererInfo = "\"-\""
		}
		uaInfo := fmt.Sprintf("\"%s\"", r.UserAgent())
		if uaInfo == "\"\"" {
			uaInfo = "\"-\""
		}
		log.Println(addr, httpInfo, lrw.statusCode, refererInfo, uaInfo)
	})
}

func ProxyMiddleware(next http.Handler) http.Handler {
	return handlers.ProxyHeaders(next)
}

func ServerHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "https://github.com/CHTJonas/httkey-server")
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				log.Println("An internal server error occurred:", err)
				log.Println(string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
