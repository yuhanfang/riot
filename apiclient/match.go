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
	SeasonID              season.Season         `datastore:",noindex"`
	QueueID               queue.Queue           `datastore:",noindex"`
	GameID                int64                 `datastore:",noindex"`
	ParticipantIdentities []ParticipantIdentity `datastore:",noindex"`
	GameVersion           string                `datastore:",noindex"`
	PlatformID            string                `datastore:",noindex"`
	GameMode              string                `datastore:",noindex"`
	MapID                 int                   `datastore:",noindex"`
	GameType              string                `datastore:",noindex"`
	Teams                 []TeamStats           `datastore:",noindex"`
	Participants          []Participant         `datastore:",noindex"`
	GameDuration          types.Milliseconds    `datastore:",noindex"`
	GameCreation          types.Milliseconds    `datastore:",noindex"`
}

type ParticipantIdentity struct {
	Player        Player
	ParticipantID int
}

type Player struct {
	CurrentPlatformID string
	SummonerName      string
	MatchHistoryUri   string
	PlatformID        string
	CurrentAccountID  int64
	ProfileIcon       int
	SummonerID        int64
	AccountID         int64
}

type TeamStats struct {
	FirstDragon          bool
	FirstInhibitor       bool
	Bans                 []TeamBans
	BaronKills           int
	FirstRiftHerald      bool
	FirstBaron           bool
	RiftHeraldKills      int
	FirstBlood           bool
	TeamID               int
	FirstTower           bool
	VilemawKills         int
	InhibitorKills       int
	TowerKills           int
	DominionVictoryScore int
	Win                  string
	DragonKills          int
}

type TeamBans struct {
	PickTurn   int
	ChampionID champion.Champion
}

type Participant struct {
	Stats                     ParticipantStats
	ParticipantID             int
	Runes                     []Rune
	Timeline                  ParticipantTimeline
	TeamID                    int
	Spell2ID                  int
	Masteries                 []Mastery
	HighestAchievedSeasonTier string
	Spell1ID                  int
	ChampionID                champion.Champion
}

type ParticipantStats struct {
	PhysicalDamageDealt             int64
	NeutralMinionsKilledTeamJungle  int
	MagicDamageDealt                int64
	TotalPlayerScore                int
	Deaths                          int
	Win                             bool
	NeutralMinionsKilledEnemyJungle int
	AltarsCaptured                  int
	LargestCriticalStrike           int
	TotalDamageDealt                int64
	MagicDamageDealtToChampions     int64
	VisionWardsBoughtInGame         int
	DamageDealtToObjectives         int64
	LargestKillingSpree             int
	Item1                           int
	QuadraKills                     int
	TeamObjective                   int
	TotalTimeCrowdControlDealt      int
	LongestTimeSpentLiving          int
	WardsKilled                     int
	FirstTowerAssist                bool
	FirstTowerKill                  bool
	Item2                           int
	Item3                           int
	Item0                           int
	FirstBloodAssist                bool
	VisionScore                     int64
	WardsPlaced                     int
	Item4                           int
	Item5                           int
	Item6                           int
	TurretKills                     int
	TripleKills                     int
	DamageSelfMitigated             int64
	ChampLevel                      int
	NodeNeutralizeAssist            int
	FirstInhibitorKill              bool
	GoldEarned                      int
	MagicalDamageTaken              int64
	Kills                           int
	DoubleKills                     int
	NodeCaptureAssist               int
	TrueDamageTaken                 int64
	NodeNeutralize                  int
	FirstInhibitorAssist            bool
	Assists                         int
	UnrealKills                     int
	NeutralMinionsKilled            int
	ObjectivePlayerScore            int
	CombatPlayerScore               int
	DamageDealtToTurrets            int64
	AltarsNeutralized               int
	PhysicalDamageDealtToChampions  int64
	GoldSpent                       int
	TrueDamageDealt                 int64
	TrueDamageDealtToChampions      int64
	ParticipantID                   int
	PentaKills                      int
	TotalHeal                       int64
	TotalMinionsKilled              int
	FirstBloodKill                  bool
	NodeCapture                     int
	LargestMultiKill                int
	SightWardsBoughtInGame          int
	TotalDamageDealtToChampions     int64
	TotalUnitsHealed                int
	InhibitorKills                  int
	TotalScoreRank                  int
	TotalDamageTaken                int64
	KillingSprees                   int
	TimeCCingOthers                 int64
	PhysicalDamageTaken             int64

	Perk0     int64
	Perk0Var1 int
	Perk0Var2 int
	Perk0Var3 int

	Perk1     int64
	Perk1Var1 int
	Perk1Var2 int
	Perk1Var3 int

	Perk2     int64
	Perk2Var1 int
	Perk2Var2 int
	Perk2Var3 int

	Perk3     int64
	Perk3Var1 int
	Perk3Var2 int
	Perk3Var3 int

	Perk4     int64
	Perk4Var1 int
	Perk4Var2 int
	Perk4Var3 int

	Perk5     int64
	Perk5Var1 int
	Perk5Var2 int
	Perk5Var3 int

	PerkPrimaryStyle int64
	PerkSubStyle     int64
}

type Rune struct {
	RuneID int
	Rank   int
}

// Interval represents a range of game time, measured in minutes.
//
// The value 999 is used to represent an endpoint that is coded as the literal
// string "end"
type Interval struct {
	Begin int
	End   int
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
	Interval Interval
	Value    float64
}

type ParticipantTimeline struct {
	Lane                        lane.Lane
	ParticipantID               int
	CSDiffPerMinDeltas          IntervalValues
	GoldPerMinDeltas            IntervalValues
	XPDiffPerMinDeltas          IntervalValues
	CreepsPerMinDeltas          IntervalValues
	XPPerMinDeltas              IntervalValues
	Role                        string
	DamageTakenDiffPerMinDeltas IntervalValues
	DamageTakenPerMinDeltas     IntervalValues
}

type Mastery struct {
	MasteryID int
	Rank      int
}

func (c *client) GetMatch(ctx context.Context, r region.Region, matchID int64) (*Match, error) {
	var res Match
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v3/matches", fmt.Sprintf("/%d", matchID), nil, &res)
	return &res, err
}

type Matchlist struct {
	Matches    []MatchReference `datastore:",noindex"`
	TotalGames int              `datastore:",noindex"`
	StartIndex int              `datastore:",noindex"`
	EndIndex   int              `datastore:",noindex"`
}

type MatchReference struct {
	Lane       lane.Lane
	GameID     int64
	Champion   champion.Champion
	PlatformID string
	Season     season.Season
	Queue      queue.Queue
	Role       string
	Timestamp  types.Milliseconds
}

// GetMatchlistOptions provides filtering options for GetMatchlist. The zero
// value means that the option will not be used in filtering.
type GetMatchlistOptions struct {
	Queue      []queue.Queue
	Season     []season.Season
	Champion   []champion.Champion
	BeginTime  *time.Time
	EndTime    *time.Time
	BeginIndex *int
	EndIndex   *int
}

func timeToUnixMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func (c *client) GetMatchlist(ctx context.Context, r region.Region, accountID int64, opts *GetMatchlistOptions) (*Matchlist, error) {
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
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v3/matchlists/by-account", fmt.Sprintf("/%d", accountID), vals, &res)
	return &res, err
}

func (c *client) GetRecentMatchlist(ctx context.Context, r region.Region, accountID int64) (*Matchlist, error) {
	var res Matchlist
	// Recent matchlists are a separate API call from matchlists, even though
	// both have the same Method. Add "recent" as a uniquifier for this method.
	_, err := c.dispatchAndUnmarshalWithUniquifier(ctx, r, "/lol/match/v3/matchlists/by-account", fmt.Sprintf("/%d/recent", accountID), nil, "recent", &res)
	return &res, err
}

type MatchTimeline struct {
	Frames        []MatchFrame       `datastore:",noindex"`
	FrameInterval types.Milliseconds `datastore:",noindex"`
}

// ParticipantFrames stores frames corresponding to each participant. The order
// is not defined (i.e. do not assume the order is ascending by participant ID).
type ParticipantFrames struct {
	Frames []MatchParticipantFrame
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
	Timestamp         types.Milliseconds
	ParticipantFrames ParticipantFrames
	Events            []MatchEvent
}

type MatchParticipantFrame struct {
	TotalGold           int
	TeamScore           int
	ParticipantID       int
	Level               int
	CurrentGold         int
	MinionsKilled       int
	DominionScore       int
	Position            MatchPosition
	XP                  int
	JungleMinionsKilled int
}

type MatchPosition struct {
	Y int
	X int
}

type MatchEvent struct {
	EventType               string
	TowerType               string
	TeamID                  int
	AscendedType            string
	KillerID                int
	LevelUpType             string
	PointCaptured           string
	AssistingParticipantIDs []int
	WardType                string
	MonsterType             string
	Type                    event.Event
	SkillSlot               int
	VictimID                int
	Timestamp               types.Milliseconds
	AfterID                 int
	MonsterSubType          string
	LaneType                lane.Type
	ItemID                  int
	ParticipantID           int
	BuildingType            string
	CreatorID               int
	Position                MatchPosition
	BeforeID                int
}

func (c *client) GetMatchTimeline(ctx context.Context, r region.Region, matchID int64) (*MatchTimeline, error) {
	var res MatchTimeline
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v3/timelines/by-match", fmt.Sprintf("/%d", matchID), nil, &res)
	return &res, err
}
