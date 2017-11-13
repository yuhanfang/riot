package apiclient

import "errors"

var (
	ErrBadHTTPStatus = errors.New("bad HTTP status returned by server")
)
