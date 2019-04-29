package uri

import (
	"strconv"
)

// Authority for URI
type Authority struct {
	uri      *URI
	UserInfo string
	Host     string
	Port     int16
}

func newAuthority(auth string, uri *URI) {
	data := regAuthority.FindStringSubmatch(auth)
	var port int16
	if data[5] == "" {
		port = -1
	} else {
		p, _ := strconv.Atoi(data[5])
		port = int16(p)
	}

	uri.Authority = &Authority{
		uri:      uri,
		UserInfo: data[2],
		Host:     data[3],
		Port:     port,
	}
}

func (a *Authority) String() (s string) {
	if a.Host != "" {
		s = a.Host
	}
	if a.UserInfo != "" {
		s = a.UserInfo + "@" + s
	}
	if a.Port >= 0 &&
		(a.uri == nil || a.Port != SchemePort[a.uri.Scheme]) {
		s += ":" + strconv.Itoa(int(a.Port))
	}
	return
}

var (
	// SchemePort for URI
	SchemePort = map[string]int16{
		"http":  80,
		"https": 443,
	}
)
