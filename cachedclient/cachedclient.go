// Package cachedclient implements a cached version of the Riot API. The client
// is backed by a Datastore that persists RPC results, and a Client that calls
// to Riot in the event an RPC is not available in the Cache.
//
// Use the New() constructor to initialize a Client.
package cachedclient

import (
	"context"
	"fmt"
	"time"

	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/season"
)

var zeroTime time.Time

type client struct {
	apiclient.Client

	d Datastore
}

// Datastore is a key-time-value store used to cache values.
type Datastore interface {
	// Get returns the value for the given key and the time that the entry was
	// written. If a non-zero time is provided, then Get returns the most recent
	// value before the given time.
	Get(ctx context.Context, key string, dest interface{}, t time.Time) (time.Time, error)

	// Put stores the value for the given key. If a non-zero time is provided,
	// then Put stores the entry by timestamp instead of overwriting a single
	// shared key.
	Put(ctx context.Context, key string, val interface{}, t time.Time) error

	// Purge removes all values for the given key except for the given number of
	// most recent entries which should be kept. For example, if keep is one, then
	// only the most recent value for the key will be kept, and the rest will be
	// deleted.
	Purge(ctx context.Context, key string, keep int) error
}

func (c *client) GetChampions(ctx context.Context, r region.Region) (*apiclient.ChampionList, error) {
	var val apiclient.ChampionList
	key := fmt.Sprintf("get-champions:%s", r)
	t, err := c.d.Get(ctx, key, &val, zeroTime)
	if err == nil && time.Since(t) < 24*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetChampions(ctx, r)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	go c.d.Purge(ctx, key, 1)
	return res, err
}

func (c *client) GetChampionsByID(ctx context.Context, r region.Region, champ champion.Champion) (*apiclient.Champion, error) {
	var val apiclient.Champion
	key := fmt.Sprintf("get-champion-by-id:%s:%d", r, champ)
	t, err := c.d.Get(ctx, key, &val, zeroTime)
	if err == nil && time.Since(t) < 24*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetChampionByID(ctx, r, champ)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	go c.d.Purge(ctx, key, 1)
	return res, err
}

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

func (c *client) GetMasterLeague(ctx context.Context, r region.Region, q queue.Queue) (*apiclient.LeagueList, error) {
	var val apiclient.LeagueList
	key := fmt.Sprintf("get-master-league:%s:%s", r, q)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 24*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetMasterLeague(ctx, r, q)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	return res, err
}

func (c *client) GetAllLeaguePositionsForSummoner(ctx context.Context, r region.Region, summonerID string) ([]apiclient.LeaguePosition, error) {
	type LeaguePositions struct {
		Positions []apiclient.LeaguePosition
	}
	var val LeaguePositions
	key := fmt.Sprintf("get-all-league-positions-for-summoner:%s:%s", r, summonerID)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 24*time.Hour {
		return val.Positions, nil
	}
	res, err := c.Client.GetAllLeaguePositionsForSummoner(ctx, r, summonerID)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, &LeaguePositions{res}, time.Now())
	go c.d.Purge(ctx, key, 1)
	return res, err
}

func (c *client) GetLeagueByID(ctx context.Context, r region.Region, leagueID string) (*apiclient.LeagueList, error) {
	var val apiclient.LeagueList
	key := fmt.Sprintf("get-league-by-id:%s:%s", r, leagueID)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 24*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetLeagueByID(ctx, r, leagueID)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	return res, err
}

func (c *client) GetMatch(ctx context.Context, r region.Region, matchID int64) (*apiclient.Match, error) {
	return c.Client.GetMatch(ctx, r, matchID)
}

func (c *client) GetMatchTimeline(ctx context.Context, r region.Region, matchID int64) (*apiclient.MatchTimeline, error) {
	return c.Client.GetMatchTimeline(ctx, r, matchID)
}

// filterMatchlist applies the given options to the matchlist, returning a new
// match list that is filtered.
func filterMatchlist(m *apiclient.Matchlist, opts *apiclient.GetMatchlistOptions) *apiclient.Matchlist {
	if opts == nil {
		return m
	}
	var res apiclient.Matchlist

	queues := make(map[queue.Queue]bool)
	for _, q := range opts.Queue {
		queues[q] = true
	}
	seasons := make(map[season.Season]bool)
	for _, s := range opts.Season {
		seasons[s] = true
	}
	champions := make(map[champion.Champion]bool)
	for _, c := range opts.Champion {
		champions[c] = true
	}

	for i, match := range m.Matches {
		if len(opts.Queue) > 0 && !queues[match.Queue] {
			continue
		}
		if len(opts.Season) > 0 && !seasons[match.Season] {
			continue
		}
		if len(opts.Champion) > 0 && !champions[match.Champion] {
			continue
		}

		ts := match.Timestamp.Time()
		if opts.BeginTime != nil && opts.EndTime != nil {
			if opts.BeginTime.After(ts) || opts.EndTime.Before(ts) {
				continue
			}
		} else if opts.EndTime != nil {
			if opts.EndTime.Before(ts) {
				continue
			}
		} else if opts.BeginTime != nil {
			if opts.BeginTime.After(ts) {
				continue
			}
		}

		if opts.BeginIndex == nil && opts.EndIndex != nil {
			beginIndex := 0
			opts.BeginIndex = &beginIndex
		} else if opts.BeginIndex != nil && opts.EndIndex == nil {
			endIndex := *opts.BeginIndex + 100
			opts.EndIndex = &endIndex
		}

		if opts.BeginIndex != nil && opts.EndIndex != nil {
			if i < *opts.BeginIndex || i >= *opts.EndIndex {
				continue
			}
		}
		res.Matches = append(res.Matches, match)
	}

	res.TotalGames = len(res.Matches)
	if opts.BeginIndex != nil {
		res.StartIndex = *opts.BeginIndex
	}
	if opts.EndIndex != nil {
		res.EndIndex = *opts.EndIndex
	}
	return &res
}

func (c *client) GetMatchlist(ctx context.Context, r region.Region, accountID string, opt *apiclient.GetMatchlistOptions) (*apiclient.Matchlist, error) {
	var val apiclient.Matchlist
	key := fmt.Sprintf("get-matchlist:%s:%d", r, accountID)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 24*time.Hour {
		return filterMatchlist(&val, opt), nil
	}
	res, err := c.Client.GetMatchlist(ctx, r, accountID, nil)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	go c.d.Purge(ctx, key, 1)
	return filterMatchlist(res, opt), err
}

func (c *client) GetRecentMatchlist(ctx context.Context, r region.Region, accountID string) (*apiclient.Matchlist, error) {
	var val apiclient.Matchlist
	key := fmt.Sprintf("get-recent-matchlist:%s:%s", r, accountID)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 1*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetRecentMatchlist(ctx, r, accountID)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	go c.d.Purge(ctx, key, 1)
	return res, err
}

func (c *client) GetFeaturedGames(ctx context.Context, r region.Region) (*apiclient.FeaturedGames, error) {
	var val apiclient.FeaturedGames
	key := fmt.Sprintf("get-featured-games:%s", r)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 1*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetFeaturedGames(ctx, r)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	go c.d.Purge(ctx, key, 1)
	return res, err
}

func (c *client) GetByAccountID(ctx context.Context, r region.Region, accountID string) (*apiclient.Summoner, error) {
	var val apiclient.Summoner
	key := fmt.Sprintf("get-by-account-id:%s:%d", r, accountID)
	_, err := c.d.Get(ctx, key, &val, zeroTime)
	if err == nil {
		return &val, nil
	}
	res, err := c.Client.GetByAccountID(ctx, r, accountID)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, zeroTime)
	return res, err
}

func (c *client) GetBySummonerName(ctx context.Context, r region.Region, name string) (*apiclient.Summoner, error) {
	var val apiclient.Summoner
	key := fmt.Sprintf("get-by-summoner-name:%s:%s", r, name)
	t, err := c.d.Get(ctx, key, &val, time.Now())
	if err == nil && time.Since(t) < 24*time.Hour {
		return &val, nil
	}
	res, err := c.Client.GetBySummonerName(ctx, r, name)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, time.Now())
	go c.d.Purge(ctx, key, 1)
	return res, err
}

func (c *client) GetBySummonerID(ctx context.Context, r region.Region, summonerID string) (*apiclient.Summoner, error) {
	var val apiclient.Summoner
	key := fmt.Sprintf("get-by-summoner-id:%s:%d", r, summonerID)
	_, err := c.d.Get(ctx, key, &val, zeroTime)
	if err == nil {
		return &val, nil
	}
	res, err := c.Client.GetBySummonerID(ctx, r, summonerID)
	if err != nil {
		return nil, err
	}
	err = c.d.Put(ctx, key, res, zeroTime)
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
