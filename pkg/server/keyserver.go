package server

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/CHTJonas/httkey-server/pkg/utils"
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
		http.Error(w, "500 Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err := utils.SplitURLToHash(host, reqPath)
	if err != nil {
		http.Error(w, "400 Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	path := filepath.Join(k.StaticPath, hash)
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
		r2 := k.setRequestPathToHash(hash, r)
		http.FileServer(http.Dir(k.StaticPath)).ServeHTTP(w, r2)
	default:
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
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
