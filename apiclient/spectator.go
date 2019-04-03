package apiclient

import (
	"context"
	"fmt"

	"github.com/yuhanfang/riot/constants/region"
)

type CurrentGameInfo struct {
	GameID            int64                    `json:"gameID",datastore:",noindex"`            // The ID of the game
	GameStartTime     int64                    `json:"gameStartTime",datastore:",noindex"`     // The game start time represented in epoch milliseconds
	PlatformID        string                   `json:"platformID",datastore:",noindex"`        // The ID of the platform on which the game is being played
	GameMode          string                   `json:"gameMode",datastore:",noindex"`          // The game mode
	MapID             int64                    `json:"mapID",datastore:",noindex"`             // The ID of the map
	GameType          string                   `json:"gameType",datastore:",noindex"`          // The game type
	BannedChampions   []BannedChampion         `json:"bannedChampions",datastore:",noindex"`   // Banned champion information
	Observers         Observer                 `json:"observers",datastore:",noindex"`         // The observer information
	Participants      []CurrentGameParticipant `json:"participants",datastore:",noindex"`      // The participant information
	GameLength        int64                    `json:"gameLength",datastore:",noindex"`        // The amount of time in seconds that has passed since the game started
	GameQueueConfigID int64                    `json:"gameQueueConfigID",datastore:",noindex"` // The queue type (queue types are documented on the Game Constants page)
}

type BannedChampion struct {
	PickTurn   int   `json:"pickTurn"`   // The turn during which the champion was banned
	ChampionID int64 `json:"championID"` // The ID of the banned champion
	TeamID     int64 `json:"teamID"`     // The ID of the team that banned the champion
}

type Observer struct {
	EncryptionKey string `json:"encryptionKey"` // Key used to decrypt the spectator grid game data for playback
}

type CurrentGameParticipant struct {
	ProfileIconId int64                              `json:"profileIconId"` // The ID of the profile icon used by this participant
	ChampionId    int64                              `json:"championId"`    // The ID of the champion played by this participant
	SummonerName  string                             `json:"summonerName"`  // The summoner name of this participant
	Runes         []CurrentGameParticipantRuneDTO    `json:"runes"`         // The runes used by this participant
	Bot           bool                               `json:"bot"`           // Flag indicating whether or not this participant is a bot
	TeamId        int64                              `json:"teamId"`        // The team ID of this participant, indicating the participant's team
	Spell2Id      int64                              `json:"spell2Id"`      // The ID of the second summoner spell used by this participant
	Masteries     []CurrentGameParticipantMasteryDTO `json:"masteries"`     // The masteries used by this participant
	Spell1Id      int64                              `json:"spell1Id"`      // The ID of the first summoner spell used by this participant
	SummonerId    string                             `json:"summonerId"`    // The encrypted summoner ID of this participant
	Perks         Perks                              `json:"perks"`
}

type Perks struct {
	PerkIDs      []int64 `json:"perkIDs"`
	PerkStyle    int64   `json:"perkStyle"`
	PerkSubStyle int64   `json:"perkSubStyle"`
}

type CurrentGameParticipantRuneDTO struct {
	Count  int   `json:"count"`  // The count of this rune used by the participant
	RuneId int64 `json:"runeId"` // The ID of the rune
}

type CurrentGameParticipantMasteryDTO struct {
	MasteryId int64 `json:"masteryId"` // The ID of the mastery
	Rank      int   `json:"rank"`      // The number of points put into this mastery by the user
}

func (c *client) GetCurrentGameInfoBySummoner(ctx context.Context, r region.Region, summonerID string) (*CurrentGameInfo, error) {
	var res CurrentGameInfo
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/spectator/v4/active-games/by-summoner", fmt.Sprintf("/%s", summonerID), nil, &res)
	return &res, err
}

type FeaturedGames struct {
	ClientRefreshInterval int64                 `json:"clientRefreshInterval",datastore:",noindex"` // The suggested interval to wait before requesting FeaturedGames again
	GameList              []FeaturedGameInfoDTO `json:"gameList",datastore:",noindex"`              // 	The list of featured games
}

type FeaturedGameInfoDTO struct {
	GameId            int64                        `json:"gameId"`            // The ID of the game
	GameStartTime     int64                        `json:"gameStartTime"`     // The game start time represented in epoch milliseconds
	PlatformId        string                       `json:"platformId"`        // The ID of the platform on which the game is being played
	GameMode          string                       `json:"gameMode"`          // The game mode
	MapId             int64                        `json:"mapId"`             // The ID of the map
	GameType          string                       `json:"gameType"`          // The game type
	BannedChampions   []BannedChampion             `json:"bannedChampions"`   // 	Banned champion information
	Observers         Observer                     `json:"observers"`         // The observer information
	Participants      []FeaturedGameParticipantDTO `json:"participants"`      //The participant information
	GameLength        int64                        `json:"gameLength"`        // The amount of time in seconds that has passed since the game started
	GameQueueConfigId int64                        `json:"gameQueueConfigId"` // The queue type (queue types are documented on the Game Constants page)
}

type FeaturedGameParticipantDTO struct {
	ProfileIconId int64  `json:"profileIconId"` // The ID of the profile icon used by this participant
	ChampionId    int64  `json:"championId"`    // The ID of the champion played by this participant
	SummonerName  string `json:"summonerName"`  // The summoner name of this participant
	Bot           bool   `json:"bot"`           // Flag indicating whether or not this participant is a bot
	Spell2Id      int64  `json:"spell2Id"`      // The ID of the second summoner spell used by this participant
	TeamId        int64  `json:"teamId"`        // The team ID of this participant, indicating the participant's team
	Spell1Id      int64  `json:"spell1Id"`      // The ID of the first summoner spell used by this participant
	Perks         Perks  `json:"perks"`
}

func (c *client) GetFeaturedGames(ctx context.Context, r region.Region) (*FeaturedGames, error) {
	var res FeaturedGames
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/spectator/v4/featured-games", "", nil, &res)
	return &res, err
}
