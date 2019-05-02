package main

import (
	"os"
)

func main() {
	err := ProxyServer(ProxyConfig{
		Addr: os.Args[1],
		//HTTPS: &TLS{
		//	Crt: "/home/moyinzi/projects/gDemo/NPE/httpProxy/server.crt",
		//	Key: "/home/moyinzi/projects/gDemo/NPE/httpProxy/server.key",
		//},
	})
	if err != nil {
		panic(err)
	}
}
