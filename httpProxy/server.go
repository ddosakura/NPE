package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ProxyConfig for C/S
type ProxyConfig struct {
	Addr  string
	HTTPS *TLS

	BasicAuth *BasicAuth
}

// TLS for Proxy
type TLS struct {
	Crt string
	Key string
}

// BasicAuth for Proxy
type BasicAuth struct {
	User string
	Pass string
}

// ProxyServer HTTP(S)
func ProxyServer(c ProxyConfig) error {
	if c.HTTPS != nil {
		fmt.Println("https server")
		return http.ListenAndServeTLS(c.Addr,
			c.HTTPS.Crt, c.HTTPS.Key,
			handler{
				config: &c,
			})
	}
	fmt.Println("http server")
	return http.ListenAndServe(c.Addr,
		handler{
			config: &c,
		})
}

type handler struct {
	config *ProxyConfig
}

/*
func (h handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("This is a proxy server!"))
		return
	}
	if h.config.BasicAuth != nil {
		u, p, ok := r.BasicAuth()
		if !ok ||
			u != h.config.BasicAuth.User ||
			p != h.config.BasicAuth.Pass {
			rw.WriteHeader(http.StatusProxyAuthRequired)
			return
		}
	}

	//rw.Header().Set("Transfer-Encoding", "chunked")
	//_, err := rw.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")//)
	//if err != nil {
	//	panic(err)
	//}

	host := r.Host
	ra := r.RemoteAddr
	// www.baidu.com:443 127.0.0.1:55906 www.baidu.com:443 {   www.baidu.com:443   false  }
	fmt.Println(host, ra, r.RequestURI, *r.URL)
	if !strings.Contains(host, ":") {
		host += ":80"
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	//defer conn.Close()
	//go func() {
	//	defer r.Body.Close()
	//	//defer conn.Close()
	//	//_, err := io.CopyBuffer(conn, r.Body, make([]byte, 1024))
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//	io.CopyBuffer(conn, r.Body, make([]byte, 1024))
	//}()
	////_, err = io.CopyBuffer(rw, conn, make([]byte, 1024))
	////if err != nil {
	////	panic(err)
	////}
	//io.CopyBuffer(rw, conn, make([]byte, 1024))

	go func() {
		// defer r.Body.Close()
		bufio.NewReader(r.Body).WriteTo(conn)
	}()
	// defer conn.Close()
	bufio.NewReader(conn).WriteTo(rw)
}
*/

var (
	regHasPort = regexp.MustCompile(`(.*)(:([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))$`)
)

// Hop-by-hop headers
// These headers are meaningful only for a single transport-level connection
// and must not be retransmitted by proxies or cached. Such headers are:
// Connection, Keep-Alive, Proxy-Authenticate, Proxy-Authorization, TE,
// Trailer, Transfer-Encoding and Upgrade.
// Note that only hop-by-hop headers may be set using the Connection general header.
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			switch k {
			case "Connection", "Keep-Alive", "Proxy-Authenticate", "Trailer", "Transfer-Encoding", "Upgrade":
			default:
				dst.Add(k, v)
			}
		}
	}
}

var (
	errAuthMethodUnsupported = errors.New("auth method unsupported")
	errAuthNull              = errors.New("user/pass is null")
)

func (h handler) parseBaseCredential(basicCredential string) (user string, pass string, err error) {
	auths := strings.SplitN(basicCredential, " ", 2)
	if len(auths) != 2 {
		return "", "", errAuthMethodUnsupported
	}
	authMethod := auths[0]
	authB64 := auths[1]
	switch authMethod {
	case "Basic":
		authstr, err := base64.StdEncoding.DecodeString(authB64)
		if err != nil {
			return "", "", err
		}
		//fmt.Println(string(authstr))
		userPwd := strings.SplitN(string(authstr), ":", 2)
		if len(userPwd) != 2 {
			return "", "", errAuthNull
		}
		user = userPwd[0]
		pass = userPwd[1]
	default:
		return "", "", errAuthMethodUnsupported
	}
	return
}

// https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
func (h handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		res, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer res.Body.Close()
		copyHeader(rw.Header(), res.Header)
		rw.WriteHeader(res.StatusCode)
		io.Copy(rw, res.Body)
		return
	}
	if h.config.BasicAuth != nil {
		//u, p, ok := r.BasicAuth()
		//if !ok ||
		//	u != h.config.BasicAuth.User ||
		//	p != h.config.BasicAuth.Pass {
		//	http.Error(rw,
		//		"Username/Password Error",
		//		http.StatusProxyAuthRequired)
		//	return
		//}

		// TODO: auth

		// https://www.jb51.net/article/89070.htm
		auth := r.Header["Proxy-Authorization"]
		if auth == nil || len(auth) < 1 {
			rw.Header().Set("Proxy-Authenticate", `Basic realm="*"`)
			http.Error(rw,
				"Username/Password Error",
				http.StatusProxyAuthRequired)
			return
		}
		u, p, err := h.parseBaseCredential(auth[0])
		if err != nil {
			http.Error(rw,
				err.Error(),
				http.StatusProxyAuthRequired)
		}
		if u != h.config.BasicAuth.User || p != h.config.BasicAuth.Pass {
			http.Error(rw,
				"Username/Password Error",
				http.StatusProxyAuthRequired)
			return
		}
	}

	// TODO: HTTP/2
	if r.ProtoMajor == 2 {
		http.Error(rw,
			"HTTP/2 not supported",
			http.StatusInternalServerError)
	}

	// TODO: remove?
	host := r.Host
	if !regHasPort.MatchString(host) {
		switch r.URL.Scheme {
		case "https":
			host += ":443"
		case "http":
			host += ":80"
		default:
			http.Error(rw,
				"Need http/https",
				http.StatusBadRequest)
			return
		}
	}
	fmt.Println(host)
	conn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		http.Error(rw,
			err.Error(),
			http.StatusServiceUnavailable)
		return
	}
	//rw.WriteHeader(http.StatusOK)
	// TODO: OK? Connection established? Connection Established?
	rw.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	fmt.Println(host, "conn ok!")

	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw,
			"Hijacking not supported",
			http.StatusInternalServerError)
		return
	}
	client, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(rw,
			err.Error(),
			http.StatusServiceUnavailable)
		return
	}

	go transfer(conn, client)
	go transfer(client, conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	//io.Copy(destination, source)
	io.CopyBuffer(destination, source, make([]byte, 1024))
}
