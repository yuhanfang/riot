package apiclient

import "errors"

var (
	ErrBadRequest           = errors.New("bad request")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrDataNotFound         = errors.New("data not found")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrInternalServerError  = errors.New("internal server error")
	ErrBadGateway           = errors.New("bad gateway")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrGatewayTimeout       = errors.New("gateway timeout")

	ErrBadHTTPStatus = errors.New("bad HTTP status returned by server")

	httpErrors = map[int]error{
		400: ErrBadRequest,
		401: ErrUnauthorized,
		403: ErrForbidden,
		404: ErrDataNotFound,
		405: ErrMethodNotAllowed,
		415: ErrUnsupportedMediaType,
		429: ErrRateLimitExceeded,
		500: ErrInternalServerError,
		502: ErrBadGateway,
		503: ErrServiceUnavailable,
		504: ErrGatewayTimeout,
	}
)
