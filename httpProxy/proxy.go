package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

var (
	//http cookie接口
	cookieJar, _ = cookiejar.New(nil)
	//跳过证书验证
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// UnsafeClient jump TLS
	UnsafeClient = &http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}
)

// ProxyClient for HTTP
type ProxyClient struct {
	Chain  *ProxyClient
	Config *ProxyConfig
}

// Proxy Builder
func Proxy(c ProxyConfig) *ProxyClient {
	p := &ProxyClient{}
	p.Config = &c
	return p

}

// Proxy Builder
func (p *ProxyClient) Proxy(c ProxyConfig) *ProxyClient {
	pc := &ProxyClient{}
	pc.Config = &c
	pc.Chain = p
	return pc
}

// Connect Proxy
func (p *ProxyClient) Connect(r io.Reader) error {
	if p.Chain != nil {
		p.Chain.Connect(r)
	}

	prefix := "http://"
	if p.Config.HTTPS != nil {
		prefix = "https://"
	}
	req, err := http.NewRequest(http.MethodConnect, prefix+p.Config.Addr, r)
	if err != nil {
		return err
	}
	if p.Config.BasicAuth != nil {
		req.SetBasicAuth(p.Config.BasicAuth.User, p.Config.BasicAuth.Pass)
	}

	client := http.DefaultClient
	if p.Config.HTTPS != nil {
		client = UnsafeClient
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(res)

	return nil
}

// Request by Proxy
func (p *ProxyClient) Request(r *http.Request) (*http.Response, error) {
	err := p.Connect(r.Body)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(r)
}
