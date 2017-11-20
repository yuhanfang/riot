// Package staticdata defines a static data client for Data Dragons.
//
// Static data is not rate limited.
package staticdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yuhanfang/riot/constants/language"
	"github.com/yuhanfang/riot/external"
)

// Client requests static data from the data dragons endpoint. It is illegal to
// directly construct an instance. Use the New() constructor instead.
type Client struct {
	d external.Doer
}

// New returns a Client configured with the given HTTP client.
func New(d external.Doer) *Client {
	return &Client{
		d: d,
	}
}

// getJSON marshals the JSON response from the URL into the destination.
func (c *Client) getJSON(ctx context.Context, url string, dest interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	res, err := c.d.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP response: %v", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

// Version is a game version.
type Version string

// Versions returns a list of game versions sorted with the most recent game
// version first.
func (c *Client) Versions(ctx context.Context) ([]Version, error) {
	var versions []Version
	err := c.getJSON(ctx, "https://ddragon.leagueoflegends.com/api/versions.json", &versions)
	return versions, err
}

// Languages returns a list of supported languages.
func (c *Client) Languages(ctx context.Context) ([]language.Language, error) {
	var langs []language.Language
	err := c.getJSON(ctx, "https://ddragon.leagueoflegends.com/cdn/languages.json", &langs)
	return langs, err
}
