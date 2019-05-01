package http

import (
	"strconv"
)

// ContentType for HTTP
type ContentType int

// Types of Content
const (
	TextPlain ContentType = iota //default

	TextHTML
)

// Content for HTTP
type Content struct {
	Type string
	Data []byte
	// only use in res cache
	//Len int
}

// NewContent for HTTP
func NewContent(t ContentType) *Content {
	T := "text/plain"
	switch t {
	case TextHTML:
		T = "text/html"
	}
	return &Content{
		Type: T,
		Data: nil,
	}
}

// TypeHeader for Content
func (c *Content) TypeHeader() string {
	return "Content-Type: " + c.Type
}

// LengthHeader for Content
func (c *Content) LengthHeader() string {
	return "Content-Length: " + strconv.Itoa(len(c.Data))
}
