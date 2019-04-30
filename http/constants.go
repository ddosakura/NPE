package http

// Version for HTTP
type Version int

// Versions
const (
	HTTP10 Version = iota
	HTTP11
	HTTP20
)

// SetVersion for Request
func (r *Request) SetVersion(v Version) {
	switch v {
	case HTTP10:
		r.Version = "HTTP/1.0"
	case HTTP11:
		r.Version = "HTTP/1.1"
	case HTTP20:
		r.Version = "HTTP/2.0"
	}
}

// Method for HTTP
type Method int

// Method in HTTP/1.1 [RFC 7231]
const (
	GET Method = iota
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
)

// SetMethod for Request
func (r *Request) SetMethod(m Method) {
	switch m {
	case GET:
		r.Method = "GET"
	case POST:
		r.Method = "POST"
	case PUT:
		r.Method = "PUT"
	case DELETE:
		r.Method = "DELETE"
	case CONNECT:
		r.Method = "CONNECT"
	case OPTIONS:
		r.Method = "OPTIONS"
	case TRACE:
		r.Method = "TRACE"
	}
}

// StatusCode for HTTP
type StatusCode int

// Message
const (
	Continue StatusCode = 100 + iota
	SwitchingProtocols
	Processing
)

// Success
const (
	OK StatusCode = 200 + iota
	Created
	Accepted
	NonAuthoritativeInformation
	NoContent
	ResetContent
	PartialContent
	MultiStatus
)

// Redirect
const (
	MultipleChoices StatusCode = 300 + iota
	MovedPermanently
	MoveTemporarily
	SeeOther
	NotModified
	UseProxy
	SwitchProxy
	TemporaryRedirect
)

// Request Error
const (
	BadRequest StatusCode = 400 + iota
	Unauthorized
	PaymentRequired
	Forbidden
	NotFound
	MethodNotAllowed
	NotAcceptable
	ProxyAuthenticationRequired
	RequestTimeout
	Conflict
	Gone
	LengthRequired
	PreconditionFailed
	RequestEntityTooLarge
	RequestURITooLong
	UnsupportedMediaType
	RequestedRangeNotSatisfiable
	ExpectationFailed
	Teapot // I'm a teapot
	TooManyConnections
	UnprocessableEntity
	Locked
	FailedDependency
	UnorderedCollection
	UpgradeRequired
	RetryWith
	UnavailableForLegalReasons
)

// Server Error
const (
	InternalServerError StatusCode = 500 + iota
	NotImplemented
	BadGateway
	ServiceUnavailable
	GatewayTimeout
	HTTPVersionNotSupported
	VariantAlsoNegotiates
	InsufficientStorage
	BandwidthLimitExceeded
	NotExtended
)

// Server Error
const (
	UnparseableResponseHeaders StatusCode = 600 + iota
)
