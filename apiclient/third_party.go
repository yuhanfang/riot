package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/region"
)

func (c *client) GetThirdPartyCodeByID(ctx context.Context, r region.Region, summonerID string) (string, error) {
	var res string
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/platform/v4/third-party-code/by-summoner", fmt.Sprintf("/%s", summonerID), nil, &res)
	return res, err
}
