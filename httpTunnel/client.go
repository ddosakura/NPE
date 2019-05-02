package tunnel

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/kr/pretty"

	"golang.org/x/net/publicsuffix"
)

type jarList struct{}

func (jarList) PublicSuffix(domain string) string {
	ps, _ := publicsuffix.EffectiveTLDPlusOne(domain)
	return ps
}

func (jarList) String() string {
	return "publicsuffix.EffectiveTLDPlusOne"
}

func proxyURL(fixedURL *url.URL) func(*http.Request) (*url.URL, error) {
	// TODO: switch proxyURI
	return func(*http.Request) (*url.URL, error) {
		// TODO: set user/pass
		return fixedURL, nil
	}
}

// Client for Tunnel
var (
	regHasPort = regexp.MustCompile(`(.*)(:([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))$`)

	urli        = url.URL{}
	urlproxy, _ = urli.Parse("http://127.0.0.1:8888")

	ClientTransport = &http.Transport{
		// 代理
		//Proxy: proxyURL(urlproxy),
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			pretty.Println(ctx, network, addr)

			//tcpAddr, err := net.ResolveTCPAddr(network, addr)
			tcpAddr, err := net.ResolveTCPAddr(network, "127.0.0.1:8888")
			if err != nil {
				return nil, err
			}
			conn, err := net.DialTCP(network, nil, tcpAddr)
			if err != nil {
				return nil, err
			}

			conn.Write([]byte(
				"CONNECT 127.0.0.1:8889 HTTP/1.1\r\n" +
					"Host: 127.0.0.1:8889\r\n" +
					"Connection: Close\r\n" +
					"\r\n",
			))
			pass(conn)

			//host := regHasPort.FindStringSubmatch(addr)[1]
			conn.Write([]byte(
				"CONNECT " + addr + " HTTP/1.1\r\n" +
					"Host: " + addr + "\r\n" +
					"Connection: Close\r\n" +
					"\r\n",
			))
			pass(conn)

			//bs, _ := ioutil.ReadAll(conn)
			//println("proxy", string(bs))

			//for {
			//	buf := bufio.NewReader(conn)
			//	bs, _, err := buf.ReadLine()
			//	if err != nil {
			//		panic(err)
			//	}
			//	println(err, ">", string(bs))
			//	if len(bs) == 0 {
			//		break
			//	}
			//}

			return conn, nil
		},
		TLSClientConfig: &tls.Config{
			// 跳过证书验证
			InsecureSkipVerify: true,
		},
	}

	ClientCookieJar, _ = cookiejar.New(
		&cookiejar.Options{
			// PublicSuffixList: publicsuffix.List,
			PublicSuffixList: &jarList{},
		},
	)

	Client = &http.Client{
		Transport: ClientTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
		Jar:     ClientCookieJar,
		Timeout: 0,
	}
)

func pass(conn io.Reader) {
	buf := bufio.NewReader(conn)
	bs, _, _ := buf.ReadLine()
	println("proxy", string(bs))
}
