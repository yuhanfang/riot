package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/region"
)

type Summoner struct {
	ProfileIconID int    `datastore:",noindex"` // ID of the summoner icon associated with the summoner.
	Name          string `datastore:",noindex"` //Summoner name.
	SummonerLevel int64  `datastore:",noindex"` // Summoner level associated with the summoner.
	RevisionDate  int64  `datastore:",noindex"` // Date summoner was last modified specified as epoch milliseconds. The following events will update this timestamp: profile icon change, playing the tutorial or advanced tutorial, finishing a game, summoner name change
	ID            int64  `datastore:",noindex"` // Summoner ID.
	AccountID     int64  `datastore:",noindex"` //Account ID.
}

func (c *client) GetByAccountID(ctx context.Context, r region.Region, accountID int64) (*Summoner, error) {
	var res Summoner
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/summoner/v3/summoners/by-account", fmt.Sprintf("/%d", accountID), nil, &res)
	return &res, err
}

func (c *client) GetBySummonerName(ctx context.Context, r region.Region, name string) (*Summoner, error) {
	var res Summoner
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/summoner/v3/summoners/by-name", fmt.Sprintf("/%s", name), nil, &res)
	return &res, err
}

func (c *client) GetBySummonerID(ctx context.Context, r region.Region, summonerID int64) (*Summoner, error) {
	var res Summoner
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/summoner/v3/summoners", fmt.Sprintf("/%d", summonerID), nil, &res)
	return &res, err
}
