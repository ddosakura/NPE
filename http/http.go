package http

import (
	"github.com/ddosakura/NPE/uri"
	"strings"
)

// Header for HTTP
type Header map[string]string

// Request for HTTP
type Request struct {
	// First Line
	Method  string
	Version string

	URI uri.URI

	Header Header
}

// MakeRequest for HTTP
func MakeRequest(URI string) *Request {
	u := uri.NewURI(URI)
	if u.Scheme != "http" && u.Scheme != "https" {
		u = uri.NewURI("http://" + URI)
	}

	if u.Scheme != "https" {
		u.Scheme = "http"
	}
	return NewRequest(*u)
}

// NewRequest for HTTP
func NewRequest(uri uri.URI) *Request {
	r := &Request{
		URI: uri,
	}
	r.SetMethod(GET)
	r.SetVersion(HTTP11)

	r.Header = map[string]string{
		"User-Agent": "DDoSakura HTTP Util",
		"Accept":     "*/*",
	}

	return r
}

func (r *Request) String() string {
	first := r.Method + " "
	if r.URI.Path == "" {
		first += "/"
	} else {
		first += r.URI.Path
	}
	first += " " + r.Version

	host := "Host: " + r.URI.Authority.HostPort()

	lines := make([]string, 0, len(r.Header)+2)
	lines = append(lines, first, host)
	for k, v := range r.Header {
		lines = append(lines, k+": "+v)
	}

	// TODO:
	// 1. \n? \r\n?
	// 2. a new line in last?
	return strings.Join(lines, "\r\n") + "\r\n"
}
