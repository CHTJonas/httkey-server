package server

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/CHTJonas/httkey-server/assets"
	"github.com/CHTJonas/httkey-server/pkg/utils"
	"github.com/cbroglie/mustache"
)

type Keyserver struct {
	StaticPath string
}

func NewKeyserver(staticPath string) *Keyserver {
	return &Keyserver{
		StaticPath: staticPath,
	}
}

func (k *Keyserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		respondWithError(w, r, "500 Internal Server Error", err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err := utils.SplitURLToHash(host, reqPath)
	if err != nil {
		respondWithError(w, r, "400 Bad Request", err.Error(), http.StatusBadRequest)
		return
	}

	path := filepath.Join(k.StaticPath, hash)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		respondWithError(w, r, "404 Not Found", "The URL does not exist.", http.StatusNotFound)
		return
	} else if err != nil {
		respondWithError(w, r, "500 Internal Server Error", err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		r2 := k.setRequestPathToHash(hash, r)
		http.FileServer(http.Dir(k.StaticPath)).ServeHTTP(w, r2)
	default:
		respondWithError(w, r, "405 Method Not Allowed", "The request method is inappropriate.", http.StatusMethodNotAllowed)
	}
}

func (k *Keyserver) setRequestPathToHash(hash string, r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = "/" + hash
	return r2
}

func respondWithError(w http.ResponseWriter, r *http.Request, title, message string, statusCode int) {
	body := fmt.Sprintf("%s\n%s\n", title, message)
	if clientAcceptsHTML(r) {
		template := assets.MustAsset("assets/error.html.mustache")
		context := map[string]string{"title": title, "message": message}
		html, err := mustache.Render(string(template), context)
		if err != nil {
			http.Error(w, message, statusCode)
			return
		}
		body = html
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, body)
}

func clientAcceptsHTML(r *http.Request) bool {
	h := r.Header.Get("Accept")
	for _, s := range strings.Split(h, ",") {
		for _, t := range strings.Split(s, ";") {
			if t == "text/html" {
				return true
			}
		}
	}
	return false
}
