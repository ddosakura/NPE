package uri

import (
	"github.com/davecgh/go-spew/spew"
	//"github.com/kr/pretty"
	"testing"
)

func doParse(t *testing.T, uri string) {
	ans := Parse(uri)
	spew.Dump(ans)
	// pretty.Println(ans)
}

func TestParse(t *testing.T) {
	doParse(t, "http://baidu.com")
	doParse(t, "https://user:pass@baidu.com/news?q=abc#2333")
}

func testNew(t *testing.T, uri string, expect URI) {
	t.Log("testing", expect.String())
	ans := NewURI(uri)
	// spew.Dump(ans)
	//pretty.Println(ans, expect)
	if ans.String() != expect.String() {
		t.Fatal(uri, "expect", expect.String(), "but", ans.String())
	}
	t.Log(ans.String(), "pass")
}

func TestNew(t *testing.T) {
	testNew(t, "http://baidu.com", URI{
		Scheme: "http",
		Authority: &Authority{
			UserInfo: "",
			Host:     "baidu.com",
			// Port:     80,
			Port: -1,
		},
		Path:     "",
		Query:    "",
		Fragment: "",
	})
	testNew(t, "http://baidu.com:8080", URI{
		Scheme: "http",
		Authority: &Authority{
			UserInfo: "",
			Host:     "baidu.com",
			Port:     8080,
		},
		Path:     "",
		Query:    "",
		Fragment: "",
	})
	testNew(t, "https://user:pass@baidu.com/news?q=abc#2333", URI{
		Scheme: "https",
		Authority: &Authority{
			UserInfo: "user:pass",
			Host:     "baidu.com",
			// Port:     443,
			Port: -1,
		},
		Path:     "/news",
		Query:    "q=abc",
		Fragment: "2333",
	})
	testNew(t, "htt#p://baidu.com", URI{
		Scheme: "",
		Authority: &Authority{
			UserInfo: "",
			Host:     "",
			Port:     -1,
		},
		Path:     "htt",
		Query:    "",
		Fragment: "p://baidu.com",
	})
	testNew(t, "baidu.com", URI{
		Scheme: "",
		Authority: &Authority{
			UserInfo: "",
			Host:     "",
			Port:     -1,
		},
		Path:     "baidu.com",
		Query:    "",
		Fragment: "",
	})

	testNew(t, "https://baidu.com", URI{
		Scheme: "https",
		Authority: &Authority{
			UserInfo: "",
			Host:     "baidu.com",
			// Port:     443,
			Port: -1,
		},
		Path:     "",
		Query:    "",
		Fragment: "",
	})
}
