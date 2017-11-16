package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/yuhanfang/riot/ratelimit"
	"github.com/yuhanfang/riot/ratelimit/service/client"
	"github.com/yuhanfang/riot/ratelimit/service/server"
)

func TestEndToEnd(t *testing.T) {
	s := server.New()
	ts := httptest.NewServer(s)
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := client.New(http.DefaultClient, u)

	ctx := context.Background()
	done, cancel, err := c.Acquire(ctx, ratelimit.Invocation{
		ApplicationKey: "key",
		Region:         "NA1",
		Method:         "/foo/bar",
		Uniquifier:     "unique",
		NoAppQuota:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = done(nil)
	if err != nil {
		t.Fatal(err)
	}
	err = cancel()
	if err == nil {
		t.Fatal("cancel should fail after done")
	}

	done, cancel, err = c.Acquire(ctx, ratelimit.Invocation{
		ApplicationKey: "key",
		Region:         "NA1",
		Method:         "/foo/bar",
		Uniquifier:     "unique",
		NoAppQuota:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = cancel()
	if err != nil {
		t.Fatal(err)
	}
	err = done(nil)
	if err == nil {
		t.Fatal("done should fail after cancel")
	}
}
