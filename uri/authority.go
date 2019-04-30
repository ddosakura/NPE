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

func newAuthority(auth string, uri *URI) (a *Authority) {
	data := regAuthority.FindStringSubmatch(auth)
	var port int16
	if data[5] == "" {
		port = -1
	} else {
		p, _ := strconv.Atoi(data[5])
		port = int16(p)
	}

	a = &Authority{
		uri:      uri,
		UserInfo: data[2],
		Host:     data[3],
		Port:     port,
	}
	if uri == nil {
		return
	}
	uri.Authority = a
	return
}

func (a *Authority) String() (s string) {
	hp := a.HostPort()
	if hp != "" && a.UserInfo != "" {
		s = a.UserInfo + "@" + hp
	}
	return
}

// HostPort return host:port
func (a *Authority) HostPort() (s string) {
	if a.Host != "" {
		s = a.Host
	} else {
		return ""
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
