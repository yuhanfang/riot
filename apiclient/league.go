package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/tier"
)

type LeagueList struct {
	LeagueID string
	Tier     tier.Tier
	Entries  []LeagueItem
	Queue    queue.Queue
	Name     string
}

type LeagueItem struct {
	Rank             string
	HotStreak        bool
	MiniSeries       MiniSeries
	Wins             int
	Veteran          bool
	Losses           int
	FreshBlood       bool
	PlayerOrTeamName string
	Inactive         bool
	PlayerOrTeamID   string
	LeaguePoints     int
}

type MiniSeries struct {
	Wins     int
	Losses   int
	Target   int
	Progress string
}

type LeaguePositions struct {
	Rank             string
	HotStreak        bool
	MiniSeries       MiniSeries
	Wins             int
	Veteran          bool
	Losses           int
	FreshBlood       bool
	PlayerOrTeamName string
	Inactive         bool
	PlayerOrTeamID   string
	LeagueID         string
	Tier             tier.Tier
	LeaguePoints     int
}

func (c *client) GetChallengerLeague(ctx context.Context, r region.Region, q queue.Queue) (*LeagueList, error) {
	var res LeagueList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/challengerleagues/by-queue", fmt.Sprintf("/%s", q.String()), nil, &res)
	return &res, err
}

func (c *client) GetMasterLeague(ctx context.Context, r region.Region, q queue.Queue) (*LeagueList, error) {
	var res LeagueList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/masterleagues/by-queue", fmt.Sprintf("/%s", q.String()), nil, &res)
	return &res, err
}

func (c *client) GetLeagueByID(ctx context.Context, r region.Region, leagueID string) (*LeagueList, error) {
	var res LeagueList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/leagues", fmt.Sprintf("/%d", leagueID), nil, &res)
	return &res, err
}

func (c *client) GetAllLeaguePositionsForSummoner(ctx context.Context, r region.Region, summonerID int64) (*LeaguePositions, error) {
	var res LeaguePositions
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/positions/by-summoner", fmt.Sprintf("/%d", summonerID), nil, &res)
	return &res, err
}
