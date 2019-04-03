package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/event"
	"github.com/yuhanfang/riot/constants/lane"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/season"
	"github.com/yuhanfang/riot/types"
)

type Match struct {
	SeasonID              season.Season         `json:"seasonID",datastore:",noindex"`              // SeasonID is the ID associated with the current season of league.
	QueueID               queue.Queue           `json:"queueID",datastore:",noindex"`               // QueueID is a constant that refers to the queue.
	GameID                int64                 `json:"gameID",datastore:",noindex"`                // GameID is the ID of the requested game.
	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities",datastore:",noindex"` // Array of player identities in requested game.
	GameVersion           string                `json:"gameVersion",datastore:",noindex"`           // Version that the game was played on.
	PlatformID            string                `json:"platformID",datastore:",noindex"`            // PlatformID
	GameMode              string                `json:"gameMode",datastore:",noindex"`              // GameMode
	MapID                 int                   `json:"mapID",datastore:",noindex"`                 // MapID is a constant of the map played on.
	GameType              string                `json:"gameType",datastore:",noindex"`              // GameType
	Teams                 []TeamStats           `json:"teams",datastore:",noindex"`                 // Teams
	Participants          []Participant         `json:"participants",datastore:",noindex"`          // Participants
	GameDuration          types.Milliseconds    `json:"gameDuration",datastore:",noindex"`          // GameDuration is the duration of the game in milliseconds
	GameCreation          types.Milliseconds    `json:"gameCreation",datastore:",noindex"`          // GameCreation is when game was created in epoch
}

type ParticipantIdentity struct {
	Player        Player `json:"player"`
	ParticipantID int    `json:"participantID"`
}

type Player struct {
	CurrentPlatformID string `json:"currentPlatformID"`
	SummonerName      string `json:"summonerName"`
	MatchHistoryUri   string `json:"matchHistoryUri"`
	PlatformID        string `json:"platformID"`
	CurrentAccountID  string `json:"currentAccountID"`
	ProfileIcon       int    `json:"profileIcon"`
	SummonerID        string `json:"summonerID"`
	AccountID         string `json:"accountID"`
}

type TeamStats struct {
	FirstDragon          bool       `json:"firstDragon"`
	FirstInhibitor       bool       `json:"firstInhibitor"`
	Bans                 []TeamBans `json:"bans"`
	BaronKills           int        `json:"baronKills"`
	FirstRiftHerald      bool       `json:"firstRiftHerald"`
	FirstBaron           bool       `json:"firstBaron"`
	RiftHeraldKills      int        `json:"riftHeraldKills"`
	FirstBlood           bool       `json:"firstBlood"`
	TeamID               int        `json:"teamID"`
	FirstTower           bool       `json:"firstTower"`
	VilemawKills         int        `json:"vilemawKills"`
	InhibitorKills       int        `json:"inhibitorKills"`
	TowerKills           int        `json:"towerKills"`
	DominionVictoryScore int        `json:"dominionVictoryScore"`
	Win                  string     `json:"win"`
	DragonKills          int        `json:"dragonKills"`
}

type TeamBans struct {
	PickTurn   int               `json:"pickTurn"`
	ChampionID champion.Champion `json:"championID"`
}

type Participant struct {
	Stats                     ParticipantStats    `json:"stats"`
	ParticipantID             int                 `json:"participantID"`
	Runes                     []Rune              `json:"runes"`
	Timeline                  ParticipantTimeline `json:"timeline"`
	TeamID                    int                 `json:"teamID"`
	Spell2ID                  int                 `json:"spell2ID"`
	Masteries                 []Mastery           `json:"masteries"`
	HighestAchievedSeasonTier string              `json:"highestAchievedSeasonTier"`
	Spell1ID                  int                 `json:"spell1ID"`
	ChampionID                champion.Champion   `json:"championID"`
}

type ParticipantStats struct {
	PhysicalDamageDealt             int64 `json:"physicalDamageDealt"`
	NeutralMinionsKilledTeamJungle  int   `json:"neutralMinionsKilledTeamJungle"`
	MagicDamageDealt                int64 `json:"magicDamageDealt"`
	TotalPlayerScore                int   `json:"totalPlayerScore"`
	Deaths                          int   `json:"deaths"`
	Win                             bool  `json:"win"`
	NeutralMinionsKilledEnemyJungle int   `json:"neutralMinionsKilledEnemyJungle"`
	AltarsCaptured                  int   `json:"altarsCaptured"`
	LargestCriticalStrike           int   `json:"largestCriticalStrike"`
	TotalDamageDealt                int64 `json:"totalDamageDealt"`
	MagicDamageDealtToChampions     int64 `json:"magicDamageDealtToChampions"`
	VisionWardsBoughtInGame         int   `json:"visionWardsBoughtInGame"`
	DamageDealtToObjectives         int64 `json:"damageDealtToObjectives"`
	LargestKillingSpree             int   `json:"largestKillingSpree"`
	Item1                           int   `json:"item1"`
	QuadraKills                     int   `json:"quadraKills"`
	TeamObjective                   int   `json:"teamObjective"`
	TotalTimeCrowdControlDealt      int   `json:"totalTimeCrowdControlDealt"`
	LongestTimeSpentLiving          int   `json:"longestTimeSpentLiving"`
	WardsKilled                     int   `json:"wardsKilled"`
	FirstTowerAssist                bool  `json:"firstTowerAssist"`
	FirstTowerKill                  bool  `json:"firstTowerKill"`
	Item2                           int   `json:"item2"`
	Item3                           int   `json:"item3"`
	Item0                           int   `json:"item0"`
	FirstBloodAssist                bool  `json:"firstBloodAssist"`
	VisionScore                     int64 `json:"visionScore"`
	WardsPlaced                     int   `json:"wardsPlaced"`
	Item4                           int   `json:"item4"`
	Item5                           int   `json:"item5"`
	Item6                           int   `json:"item6"`
	TurretKills                     int   `json:"turretKills"`
	TripleKills                     int   `json:"tripleKills"`
	DamageSelfMitigated             int64 `json:"damageSelfMitigated"`
	ChampLevel                      int   `json:"champLevel"`
	NodeNeutralizeAssist            int   `json:"nodeNeutralizeAssist"`
	FirstInhibitorKill              bool  `json:"firstInhibitorKill"`
	GoldEarned                      int   `json:"goldEarned"`
	MagicalDamageTaken              int64 `json:"magicalDamageTaken"`
	Kills                           int   `json:"kills"`
	DoubleKills                     int   `json:"doubleKills"`
	NodeCaptureAssist               int   `json:"nodeCaptureAssist"`
	TrueDamageTaken                 int64 `json:"trueDamageTaken"`
	NodeNeutralize                  int   `json:"nodeNeutralize"`
	FirstInhibitorAssist            bool  `json:"firstInhibitorAssist"`
	Assists                         int   `json:"assists"`
	UnrealKills                     int   `json:"unrealKills"`
	NeutralMinionsKilled            int   `json:"neutralMinionsKilled"`
	ObjectivePlayerScore            int   `json:"objectivePlayerScore"`
	CombatPlayerScore               int   `json:"combatPlayerScore"`
	DamageDealtToTurrets            int64 `json:"damageDealtToTurrets"`
	AltarsNeutralized               int   `json:"altarsNeutralized"`
	PhysicalDamageDealtToChampions  int64 `json:"physicalDamageDealtToChampions"`
	GoldSpent                       int   `json:"goldSpent"`
	TrueDamageDealt                 int64 `json:"trueDamageDealt"`
	TrueDamageDealtToChampions      int64 `json:"trueDamageDealtToChampions"`
	ParticipantID                   int   `json:"participantID"`
	PentaKills                      int   `json:"pentaKills"`
	TotalHeal                       int64 `json:"totalHeal"`
	TotalMinionsKilled              int   `json:"totalMinionsKilled"`
	FirstBloodKill                  bool  `json:"firstBloodKill"`
	NodeCapture                     int   `json:"nodeCapture"`
	LargestMultiKill                int   `json:"largestMultiKill"`
	SightWardsBoughtInGame          int   `json:"sightWardsBoughtInGame"`
	TotalDamageDealtToChampions     int64 `json:"totalDamageDealtToChampions"`
	TotalUnitsHealed                int   `json:"totalUnitsHealed"`
	InhibitorKills                  int   `json:"inhibitorKills"`
	TotalScoreRank                  int   `json:"totalScoreRank"`
	TotalDamageTaken                int64 `json:"totalDamageTaken"`
	KillingSprees                   int   `json:"killingSprees"`
	TimeCCingOthers                 int64 `json:"timeCCingOthers"`
	PhysicalDamageTaken             int64 `json:"physicalDamageTaken"`

	Perk0     int64 `json:"perk0"`
	Perk0Var1 int   `json:"perk0Var1"`
	Perk0Var2 int   `json:"perk0Var2"`
	Perk0Var3 int   `json:"perk0Var3"`

	Perk1     int64 `json:"perk1"`
	Perk1Var1 int   `json:"perk1Var1"`
	Perk1Var2 int   `json:"perk1Var2"`
	Perk1Var3 int   `json:"perk1Var3"`

	Perk2     int64 `json:"perk2"`
	Perk2Var1 int   `json:"perk2Var1"`
	Perk2Var2 int   `json:"perk2Var2"`
	Perk2Var3 int   `json:"perk2Var3"`

	Perk3     int64 `json:"perk3"`
	Perk3Var1 int   `json:"perk3Var1"`
	Perk3Var2 int   `json:"perk3Var2"`
	Perk3Var3 int   `json:"perk3Var3"`

	Perk4     int64 `json:"perk4"`
	Perk4Var1 int   `json:"perk4Var1"`
	Perk4Var2 int   `json:"perk4Var2"`
	Perk4Var3 int   `json:"perk4Var3"`

	Perk5     int64 `json:"perk5"`
	Perk5Var1 int   `json:"perk5Var1"`
	Perk5Var2 int   `json:"perk5Var2"`
	Perk5Var3 int   `json:"perk5Var3"`

	PerkPrimaryStyle int64 `json:"perkPrimaryStyle"`
	PerkSubStyle     int64 `json:"perkSubStyle"`
}

type Rune struct {
	RuneID int `json:"runeID"`
	Rank   int `json:"rank"`
}

// Interval represents a range of game time, measured in minutes.
//
// The value 999 is used to represent an endpoint that is coded as the literal
// string "end"
type Interval struct {
	Begin int `json:"begin"`
	End   int `json:"end"`
}

// IntervalValues represents a mapping from intervals to values.
type IntervalValues []IntervalValue

func (i *IntervalValues) UnmarshalJSON(b []byte) error {
	var obj map[string]float64
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return err
	}
	var vals []IntervalValue
	for k, v := range obj {
		intervals := strings.Split(k, "-")
		if len(intervals) != 2 {
			return fmt.Errorf("unable to parse intervals: %v", intervals)
		}
		intervals[0] = strings.TrimSpace(intervals[0])
		intervals[1] = strings.TrimSpace(intervals[1])

		begin, err := strconv.ParseInt(intervals[0], 10, 64)
		if err != nil {
			return err
		}
		var end int64
		if intervals[1] == "end" {
			end = 999
		} else {
			end, err = strconv.ParseInt(intervals[1], 10, 64)
			if err != nil {
				return err
			}
		}
		vals = append(vals, IntervalValue{
			Interval: Interval{
				Begin: int(begin),
				End:   int(end),
			},
			Value: v,
		})
	}
	*i = IntervalValues(vals)
	return nil
}

type IntervalValue struct {
	Interval Interval `json:"interval"`
	Value    float64  `json:"value"`
}

type ParticipantTimeline struct {
	Lane                        lane.Lane      `json:"lane"`
	ParticipantID               int            `json:"participantID"`
	CSDiffPerMinDeltas          IntervalValues `json:"cSDiffPerMinDeltas"`
	GoldPerMinDeltas            IntervalValues `json:"goldPerMinDeltas"`
	XPDiffPerMinDeltas          IntervalValues `json:"xPDiffPerMinDeltas"`
	CreepsPerMinDeltas          IntervalValues `json:"creepsPerMinDeltas"`
	XPPerMinDeltas              IntervalValues `json:"xPPerMinDeltas"`
	Role                        string         `json:"role"`
	DamageTakenDiffPerMinDeltas IntervalValues `json:"damageTakenDiffPerMinDeltas"`
	DamageTakenPerMinDeltas     IntervalValues `json:"damageTakenPerMinDeltas"`
}

type Mastery struct {
	MasteryID int `json:"masteryID"`
	Rank      int `json:"rank"`
}

func (c *client) GetMatch(ctx context.Context, r region.Region, matchID int64) (*Match, error) {
	var res Match
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v4/matches", fmt.Sprintf("/%d", matchID), nil, &res)
	return &res, err
}

type Matchlist struct {
	Matches    []MatchReference `json:"matches",datastore:",noindex"`
	TotalGames int              `json:"totalGames",datastore:",noindex"`
	StartIndex int              `json:"startIndex",datastore:",noindex"`
	EndIndex   int              `json:"endIndex",datastore:",noindex"`
}

type MatchReference struct {
	Lane       lane.Lane          `json:"lane"`
	GameID     int64              `json:"gameID"`
	Champion   champion.Champion  `json:"champion"`
	PlatformID string             `json:"platformID"`
	Season     season.Season      `json:"season"`
	Queue      queue.Queue        `json:"queue"`
	Role       string             `json:"role"`
	Timestamp  types.Milliseconds `json:"timestamp"`
}

// GetMatchlistOptions provides filtering options for GetMatchlist. The zero
// value means that the option will not be used in filtering.
type GetMatchlistOptions struct {
	Queue      []queue.Queue       `json:"queue"`
	Season     []season.Season     `json:"season"`
	Champion   []champion.Champion `json:"champion"`
	BeginTime  *time.Time          `json:"beginTime"`
	EndTime    *time.Time          `json:"endTime"`
	BeginIndex *int                `json:"beginIndex"`
	EndIndex   *int                `json:"endIndex"`
}

func timeToUnixMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func (c *client) GetMatchlist(ctx context.Context, r region.Region, accountID string, opts *GetMatchlistOptions) (*Matchlist, error) {
	var (
		res  Matchlist
		vals url.Values
	)

	if opts != nil {
		vals = url.Values(make(map[string][]string))
		if len(opts.Queue) != 0 {
			for _, v := range opts.Queue {
				vals.Add("queue", fmt.Sprintf("%d", v))
			}
		}
		if len(opts.Season) != 0 {
			for _, v := range opts.Season {
				vals.Add("season", fmt.Sprintf("%d", v))
			}
		}
		if len(opts.Champion) != 0 {
			for _, v := range opts.Champion {
				vals.Add("champion", fmt.Sprintf("%d", v))
			}
		}
		if opts.BeginTime != nil {
			vals.Add("beginTime", fmt.Sprintf("%d", timeToUnixMilliseconds(*opts.BeginTime)))
		}
		if opts.EndTime != nil {
			vals.Add("endTime", fmt.Sprintf("%d", timeToUnixMilliseconds(*opts.EndTime)))
		}
		if opts.BeginIndex != nil {
			vals.Add("beginIndex", fmt.Sprintf("%d", *opts.BeginIndex))
		}
		if opts.EndIndex != nil {
			vals.Add("endIndex", fmt.Sprintf("%d", *opts.EndIndex))
		}
	}
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v4/matchlists/by-account", fmt.Sprintf("/%s", accountID), vals, &res)
	return &res, err
}

func (c *client) GetRecentMatchlist(ctx context.Context, r region.Region, accountID string) (*Matchlist, error) {
	var res Matchlist
	// Recent matchlists are a separate API call from matchlists, even though
	// both have the same Method. Add "recent" as a uniquifier for this method.
	_, err := c.dispatchAndUnmarshalWithUniquifier(ctx, r, "/lol/match/v4/matchlists/by-account", fmt.Sprintf("/%s/recent", accountID), nil, "recent", &res)
	return &res, err
}

type MatchTimeline struct {
	Frames        []MatchFrame       `json:"frames",datastore:",noindex"`
	FrameInterval types.Milliseconds `json:"frameInterval",datastore:",noindex"`
}

// ParticipantFrames stores frames corresponding to each participant. The order
// is not defined (i.e. do not assume the order is ascending by participant ID).
type ParticipantFrames struct {
	Frames []MatchParticipantFrame `json:"frames"`
}

func (p *ParticipantFrames) UnmarshalJSON(b []byte) error {
	var obj map[int]MatchParticipantFrame
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return err
	}
	var vals []MatchParticipantFrame

	for _, v := range obj {
		vals = append(vals, v)
	}
	p.Frames = vals
	return nil
}

type MatchFrame struct {
	Timestamp         types.Milliseconds `json:"timestamp"`
	ParticipantFrames ParticipantFrames  `json:"participantFrames"`
	Events            []MatchEvent       `json:"events"`
}

type MatchParticipantFrame struct {
	TotalGold           int           `json:"totalGold"`
	TeamScore           int           `json:"teamScore"`
	ParticipantID       int           `json:"participantID"`
	Level               int           `json:"level"`
	CurrentGold         int           `json:"currentGold"`
	MinionsKilled       int           `json:"minionsKilled"`
	DominionScore       int           `json:"dominionScore"`
	Position            MatchPosition `json:"position"`
	XP                  int           `json:"xP"`
	JungleMinionsKilled int           `json:"jungleMinionsKilled"`
}

type MatchPosition struct {
	Y int `json:"y"`
	X int `json:"x"`
}

type MatchEvent struct {
	EventType               string             `json:"eventType"`
	TowerType               string             `json:"towerType"`
	TeamID                  int                `json:"teamID"`
	AscendedType            string             `json:"ascendedType"`
	KillerID                int                `json:"killerID"`
	LevelUpType             string             `json:"levelUpType"`
	PointCaptured           string             `json:"pointCaptured"`
	AssistingParticipantIDs []int              `json:"assistingParticipantIDs"`
	WardType                string             `json:"wardType"`
	MonsterType             string             `json:"monsterType"`
	Type                    event.Event        `json:"type"`
	SkillSlot               int                `json:"skillSlot"`
	VictimID                int                `json:"victimID"`
	Timestamp               types.Milliseconds `json:"timestamp"`
	AfterID                 int                `json:"afterID"`
	MonsterSubType          string             `json:"monsterSubType"`
	LaneType                lane.Type          `json:"laneType"`
	ItemID                  int                `json:"itemID"`
	ParticipantID           int                `json:"participantID"`
	BuildingType            string             `json:"buildingType"`
	CreatorID               int                `json:"creatorID"`
	Position                MatchPosition      `json:"position"`
	BeforeID                int                `json:"beforeID"`
}

func (c *client) GetMatchTimeline(ctx context.Context, r region.Region, matchID int64) (*MatchTimeline, error) {
	var res MatchTimeline
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v4/timelines/by-match", fmt.Sprintf("/%d", matchID), nil, &res)
	return &res, err
}
