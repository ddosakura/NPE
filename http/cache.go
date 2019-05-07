package http

import (
	"fmt"
	"net"
	"sync"
)

var (
	connMutex sync.Mutex
	connMap   map[string]*net.TCPConn = make(map[string]*net.TCPConn)
)

func cacheConn(addr *net.TCPAddr) *net.TCPConn {
	defer connMutex.Unlock()
	connMutex.Lock()
	fmt.Println("get", addr.String())
	return connMap[addr.String()]
}

func connCache(addr *net.TCPAddr, conn *net.TCPConn) {
	defer connMutex.Unlock()
	connMutex.Lock()
	fmt.Println("set", addr.String())
	connMap[addr.String()] = conn
}
