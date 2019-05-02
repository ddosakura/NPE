package http

import (
	"bufio"

	"github.com/ddosakura/NPE/uri"

	//"github.com/kr/pretty"
	"io"
	"net"
	"strconv"
	"strings"
)

// Header for HTTP
type Header map[string]string

// Request for HTTP
type Request struct {
	// First Line
	Method  string
	Version string

	URI uri.URI

	Header  Header
	Content *Content

	ProxyAddr string
}

// MakeRequest for HTTP
func MakeRequest(URI string) *Request {
	u := uri.NewURI(URI)
	// pretty.Println(URI, u)
	if u.Scheme != "http" && u.Scheme != "https" {
		u = uri.NewURI("http://" + URI)
	}

	if u.Scheme != "https" {
		u.Scheme = "http"
	}
	return NewRequest(*u)
}

// NewRequest for HTTP
func NewRequest(uri uri.URI) *Request {
	r := &Request{
		URI: uri,
	}
	r.SetMethod(GET)
	r.SetVersion(HTTP11)

	r.Header = map[string]string{
		// "User-Agent": "DDoSakura HTTP Util",
		"User-Agent": "curl/7.64.1",
		"Accept":     "*/*",
	}

	return r
}

func (r *Request) str() string {
	first := r.Method + " "
	if r.URI.Path == "" {
		first += "/"
	} else {
		first += r.URI.Path
	}
	first += " " + r.Version

	host := "Host: " + r.URI.Authority.HostPort()

	ex := 2
	if r.Content != nil {
		ex += 2
	}
	lines := make([]string, 0, len(r.Header)+ex)
	lines = append(lines, first, host)
	if r.Content != nil {
		lines = append(lines, r.Content.TypeHeader(), r.Content.LengthHeader())
	}
	for k, v := range r.Header {
		lines = append(lines, k+": "+v)
	}

	return strings.Join(lines, "\r\n") + "\r\n" + "\r\n"
}

func (r *Request) String() string {
	body := ""
	if r.Content != nil {
		body = string(r.Content.Data)
	}
	return r.str() + body
}

// Bytes of Request
func (r *Request) Bytes() (bs []byte) {
	bs = []byte(r.str())
	if r.Content != nil {
		bs = append(bs, r.Content.Data...)
	}
	return
}

// Response for HTTP
type Response struct {
	Version string
	Code    StatusCode
	Msg     string

	Header  Header
	Content *Content
}

// ParseResponse for HTTP
func ParseResponse(r io.Reader) *Response {
	//ioutil.ReadAll(r)
	buf := bufio.NewReader(r)
	bs, _, err := buf.ReadLine()
	if err != nil {
		return nil
	}
	line := string(bs)
	cs := strings.Split(line, " ")
	ver := cs[0]
	code, err := strconv.Atoi(cs[1])
	if err != nil {
		return nil
	}
	res := &Response{
		Version: ver,
		Code:    StatusCode(code),
		Msg:     strings.Join(cs[2:], ""),
		Header:  make(Header),
	}
	for {
		line = ""
		for {
			l, isPrefix, err := buf.ReadLine()
			if err != nil {
				return nil
			}
			line += string(l)
			if !isPrefix {
				break
			}
		}
		if line == "" {
			break
		}
		cs := strings.Split(line, ": ")
		k := cs[0]
		v := strings.Join(cs[1:], ": ")
		if strings.HasPrefix(k, "Content-") {
			if res.Content == nil {
				res.Content = &Content{
					Type: "text/plain",
					Data: nil,
				}
			}
			switch k {
			case "Content-Type":
				res.Content.Type = v
			case "Content-Length":
				len, err := strconv.Atoi(v)
				if err != nil {
					return nil
				}
				// res.Content.Len = len

				//res.Content.Data = make([]byte, 0, len)
				res.Content.Data = make([]byte, len)
			}
		} else {
			res.Header[k] = v
		}
	}
	n, err := buf.Read(res.Content.Data)
	if err != nil && n == len(res.Content.Data) {
		return nil
	}
	return res
}

func (r *Response) str() string {
	first := r.Version + " " + strconv.Itoa(int(r.Code)) + " " + r.Msg

	ex := 1
	if r.Content != nil {
		ex += 2
	}
	lines := make([]string, 0, len(r.Header)+ex)
	lines = append(lines, first)
	if r.Content != nil {
		lines = append(lines, r.Content.TypeHeader(), r.Content.LengthHeader())
	}
	for k, v := range r.Header {
		lines = append(lines, k+": "+v)
	}

	return strings.Join(lines, "\r\n") + "\r\n" + "\r\n"
}

func (r *Response) String() string {
	body := ""
	if r.Content != nil {
		body = string(r.Content.Data)
	}
	return r.str() + body
}

// R Flow
type R struct {
	req *Request
}

// Build Request
func Build(uri string) *R {
	return &R{
		req: MakeRequest(uri),
	}
}

// Build Request
func (r *R) Build(fn func(*Request)) *R {
	fn(r.req)
	return r
}

// DoSync Request
func (r *R) DoSync() (*Response, error) {
	host := r.req.URI.Authority.Host
	var tcpAddr *net.TCPAddr
	var err error
	if r.req.ProxyAddr == "" {
		port := r.req.URI.Authority.Port
		if port < 0 {
			port = uri.SchemePort[r.req.URI.Scheme]
		}
		tcpAddr, err = net.ResolveTCPAddr("tcp4", host+":"+strconv.Itoa(int(port)))
		if err != nil {
			return nil, err
		}
	} else {
		tcpAddr, err = net.ResolveTCPAddr("tcp4", r.req.ProxyAddr)
		if err != nil {
			return nil, err
		}
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(MakeRequest(host).Bytes())
	if err != nil {
		return nil, err
	}
	return ParseResponse(conn), nil
}

// Do Request
func (r *R) Do(fn func(*Response), e func(error)) {
	go func() {
		res, err := r.DoSync()
		if err != nil {
			e(err)
		}
		fn(res)
	}()
}
