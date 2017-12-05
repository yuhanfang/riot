// Package apiclient accesses the official Riot API.
//
// Construct a client with the New() function, and call the various client
// methods to retrieve data from the API.
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

	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/external"
	"github.com/yuhanfang/riot/ratelimit"
)

// Client accesses the Riot API. Use New() to retrieve a valid instance.
type Client interface {
	// ----- Champion Mastery API -----

	// GetAllChampionMasteries returns all champion mastery entries sorted by
	// number of champion points descending.
	GetAllChampionMasteries(ctx context.Context, r region.Region, summonerID int64) ([]ChampionMastery, error)

	// GetChampionMastery returns champion mastery by summoner ID and champion.
	GetChampionMastery(ctx context.Context, r region.Region, summonerID int64, champ champion.Champion) (*ChampionMastery, error)

	// GetChampionMasteryScore returns a player's total champion mastery score,
	// which is the sum of individual champion mastery levels.
	GetChampionMasteryScore(ctx context.Context, r region.Region, summonerID int64) (int, error)

	// ----- Champions API -----

	// GetChampions returns all champions.
	GetChampions(ctx context.Context, r region.Region) (*ChampionList, error)

	// GetChampionByID returns champion information for a specific champion.
	GetChampionByID(ctx context.Context, r region.Region, champ champion.Champion) (*Champion, error)

	// ----- League API -----

	// GetChallengerLeague returns the challenger league for the given queue.
	GetChallengerLeague(context.Context, region.Region, queue.Queue) (*LeagueList, error)

	// GetMasterLeague returns the master league for the given queue.
	GetMasterLeague(context.Context, region.Region, queue.Queue) (*LeagueList, error)

	// GetAllLeaguePositionsForSummoner returns league positions in all queues
	// for the given summoner ID.
	GetAllLeaguePositionsForSummoner(ctx context.Context, r region.Region, summonerID int64) ([]LeaguePosition, error)

	// GetLeagueByID returns the league with given ID, including inactive
	// entries.
	GetLeagueByID(ctx context.Context, r region.Region, leagueID string) (*LeagueList, error)

	// ----- Match API -----

	// GetMatch returns a match by match ID.
	GetMatch(ctx context.Context, r region.Region, matchID int64) (*Match, error)

	GetMatchTimeline(ctx context.Context, r region.Region, matchID int64) (*MatchTimeline, error)

	// GetMatchlist returns a matchlist for games played on a given account ID
	// and filtered using given filter parameters, if any.
	GetMatchlist(ctx context.Context, r region.Region, accountID int64, opts *GetMatchlistOptions) (*Matchlist, error)

	// GetRecentMatchlist returns the last 20 matches played on the given account ID.
	GetRecentMatchlist(ctx context.Context, r region.Region, accountID int64) (*Matchlist, error)

	// ----- Spectator API -----

	// GetFeaturedGames returns a list of featured games.
	GetFeaturedGames(ctx context.Context, r region.Region) (*FeaturedGames, error)

	// GetCurrentGameInfoBySummoner returns current game information for a given
	// summoner ID.
	GetCurrentGameInfoBySummoner(ctx context.Context, r region.Region, summonerID int64) (*CurrentGameInfo, error)

	// ----- Summoner API -----

	// GetByAccountID returns a summoner by account ID.
	GetByAccountID(ctx context.Context, r region.Region, accountID int64) (*Summoner, error)

	// GetBySummonerName returns a summoner by summoner name.
	GetBySummonerName(ctx context.Context, r region.Region, name string) (*Summoner, error)

	// GetBySummonerID returns a sumoner by summoner ID.
	GetBySummonerID(ctx context.Context, r region.Region, summonerID int64) (*Summoner, error)
}

// client is the internal implementation of Client.
type client struct {
	key string
	c   external.Doer
	r   ratelimit.Limiter
}

// New returns a Client configured for the given API client and underlying HTTP
// client. The returned Client is threadsafe.
func New(key string, httpClient external.Doer, limiter ratelimit.Limiter) Client {
	return &client{
		key: key,
		c:   httpClient,
		r:   limiter,
	}
}

// dispatchAndUnmarshalWithUniquifier is the same as dispatchAndUnmarshal,
// except with an additional uniquifier parameter that allows special case
// handling of certain methods that have different quota buckets depending on
// the relative path.
func (c *client) dispatchAndUnmarshalWithUniquifier(ctx context.Context, r region.Region, m string, relativePath string, v url.Values, u string, dest interface{}) (*http.Response, error) {
	res, err := c.dispatchMethod(ctx, r, m, relativePath, v, u)
	if err != nil {
		return res, err
	}
	if res.StatusCode != http.StatusOK {
		err, ok := httpErrors[res.StatusCode]
		if !ok {
			err = ErrBadHTTPStatus
		}
		return res, err
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

// dispatchAndUnmarshal dispatches the method (see dispatchMethod). If the
// method returns HTTP okay, then read the body into a buffer and attempt to
// unmarshal it into the supplied destination. Otherwise, the method returns
// one of the documented errors. In any case, the body is set to read from the
// beginning of the stream and is left open, as if the response were returned
// directly from an HTTP request.
func (c *client) dispatchAndUnmarshal(ctx context.Context, r region.Region, m string, relativePath string, v url.Values, dest interface{}) (*http.Response, error) {
	return c.dispatchAndUnmarshalWithUniquifier(ctx, r, m, relativePath, v, "", dest)
}

// dispatchMethod calls the given API method for the given region. The
// relativePath is appended to the method to form the REST endpoint. The given
// URL values are encoded and passed as URL parameters following the REST
// endpoint.
func (c *client) dispatchMethod(ctx context.Context, r region.Region, m string, relativePath string, v url.Values, uniquifier string) (*http.Response, error) {
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
		Region:         strings.ToUpper(string(r)),
		Method:         strings.ToLower(m),
		Uniquifier:     uniquifier,
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
