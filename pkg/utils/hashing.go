package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/spaolacci/murmur3"
)

func RawURLToHash(rawurl string) (string, error) {
	u, err := ParseURL(rawurl)
	if err != nil {
		return "", err
	}
	return URLToHash(u)
}

func URLToHash(u *url.URL) (string, error) {
	return SplitURLToHash(u.Host, u.Path)
}

func SplitURLToHash(host, path string) (string, error) {
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
	return fmt.Sprintf("%s%s", uintToString(h1), uintToString(h2))
}

func uintToString(x uint64) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, x)
	return hex.EncodeToString(buf)
}
