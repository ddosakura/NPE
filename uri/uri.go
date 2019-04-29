package uri

import (
	"regexp"
)

var (
	regURI = regexp.MustCompile(`^(([^:/?#]+):)?(\/\/([^/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)
	// exclude IPv6
	regAuthority = regexp.MustCompile(`^(([^/?#@]*)@)?([^/?#:]*)(:([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))?$`)
)

// Parse URI string
func Parse(uri string) []string {
	return regURI.FindStringSubmatch(uri)
}

// URI struct
type URI struct {
	Scheme string
	// [ userinfo "@" ] host [ ":" port ]
	Authority *Authority // ([^/?#]*)
	Path      string
	Query     string
	Fragment  string
}

// NewURI by Parse
func NewURI(uri string) *URI {
	data := Parse(uri)
	u := &URI{
		Scheme:   data[2],
		Path:     data[5],
		Query:    data[7],
		Fragment: data[9],
	}
	newAuthority(data[4], u)
	return u
}

func (u *URI) String() (s string) {
	if u.Scheme != "" {
		s += u.Scheme + ":"
	}
	if u.Authority.Host != "" {
		s += "//" + u.Authority.String()
	}
	if u.Path != "" {
		s += u.Path
	}
	if u.Query != "" {
		s += "?" + u.Query
	}
	if u.Fragment != "" {
		s += "#" + u.Fragment
	}
	return
}
