package http

type StatusCode int

const (
	Continue StatusCode = 100 + iota
	SwitchingProtocols
	Processing
)
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
const (
	UnparseableResponseHeaders StatusCode = 600 + iota
)
