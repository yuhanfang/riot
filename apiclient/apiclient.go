package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

// Client accesses the Riot API. Use New() to retrieve a valid instance.
type Client interface {
	// League API

	GetChallengerLeague(context.Context, region.Region, queue.Queue) (*LeagueListDTO, error)
	GetMasterLeague(context.Context, region.Region, queue.Queue) (*LeagueListDTO, error)
	GetAllLeaguePositionsForSummoner(context.Context, region.Region, int64) (*LeaguePositionDTO, error)

	// Match API

	GetMatch(context.Context, region.Region, int64) (*MatchDTO, error)
	GetMatchlist(context.Context, region.Region, int64, *GetMatchlistOptions) (*MatchlistDTO, error)
	GetRecentMatchlist(context.Context, region.Region, int64) (*MatchlistDTO, error)
}

// client is the internal implementation of Client.
type client struct {
	key string
	c   Doer
	r   ratelimit.Limiter
}

// Doer executes arbitrary HTTP requests, and is usually a *http.Client.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// New returns a Client configured for the given API client and underlying HTTP
// client. The returned Client is threadsafe.
func New(key string, httpClient Doer, limiter ratelimit.Limiter) Client {
	return &client{
		key: key,
		c:   httpClient,
		r:   limiter,
	}
}

// dispatchAndUnmarshal dispatches the method (see dispatchMethod). If the
// method returns HTTP okay, then read the body into a buffer and attempt to
// unmarshal it into the supplied destination. Otherwise, the method returns
// ErrBadHTTPStatus. In any case, the body is set to read from the beginning of
// the stream and is left open, as if the response were returned directly from
// an HTTP request.
func (c *client) dispatchAndUnmarshal(ctx context.Context, r region.Region, m string, relativePath string, v url.Values, dest interface{}) (*http.Response, error) {
	res, err := c.dispatchMethod(ctx, r, m, relativePath, v)
	if err != nil {
		return res, err
	}
	if res.StatusCode != http.StatusOK {
		return res, ErrBadHTTPStatus
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewReader(b))

	// The body is in good state, so now we can return if there was an IO problem.
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, dest)

	return res, err
}

// dispatchMethod calls the given API method for the given region. The
// relativePath is appended to the method to form the REST endpoint. The given
// URL values are encoded and passed as URL parameters following the REST
// endpoint.
func (c *client) dispatchMethod(ctx context.Context, r region.Region, m string, relativePath string, v url.Values) (*http.Response, error) {
	var suffix, separator string

	if len(v) > 0 {
		suffix = fmt.Sprintf("?%s", v.Encode())
	}
	if !strings.HasPrefix(relativePath, "/") {
		separator = "/"
	}
	path := r.Host() + m + separator + relativePath + suffix
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Add("X-Riot-Token", c.key)

	done, _, err := c.r.Acquire(ctx, ratelimit.Invocation{
		ApplicationKey: c.key,
		Region:         string(r),
		Method:         m,
	})

	if err != nil {
		return nil, err
	}

	// If either the done() or the HTTP request is an error, then return error.
	res, err := c.c.Do(req)
	derr := done(res)
	if err == nil {
		err = derr
	}
	return res, err
}
