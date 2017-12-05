// Package client implements rate limiting by connecting to a centralized rate
// limiting server.
package client

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/yuhanfang/riot/external"
	"github.com/yuhanfang/riot/ratelimit"
)

// client implemnts the ratelimit.Limiter interface by querying a rate limit
// server.
type client struct {
	base *url.URL
	d    external.Doer
}

// Acquire acquires quota for the given invocation. The caller must call done()
// or cancel() within one minute of a successful call, or the quota will be
// assumed to have been used, and will refresh after the maximum time.
func (c *client) Acquire(ctx context.Context, inv ratelimit.Invocation) (ratelimit.Done, ratelimit.Cancel, error) {
	address := c.base.String() + "/acquire/" + inv.ApplicationKey + "/" + inv.Region
	values := url.Values(make(map[string][]string))
	if inv.Method != "" {
		values.Add("method", inv.Method)
	}
	if inv.Uniquifier != "" {
		values.Add("uniquifier", inv.Uniquifier)
	}
	if inv.NoAppQuota {
		values.Add("noappquota", "T")
	}
	req, err := http.NewRequest("POST", address, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)
	res, err := c.d.Do(req)
	err = getError(res, err)
	if err != nil {
		return nil, nil, err
	}
	tok, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	token := string(tok)

	done := func(res *http.Response) error {
		address := c.base.String() + "/done/" + token
		req, err := http.NewRequest("POST", address, nil)
		if err != nil {
			return err
		}
		if res != nil {
			req.Header = res.Header
		}
		req = req.WithContext(ctx)
		res, err = c.d.Do(req)
		err = getError(res, err)
		return err
	}

	cancel := func() error {
		address := c.base.String() + "/cancel/" + token
		req, err := http.NewRequest("POST", address, nil)
		if err != nil {
			return err
		}
		req = req.WithContext(ctx)
		res, err = c.d.Do(req)
		err = getError(res, err)
		return err
	}

	return done, cancel, nil
}

// getError returns the error on bad response or if err is non-nil.
func getError(res *http.Response, err error) error {
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(b))
	}
	return nil
}

// New returns a Limiter configured with the given http client (usually
// http.DefaultClient) and base URL of the server.
func New(doer external.Doer, base *url.URL) ratelimit.Limiter {
	return &client{
		d:    doer,
		base: base,
	}
}
