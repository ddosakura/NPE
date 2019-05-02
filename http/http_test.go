package http

import (
	"fmt"
	"github.com/ddosakura/NPE/uri"
	"time"
	//"github.com/kr/pretty"
	"testing"
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

// go test -v -run ^TestRequest$
func TestRequest(t *testing.T) {
	r, e := Build("https://baidu.com").DoSync()
	if e != nil {
		t.Fatal(e)
	}
	if r.Code != 400 {
		t.Fatal("code error!")
	}

	r, e = Build("www.baidu.com").Build(func(r *Request) {
		fmt.Println(r)
		fmt.Println("--- --- ---")
	}).DoSync()
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r.String())
}

func TestOptions(t *testing.T) {
	r, e := Build("baidu.com").Build(func(r *Request) {
		r.SetMethod(OPTIONS)
		r.URI.Path = "*"
		fmt.Println(r)
		fmt.Println("--- --- ---")
	}).DoSync()
	if e != nil {
		t.Fatal(e)
	}

	// 貌似没区分Options请求
	fmt.Println(r.String())
}

// TODO: test after dohttps
func TestModified(t *testing.T) {
	r, e := Build("https://blog.moyinzi.top").Build(func(r *Request) {
		r.Header["If-Modified-Since"] = time.Now().UTC().String()
		fmt.Println(r)
		fmt.Println("--- --- ---")
	}).DoSync()
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r.String())
}
