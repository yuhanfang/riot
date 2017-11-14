package apiclient

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/event"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/season"
)

type Match struct {
	SeasonID              int
	QueueID               int
	GameID                int64
	ParticipantIdentities []ParticipantIdentity
	GameVersion           string
	PlatformID            string
	GameMode              string
	MapID                 int
	GameType              string
	Teams                 []TeamStats
	Participants          []Participant
	GameDuration          int64
	GameCreation          int64
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
	TeamId               int
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
	ParticipantId             int
	Runes                     []Rune
	Timeline                  ParticipantTimeline
	TeamId                    int
	Spell2Id                  int
	Masteries                 []Mastery
	HighestAchievedSeasonTier string
	Spell1Id                  int
	ChampionId                champion.Champion
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
}

type Rune struct {
	RuneID int
	Rank   int
}

type ParticipantTimeline struct {
	Lane                        string
	ParticipantID               int
	CSDiffPerMinDeltas          map[string]float64
	GoldPerMinDeltas            map[string]float64
	XPDiffPerMinDeltas          map[string]float64
	CreepsPerMinDeltas          map[string]float64
	XPPerMinDeltas              map[string]float64
	Role                        string
	DamageTakenDiffPerMinDeltas map[string]float64
	DamageTakenPerMinDeltas     map[string]float64
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
	Matches    []MatchReference
	TotalGames int
	StartIndex int
	EndIndex   int
}

type MatchReference struct {
	Lane       string
	GameID     int64
	Champion   champion.Champion
	PlatformID string
	Season     season.Season
	Queue      int
	Role       string
	Timestamp  int64
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
	return t.UnixNano() / int64(time.Millisecond)
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
	Frames        []MatchFrame
	FrameInterval int64
}

type MatchFrame struct {
	Timestamp         int64
	ParticipantFrames map[int]MatchParticipantFrame
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
	Timestamp               int64
	AfterID                 int
	MonsterSubType          string
	LaneType                string
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
