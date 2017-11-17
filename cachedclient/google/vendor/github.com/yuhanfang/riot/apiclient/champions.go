package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/region"
)

type ChampionList struct {
	champions []Champion // The collection of champion information.
}

type Champion struct {
	RankedPlayEnabled bool  // Ranked play enabled flag.
	BotEnabled        bool  // Bot enabled flag (for custom games).
	BotMmEnabled      bool  // Bot Match Made enabled flag (for Co-op vs. AI games).
	Active            bool  // Indicates if the champion is active.
	FreeToPlay        bool  // Indicates if the champion is free to play. Free to play champions are rotated periodically.
	ID                int64 // Champion ID. For static information correlating to champion IDs, please refer to the LoL Static Data API.
}

func (c *client) GetChampions(ctx context.Context, r region.Region) (*ChampionList, error) {
	var res ChampionList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/platform/v3/champions", "", nil, &res)
	return &res, err
}

func (c *client) GetChampionByID(ctx context.Context, r region.Region, champ champion.Champion) (*Champion, error) {
	var res Champion
	_, err := c.dispatchAndUnmarshalWithUniquifier(ctx, r, "/lol/platform/v3/champions", fmt.Sprintf("/%d", champ), nil, "by-id", &res)
	return &res, err
}
