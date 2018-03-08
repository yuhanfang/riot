// Package esports implements an API for interacting with lolesports.
//
// Use the NewClient() constructor to construct a client, and call the client
// methods to interact with lolesports. Rate limiting is not implemented.
package esports

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yuhanfang/riot/external"
)

type Client struct {
	d   external.Doer
	log bool
}

func NewClient(doer external.Doer) *Client {
	return &Client{
		d: doer,
	}
}

func (c *Client) EnableLogs() {
	c.log = true
}

func (c *Client) DisableLogs() {
	c.log = false
}

func (c Client) doJSON(ctx context.Context, req *http.Request, dest interface{}) (*http.Response, error) {
	res, err := c.d.Do(req.WithContext(ctx))
	if err != nil {
		return res, err
	}
	if dest == nil {
		return res, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewReader(body))
	if err != nil {
		return res, err
	}
	if c.log {
		log.Printf("%s %s = %s", req.Method, req.URL, string(body))
	}
	return res, json.Unmarshal(body, dest)
}
