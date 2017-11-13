package apiclient

import (
	"fmt"
	"net/url"
	"time"

	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/event"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/constants/season"
	"golang.org/x/net/context"
)

type MatchDTO struct {
	SeasonID              int
	QueueID               int
	GameID                int64
	ParticipantIdentities []ParticipantIdentityDTO
	GameVersion           string
	PlatformID            string
	GameMode              string
	MapID                 int
	GameType              string
	Teams                 []TeamStatsDTO
	Participants          []ParticipantDTO
	GameDuration          int64
	GameCreation          int64
}

type ParticipantIdentityDTO struct {
	Player        PlayerDTO
	ParticipantID int
}

type PlayerDTO struct {
	CurrentPlatformID string
	SummonerName      string
	MatchHistoryUri   string
	PlatformID        string
	CurrentAccountID  int64
	ProfileIcon       int
	SummonerID        int64
	AccountID         int64
}

type TeamStatsDTO struct {
	FirstDragon          bool
	FirstInhibitor       bool
	Bans                 []TeamBansDTO
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

type TeamBansDTO struct {
	PickTurn   int
	ChampionID champion.Champion
}

type ParticipantDTO struct {
	Stats                     ParticipantStatsDTO
	ParticipantId             int
	Runes                     []RuneDTO
	Timeline                  ParticipantTimelineDTO
	TeamId                    int
	Spell2Id                  int
	Masteries                 []MasteryDTO
	HighestAchievedSeasonTier string
	Spell1Id                  int
	ChampionId                champion.Champion
}

type ParticipantStatsDTO struct {
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

type RuneDTO struct {
	RuneID int
	Rank   int
}

type ParticipantTimelineDTO struct {
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

type MasteryDTO struct {
	MasteryId int
	Rank      int
}

func (c *client) GetMatch(ctx context.Context, r region.Region, matchID int64) (*MatchDTO, error) {
	var res MatchDTO
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v3/matches", fmt.Sprintf("/%d", matchID), nil, &res)
	return &res, err
}

type MatchlistDTO struct {
	Matches    []MatchReferenceDTO
	TotalGames int
	StartIndex int
	EndIndex   int
}

type MatchReferenceDTO struct {
	Lane       string
	GameID     int64
	Champion   champion.Champion
	PlatformID string
	Season     season.Season
	Queue      int
	Role       string
	Timestamp  int64
}

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

func (c *client) GetMatchlist(ctx context.Context, r region.Region, accountID int64, opts *GetMatchlistOptions) (*MatchlistDTO, error) {
	var (
		res  MatchlistDTO
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

func (c *client) GetRecentMatchlist(ctx context.Context, r region.Region, accountID int64) (*MatchlistDTO, error) {
	var res MatchlistDTO
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v3/matchlists/by-account", fmt.Sprintf("/%d/recent", accountID), nil, &res)
	return &res, err
}

type MatchTimelineDTO struct {
	Frames        []MatchFrameDTO
	FrameInterval int64
}

type MatchFrameDTO struct {
	Timestamp         int64
	ParticipantFrames map[int]MatchParticipantFrameDTO
	Events            []MatchEventDTO
}

type MatchParticipantFrameDTO struct {
	TotalGold           int
	TeamScore           int
	ParticipantID       int
	Level               int
	CurrentGold         int
	MinionsKilled       int
	DominionScore       int
	Position            MatchPositionDTO
	XP                  int
	JungleMinionsKilled int
}

type MatchPositionDTO struct {
	Y int
	X int
}

type MatchEventDTO struct {
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
	Position                MatchPositionDTO
	BeforeID                int
}

func (c *client) GetMatchTimeline(ctx context.Context, r region.Region, matchID int64) (*MatchTimelineDTO, error) {
	var res MatchTimelineDTO
	_, err := c.dispatchAndUnmarshal(ctx, r, "/lol/match/v3/timelines/by-match", fmt.Sprintf("/%d", matchID), nil, &res)
	return &res, err
}
