package main

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/spaolacci/murmur3"
)

func RawURLToHash(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return URLToHash(u)
}

func URLToHash(u *url.URL) (string, error) {
	return SplitURLToHash(u.Host, u.Path)
}

func SplitURLToHash(host, path string) (string, error) {
	//
	_, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return hash(host, path), nil
}

func hash(hostname, path string) string {
	seed := murmur3.Sum32([]byte(hostname))
	hasher := murmur3.New128WithSeed(seed)
	hasher.Write([]byte(path))
	h1, h2 := hasher.Sum128()
	return fmt.Sprintf("%d%d", h1, h2)
}

// func splitHostPort(hostport string) (host, port string) {
// 	host = hostport
// 	colon := strings.LastIndexByte(host, ':')
// 	if colon != -1 && validOptionalPort(host[colon:]) {
// 		host, port = host[:colon], host[colon+1:]
// 	}
// 	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
// 		host = host[1 : len(host)-1]
// 	}
// 	return
// }
