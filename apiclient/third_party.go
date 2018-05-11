package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/region"
)

func (c *client) GetThirdPartyCodeByID(ctx context.Context, r region.Region, summonerID int64) (string, error) {
	var res string
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/platform/v3/third-party-code/by-summoner", fmt.Sprintf("/%d", summonerID), nil, &res)
	return res, err
}
