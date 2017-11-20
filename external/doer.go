package external

import "net/http"

// Doer executes arbitrary HTTP requests, ans is usually a *http.Client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
