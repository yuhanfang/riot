package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/tier"
)

type LeagueList struct {
	LeagueID string       `datastore:",noindex"`
	Tier     tier.Tier    `datastore:",noindex"`
	Entries  []LeagueItem `datastore:",noindex"`
	Queue    queue.Queue  `datastore:",noindex"`
	Name     string       `datastore:",noindex"`
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

type LeaguePosition struct {
	Rank             string     `datastore:",noindex"`
	HotStreak        bool       `datastore:",noindex"`
	MiniSeries       MiniSeries `datastore:",noindex"`
	Wins             int        `datastore:",noindex"`
	Veteran          bool       `datastore:",noindex"`
	Losses           int        `datastore:",noindex"`
	FreshBlood       bool       `datastore:",noindex"`
	PlayerOrTeamName string     `datastore:",noindex"`
	Inactive         bool       `datastore:",noindex"`
	PlayerOrTeamID   string     `datastore:",noindex"`
	LeagueID         string     `datastore:",noindex"`
	Tier             tier.Tier  `datastore:",noindex"`
	LeaguePoints     int        `datastore:",noindex"`
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
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/leagues", fmt.Sprintf("/%s", leagueID), nil, &res)
	return &res, err
}

func (c *client) GetAllLeaguePositionsForSummoner(ctx context.Context, r region.Region, summonerID int64) ([]LeaguePosition, error) {
	var res []LeaguePosition
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/positions/by-summoner", fmt.Sprintf("/%d", summonerID), nil, &res)
	return res, err
}
