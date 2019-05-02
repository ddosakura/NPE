package main

import (
	"github.com/kr/pretty"
	"fmt"
	"net/http"
	"testing"
)

func TestProxy(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		"http://studyzy.cnblogs.com",
		nil)
	if err != nil {
		panic(err)
	}
	pretty.Println(req)

	p := Proxy(ProxyConfig{
		Addr: "127.0.0.1:8888",
	})

	res, err := p.Request(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
