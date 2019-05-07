package http

import (
	"fmt"
	"net"
	"time"

	"github.com/ddosakura/NPE/uri"

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

	// 貌似 baidu 没区分 Options 请求
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

func TestMultiRequest(t *testing.T) {
	/*
		get 180.97.33.107:80
		set 180.97.33.107:80
		get 180.97.33.108:80
		http://www.baidu.com
		set 180.97.33.108:80
		http://www.baidu.com/s?wd=golang

		golang 构造 TCPAddr 的时候会进行 DNS 查询?
	*/
	//Build("http://www.baidu.com").
	//	Conn(nil).
	//	Do(func(r *Response) {
	//		fmt.Println("http://www.baidu.com")
	//		fmt.Println(r.String())
	//	}, func(e error) {
	//		t.Log(e)
	//	})
	//Build("http://www.baidu.com/s?wd=golang").
	//	Conn(nil).
	//	Do(func(r *Response) {
	//		fmt.Println("http://www.baidu.com/s?wd=golang")
	//		fmt.Println(r.String())
	//	}, func(e error) {
	//		t.Log(e)
	//	})

	addr, err := net.ResolveTCPAddr("tcp4", "180.97.33.107:80")
	if err != nil {
		t.Fatal(err)
	}

	Build("http://baidu.com").
		Conn(addr).
		Do(func(r *Response) {
			fmt.Println("http://baidu.com")
			//fmt.Println("---", r)
			fmt.Println(r.String())
		}, func(e error) {
			t.Log(e)
		})
	Build("http:///www.baidu.com").
		Conn(addr).
		Do(func(r *Response) {
			fmt.Println("http://www.baidu.com")
			//fmt.Println("---", r)
			fmt.Println(r.String())
		}, func(e error) {
			t.Log(e)
		})

	time.Sleep(time.Second * 60)
}
