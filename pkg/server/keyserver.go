package server

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

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
		errorHTML(w, "500 Internal Server Error", err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err := utils.SplitURLToHash(host, reqPath)
	if err != nil {
		errorHTML(w, "400 Bad Request", err.Error(), http.StatusBadRequest)
		return
	}

	path := filepath.Join(k.StaticPath, hash)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		errorHTML(w, "404 Not Found", "The URL does not exist.", http.StatusNotFound)
		return
	} else if err != nil {
		errorHTML(w, "500 Internal Server Error", err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		r2 := k.setRequestPathToHash(hash, r)
		http.FileServer(http.Dir(k.StaticPath)).ServeHTTP(w, r2)
	default:
		errorHTML(w, "405 Method Not Allowed", "The request method is inappropriate.", http.StatusMethodNotAllowed)
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

func errorHTML(w http.ResponseWriter, title, message string, statusCode int) {
	template := assets.MustAsset("assets/error.html.mustache")
	context := map[string]string{"title": title, "message": message}
	body, err := mustache.Render(string(template), context)
	if err != nil {
		http.Error(w, message, statusCode)
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, body)
}
