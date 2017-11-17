package cachedclient

import (
	"context"
	"fmt"
	"time"

	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
)

type client struct {
	apiclient.Client

	d Datastore
}

type Datastore interface {
	// Get returns the value for the given key and the time that the entry was
	// written. If a non-zero time is provided, then Get returns the most recent
	// value before the given time.
	Get(ctx context.Context, key string, dest interface{}, t time.Time) (time.Time, error)

	// Put stores the value for the given key. If a non-zero time is provided,
	// then Put stores the entry by timestamp instead of overwriting a single
	// shared key.
	Put(ctx context.Context, key string, val interface{}, t time.Time) error
}

// GetChallengerLeague returns a cached call from within 24 hours ago. If no
// cached call exists, then the function makes a call and caches it.
func (c *client) GetChallengerLeague(ctx context.Context, r region.Region, q queue.Queue) (*apiclient.LeagueList, error) {
	var val apiclient.LeagueList
	key := fmt.Sprintf("get-challenger-league:%s:%s", r, q)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 24*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetChallengerLeague(ctx, r, q)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	return res, err
}

// New returns a cached client, using the underlying client to query non-cached
// values, and the underlying datastore as the cache location.
func New(c apiclient.Client, d Datastore) apiclient.Client {
	return &client{
		Client: c,
		d:      d,
	}
}
