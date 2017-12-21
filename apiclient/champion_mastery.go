package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/region"
)

type ChampionMastery struct {
	ChestGranted                 bool              `datastore:",noindex"` // Is chest granted for this champion or not in current season.
	ChampionLevel                int               `datastore:",noindex"` // Champion level for specified player and champion combination.
	ChampionPoints               int               `datastore:",noindex"` // Total number of champion points for this player and champion combination - they are used to determine championLevel.
	ChampionID                   champion.Champion `datastore:",noindex"` // Champion ID for this entry.
	PlayerID                     int64             `datastore:",noindex"` // Player ID for this entry.
	ChampionPointsUntilNextLevel int64             `datastore:",noindex"` // Number of points needed to achieve next level. Zero if player reached maximum champion level for this champion.
	ChampionPointsSinceLastLevel int64             `datastore:",noindex"` // Number of points earned since current level has been achieved. Zero if player reached maximum champion level for this champion.
	LastPlayTime                 int64             `datastore:",noindex"` // Last time this champion was played by this player - in Unix milliseconds time format.
}

func (c *client) GetAllChampionMasteries(ctx context.Context, r region.Region, summonerID int64) ([]ChampionMastery, error) {
	var res []ChampionMastery
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/champion-mastery/v3/champion-masteries/by-summoner", fmt.Sprintf("/%d", summonerID), nil, &res)
	return res, err
}

func (c *client) GetChampionMastery(ctx context.Context, r region.Region, summonerID int64, champ champion.Champion) (*ChampionMastery, error) {
	var res ChampionMastery
	_, err := c.dispatchAndUnmarshalWithUniquifier(ctx, r, "/lol/champion-mastery/v3/champion-masteries/by-summoner", fmt.Sprintf("/%d/by-champion/%d", summonerID, champ), nil, "by-champion", &res)
	return &res, err
}

func (c *client) GetChampionMasteryScore(ctx context.Context, r region.Region, summonerID int64) (int, error) {
	var res int
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/champion-mastery/v3/scores/by-summoner", fmt.Sprintf("/%d", summonerID), nil, &res)
	return res, err
}
