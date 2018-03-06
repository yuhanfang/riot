package league

import "fmt"

// League represents an esports league supported by the API.
type League int64

const (
	LeagueInvalid               League = 0
	LeagueAllStar               League = 1
	LeagueNALCS                 League = 2
	LeagueEULCS                 League = 3
	LeagueNACS                  League = 4
	LeagueEUCS                  League = 5
	LeagueLCK                   League = 6
	LeagueLPL                   League = 7
	LeagueLMS                   League = 8
	LeagueWorlds                League = 9
	LeagueMidSeasonInvitational League = 10
	LeagueInternationalWildcard League = 12
	LeagueOPL                   League = 13
	LeagueCBLOL                 League = 14
)

// SlugToLeague returns a League constant for a given slug.
func SlugToLeague(slug string) League {
	switch slug {
	case "all-star":
		return LeagueAllStar
	case "na-lcs":
		return LeagueNALCS
	case "eu-lcs":
		return LeagueEULCS
	case "na-cs":
		return LeagueNACS
	case "eu-cs":
		return LeagueEUCS
	case "lck":
		return LeagueLCK
	case "lpl-china":
		return LeagueLPL
	case "lms":
		return LeagueLMS
	case "worlds":
		return LeagueWorlds
	case "msi":
		return LeagueMidSeasonInvitational
	case "iwc":
		return LeagueInternationalWildcard
	case "oce-opl":
		return LeagueOPL
	case "cblol-brazil":
		return LeagueCBLOL
	default:
		return LeagueInvalid
	}
}

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
