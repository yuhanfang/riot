package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/region"
)

type Summoner struct {
	ProfileIconID int    `datastore:",noindex"` // ID of the summoner icon associated with the summoner.
	Name          string `datastore:",noindex"` // Summoner name.
	PUUID		  string `datastore:",noinex"`  // PUUID is the player universally unique identifier.
	SummonerLevel int64  `datastore:",noindex"` // Summoner level associated with the summoner.
	AccountID     string `datastore:",noindex"` // Encrypted account ID.
	ID            string `datastore:",noindex"` // Encrypted summoner ID.
	RevisionDate  int64  `datastore:",noindex"` // Date summoner was last modified specified as epoch milliseconds. The following events will update this timestamp: profile icon change, playing the tutorial or advanced tutorial, finishing a game, summoner name change
}

func (c *client) GetByAccountID(ctx context.Context, r region.Region, accountID string) (*Summoner, error) {
	var res Summoner
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/summoner/v4/summoners/by-account", fmt.Sprintf("/%s", accountID), nil, &res)
	return &res, err
}

func (c *client) GetBySummonerName(ctx context.Context, r region.Region, name string) (*Summoner, error) {
	var res Summoner
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/summoner/v4/summoners/by-name", fmt.Sprintf("/%s", name), nil, &res)
	return &res, err
}

func (c *client) GetBySummonerPUUID(ctx context.Context, r region.Region, puuid string) (*Summoner, error) {
	var res Summoner
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/summoner/v4/summoners/by-puuid", fmt.Sprintf("/%s", puuid), nil, &res)
	return &res, err
}

func (c *client) GetBySummonerID(ctx context.Context, r region.Region, summonerID string) (*Summoner, error) {
	var res Summoner
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/summoner/v4/summoners", fmt.Sprintf("/%s", summonerID), nil, &res)
	return &res, err
}
