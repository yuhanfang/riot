package apiclient

import (
	"fmt"

	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/tier"
	"golang.org/x/net/context"
)

type LeagueListDTO struct {
	LeagueID string
	Tier     tier.Tier
	Entries  []LeagueItemDTO
	Queue    queue.Queue
	Name     string
}

type LeagueItemDTO struct {
	Rank             string
	HotStreak        bool
	MiniSeries       MiniSeriesDTO
	Wins             int
	Veteran          bool
	Losses           int
	FreshBlood       bool
	PlayerOrTeamName string
	Inactive         bool
	PlayerOrTeamID   string
	LeaguePoints     int
}

type MiniSeriesDTO struct {
	Wins     int
	Losses   int
	Target   int
	Progress string
}

type LeaguePositionDTO struct {
	Rank             string
	HotStreak        bool
	MiniSeries       MiniSeriesDTO
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

func (c *client) GetChallengerLeague(ctx context.Context, r region.Region, q queue.Queue) (*LeagueListDTO, error) {
	var res LeagueListDTO
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/challengerleagues/by-queue", fmt.Sprintf("/%s", q.String()), nil, &res)
	return &res, err
}

func (c *client) GetMasterLeague(ctx context.Context, r region.Region, q queue.Queue) (*LeagueListDTO, error) {
	var res LeagueListDTO
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/masterleagues/by-queue", fmt.Sprintf("/%s", q.String()), nil, &res)
	return &res, err
}

func (c *client) GetLeagueByID(ctx context.Context, r region.Region, leagueID string) (*LeagueListDTO, error) {
	var res LeagueListDTO
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/leagues", fmt.Sprintf("/%d", leagueID), nil, &res)
	return &res, err
}

func (c *client) GetAllLeaguePositionsForSummoner(ctx context.Context, r region.Region, summonerID int64) (*LeaguePositionDTO, error) {
	var res LeaguePositionDTO
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/league/v3/positions/by-summoner", fmt.Sprintf("/%d", summonerID), nil, &res)
	return &res, err
}
