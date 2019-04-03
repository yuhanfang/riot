package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/tier"
)

type LeagueList struct {
	LeagueID string       `json:"leagueID",datastore:",noindex"`
	Tier     tier.Tier    `json:"tier",datastore:",noindex"`
	Entries  []LeagueItem `json:"entries",datastore:",noindex"`
	Queue    queue.Queue  `json:"queue",datastore:",noindex"`
	Name     string       `json:"name",datastore:",noindex"`
}

type LeagueItem struct {
	Rank         string     `json:"rank"`
	HotStreak    bool       `json:"hotStreak"`
	MiniSeries   MiniSeries `json:"miniSeries"`
	Wins         int        `json:"wins"`
	Veteran      bool       `json:"veteran"`
	Losses       int        `json:"losses"`
	FreshBlood   bool       `json:"freshBlood"`
	SummonerName string     `json:"summonerName"`
	Inactive     bool       `json:"inactive"`
	SummonerID   string     `json:"summonerID"`
	LeaguePoints int        `json:"leaguePoints"`
}

type MiniSeries struct {
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Target   int    `json:"target"`
	Progress string `json:"progress"`
}

type LeaguePosition struct {
	Rank             string     `json:"rank",datastore:",noindex"`
	QueueType        string     `json:"queueType",datastore:",noindex"`
	HotStreak        bool       `json:"hotStreak",datastore:",noindex"`
	MiniSeries       MiniSeries `json:"miniSeries",datastore:",noindex"`
	Wins             int        `json:"wins",datastore:",noindex"`
	Veteran          bool       `json:"veteran",datastore:",noindex"`
	Losses           int        `json:"losses",datastore:",noindex"`
	FreshBlood       bool       `json:"freshBlood",datastore:",noindex"`
	PlayerOrTeamName string     `json:"playerOrTeamName",datastore:",noindex"`
	Inactive         bool       `json:"inactive",datastore:",noindex"`
	PlayerOrTeamID   string     `json:"playerOrTeamID",datastore:",noindex"`
	LeagueID         string     `json:"leagueID",datastore:",noindex"`
	Tier             tier.Tier  `json:"tier",datastore:",noindex"`
	LeaguePoints     int        `json:"leaguePoints",datastore:",noindex"`
}

func (c *client) GetChallengerLeague(ctx context.Context, r region.Region, q queue.Queue) (*LeagueList, error) {
	var res LeagueList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v4/challengerleagues/by-queue", fmt.Sprintf("/%s", q.String()), nil, &res)
	return &res, err
}

func (c *client) GetGrandmasterLeague(ctx context.Context, r region.Region, q queue.Queue) (*LeagueList, error) {
	var res LeagueList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v4/grandmasterleagues/by-queue", fmt.Sprintf("/%s", q.String()), nil, &res)
	return &res, err
}

func (c *client) GetMasterLeague(ctx context.Context, r region.Region, q queue.Queue) (*LeagueList, error) {
	var res LeagueList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v4/masterleagues/by-queue", fmt.Sprintf("/%s", q.String()), nil, &res)
	return &res, err
}

func (c *client) GetLeagueByID(ctx context.Context, r region.Region, leagueID string) (*LeagueList, error) {
	var res LeagueList
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v4/leagues", fmt.Sprintf("/%s", leagueID), nil, &res)
	return &res, err
}

func (c *client) GetAllLeaguePositionsForSummoner(ctx context.Context, r region.Region, summonerID string) ([]LeaguePosition, error) {
	var res []LeaguePosition
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/positions/by-summoner", fmt.Sprintf("/%s", summonerID), nil, &res)
	return res, err
}
