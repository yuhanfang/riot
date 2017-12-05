package esports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yuhanfang/riot/constants/language"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/types"
)

type Leagues struct {
	Leagues               []Leagues_League
	HighlanderTournaments []Leagues_HighlanderTournament
	HighlanderRecords     []Leagues_HighlanderRecord
	Teams                 []Leagues_Team
	// Players
}

type Leagues_League struct {
	ID          int64
	Slug        string
	Name        string
	GUID        string
	Region      string
	DrupalID    int64
	LogoURL     string
	CreatedAt   string // 2015-04-28T21:08:42.000Z
	UpdatedAt   string
	Abouts      map[language.Language]string
	Names       map[language.Language]string
	Tournaments []string
}

type Leagues_HighlanderTournament struct {
	ID              string
	Title           string
	Description     string
	LeagueReference string
	// RosteringStrategy
	// Queues
	// Rosters
	Published bool
	// Breakpoints
	Brackets map[string]Leagues_HighlanderTournament_Bracket
	// LiveMatches
	StartDate   string // YYYY-MM-DD
	EndDate     string // YYYY-MM-DD
	LeagueID    string
	PlatformIDs []string
	GameIDs     []string
	League      string
}

type Leagues_HighlanderTournament_Bracket struct {
	ID             string
	Name           string
	Position       int
	GroupPosition  int
	GroupName      string
	CanManufacture bool
	State          string
	// BracketType
	// MatchType
	// GameMode
	// Input
	Matches map[string]Leagues_HighlanderTournament_Bracket_Match
	// Standings
}

type Leagues_HighlanderTournament_Bracket_Match struct {
	ID            string
	Name          string
	Position      int
	State         string
	GroupPosition int
	Games         map[string]Leagues_HighlanderTournament_Bracket_Match_Game
}

type Leagues_HighlanderTournament_Bracket_Match_Game struct {
	ID            string
	Name          string
	GeneratedName string
	// GameMode
	// Input
	// Standings
	// Scores
	GameID     types.StringInt64
	GameRealm  region.Region
	PlatformID string
	Revision   int
}

type Leagues_HighlanderRecord struct {
	Wins       int
	Losses     int
	Ties       int
	Score      int
	Roster     string
	Tournament string
	Bracket    string
	ID         string
}

type Leagues_Team struct {
	ID           int64
	Slug         string
	Name         string
	GUID         string
	TeamPhotoURL string
	LogoURL      string
	Acroynm      string
	HomeLeague   string
	AltLogoURL   string
	CreatedAt    string // YYYY-MM-DDT18:34:47.000Z
	UpdatedAt    string
	Bios         map[language.Language]string
	// ForeignIDs
	Players  []int64
	Starters []int64
	Subs     []int64
}

// League represents an esports league supported by the API.
type League int64

const (
	LeagueAllStar               League = 1
	LeagueNALCS                        = 2
	LeagueEULCS                        = 3
	LeagueNACS                         = 4
	LeagueEUCS                         = 5
	LeagueLCK                          = 6
	LeagueLPL                          = 7
	LeagueLMS                          = 8
	LeagueWorlds                       = 9
	LeagueMidSeasonInvitational        = 10
	LeagueInternationalWildcard        = 12
	LeagueOPL                          = 13
	LeagueCBLOL                        = 14
)

// Slug returns a short string identifier for the league.
func (l League) Slug() string {
	switch l {
	case LeagueAllStar:
		return "all-star"
	case LeagueNALCS:
		return "na-lcs"
	case LeagueEULCS:
		return "eu-lcs"
	case LeagueNACS:
		return "na-cs"
	case LeagueEUCS:
		return "eu-cs"
	case LeagueLCK:
		return "lck"
	case LeagueLPL:
		return "lpl-china"
	case LeagueLMS:
		return "lms"
	case LeagueWorlds:
		return "worlds"
	case LeagueMidSeasonInvitational:
		return "msi"
	case LeagueInternationalWildcard:
		return "iwc"
	case LeagueOPL:
		return "oce-opl"
	case LeagueCBLOL:
		return "cblol-brazil"
	default:
		panic(fmt.Sprintf("unsupported league %d", l))
	}
}

// String returns a canonical name for the league.
func (l League) String() string {
	switch l {
	case LeagueAllStar:
		return "All-Star"
	case LeagueNALCS:
		return "NA LCS"
	case LeagueEULCS:
		return "EU LCS"
	case LeagueNACS:
		return "NA Challenger Series"
	case LeagueEUCS:
		return "EU Challenger Series"
	case LeagueLCK:
		return "LCK - Champions Korea"
	case LeagueLPL:
		return "LPL China"
	case LeagueLMS:
		return "LMS Taiwan"
	case LeagueWorlds:
		return "World Championship"
	case LeagueMidSeasonInvitational:
		return "Mid-Season Invitational"
	case LeagueInternationalWildcard:
		return "International Wildcard"
	case LeagueOPL:
		return "OPL"
	case LeagueCBLOL:
		return "CBLOL Brazil"
	default:
		panic(fmt.Sprintf("unsupported league %d", l))
	}
}

func (c Client) GetLeagues(ctx context.Context, id League) (*Leagues, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.lolesports.com/api/v1/leagues?id=%d", id), nil)
	if err != nil {
		return nil, err
	}
	var res Leagues
	_, err = c.doJSON(ctx, req, &res)
	return &res, err
}

func (c Client) GetLeaguesBySlug(ctx context.Context, slug string) (*Leagues, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.lolesports.com/api/v1/leagues?slug=%s", slug), nil)
	if err != nil {
		return nil, err
	}
	var res Leagues
	_, err = c.doJSON(ctx, req, &res)
	return &res, err
}
