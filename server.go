package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var staticPath = "/tmp"

func myMiddleware(next http.Handler) http.Handler {
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
		log.Println(u.String())
		w.Header().Set("Server", "https://github.com/CHTJonas/httkey-server")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.MatcherFunc(alwaysMatch).HandlerFunc(serveFileHandler)
	r.Use(myMiddleware)

	srv := &http.Server{
		Handler:     r,
		Addr:        "127.0.0.1:8000",
		ReadTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
	log.Println("Starting server")

	c := make(chan os.Signal, 1)
	// SIGQUIT or SIGTERM will not be caught.
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Println("Waiting up to 30 seconds for server to shutdown")
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Shutdown error:", err.Error())
	}
	log.Println("Goodbye!")
	os.Exit(0)
}

func alwaysMatch(_ *http.Request, _ *mux.RouteMatch) bool {
	return true
}

func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err := SplitURLToHash(host, reqPath)
	if err != nil {
		http.Error(w, "400 Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	path := filepath.Join(staticPath, hash)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = "/" + hash
		http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r2)
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
