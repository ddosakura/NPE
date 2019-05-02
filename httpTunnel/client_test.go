package tunnel

import (
	"io/ioutil"
	"testing"
)

func TestClient(t *testing.T) {
	//res, err := Client.Get("http://studyzy.cnblogs.com")
	res, err := Client.Get("http://baidu.com/")
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	println(string(bs))
}
