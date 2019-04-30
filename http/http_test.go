package http

import (
	"fmt"
	"github.com/ddosakura/NPE/uri"
	// "testing"
)

func ExampleNewRequest() {
	fmt.Println(NewRequest(*uri.NewURI("http://baidu.com")).String())
	fmt.Println(NewRequest(*uri.NewURI("http://baidu.com:8080")).String())
	fmt.Println(NewRequest(*uri.NewURI("https://baidu.com")).String())
	fmt.Println(NewRequest(*uri.NewURI("https://baidu.com:8443")).String())
	fmt.Println(MakeRequest("baidu.com").String())
	fmt.Println(MakeRequest("baidu.com:8080").String())
	fmt.Println(MakeRequest("http://baidu.com").String())
	fmt.Println(MakeRequest("http://baidu.com:8080").String())
	fmt.Println(MakeRequest("https://baidu.com").String())
	fmt.Println(MakeRequest("https://baidu.com:8443").String())
	fmt.Println(MakeRequest("baidu.com:80/news").String())
	fmt.Println(MakeRequest("baidu.com:8080/news").String())

	// Output:
	//
	// ...
}
