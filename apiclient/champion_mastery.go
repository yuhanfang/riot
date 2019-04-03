package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/region"
)

type ChampionMastery struct {
	ChestGranted                 bool              `json:"chestGranted",datastore:",noindex"`                 // Is chest granted for this champion or not in current season.
	ChampionLevel                int               `json:"championLevel",datastore:",noindex"`                // Champion level for specified player and champion combination.
	ChampionPoints               int               `json:"championPoints",datastore:",noindex"`               // Total number of champion points for this player and champion combination - they are used to determine championLevel.
	ChampionID                   champion.Champion `json:"championID",datastore:",noindex"`                   // Champion ID for this entry.
	PlayerID                     string            `json:"playerID",datastore:",noindex"`                     // Encrypted Player ID for this entry.
	ChampionPointsUntilNextLevel int64             `json:"championPointsUntilNextLevel",datastore:",noindex"` // Number of points needed to achieve next level. Zero if player reached maximum champion level for this champion.
	ChampionPointsSinceLastLevel int64             `json:"championPointsSinceLastLevel",datastore:",noindex"` // Number of points earned since current level has been achieved. Zero if player reached maximum champion level for this champion.
	LastPlayTime                 int64             `json:"lastPlayTime",datastore:",noindex"`                 // Last time this champion was played by this player - in Unix milliseconds time format.
}

func (c *client) GetAllChampionMasteries(ctx context.Context, r region.Region, summonerID string) ([]ChampionMastery, error) {
	var res []ChampionMastery
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/champion-mastery/v4/champion-masteries/by-summoner", fmt.Sprintf("/%s", summonerID), nil, &res)
	return res, err
}

func (c *client) GetChampionMastery(ctx context.Context, r region.Region, summonerID string, champ champion.Champion) (*ChampionMastery, error) {
	var res ChampionMastery
	_, err := c.dispatchAndUnmarshalWithUniquifier(ctx, r, "/lol/champion-mastery/v4/champion-masteries/by-summoner", fmt.Sprintf("/%s/by-champion/%d", summonerID, champ), nil, "by-champion", &res)
	return &res, err
}

func (c *client) GetChampionMasteryScore(ctx context.Context, r region.Region, summonerID string) (int, error) {
	var res int
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/champion-mastery/v4/scores/by-summoner", fmt.Sprintf("/%s", summonerID), nil, &res)
	return res, err
}
