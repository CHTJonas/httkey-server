package utils

import (
	"errors"
	"net/url"
)

var (
	ErrInvalidScheme = errors.New("Invalid URL scheme")
	ErrNoHost        = errors.New("No URL hostname")
)

func ParseURL(rawurl string) (*url.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrInvalidScheme
	}
	if u.Host == "" {
		return nil, ErrNoHost
	}
	return u, nil
}
