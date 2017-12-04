// Package esports implements an API for interacting with lolesports.
//
// This package is currently unstable. Some attributes are commented out; those
// are available in the API, but currently ignored by us. Contributions are
// welcome!
package esports

type Leagues struct {
	Leagues               []League
	HighlanderTournaments []HighlanderTournament
	HighlanderRecords     []HighlanderRecord
	Teams                 []Team
	// Players
}

type League struct {
	ID          int64
	Slug        string
	Name        string
	GUID        string
	Region      string
	DrupalID    int64
	LogoURL     string
	CreatedAt   string // 2015-04-28T21:08:42.000Z
	UpdatedAt   string
	Abouts      map[constants.Language]string
	Names       map[constants.Language]string
	Tournaments []string
}

type HighlanderTournament struct {
	ID              string
	Title           string
	Description     string
	LeagueReference string
	// RosteringStrategy
	// Queues
	// Rosters
	Published bool
	// Breakpoints
	// Brackets
	// LiveMatches
	StartDate   string // YYYY-MM-DD
	EndDate     string // YYYY-MM-DD
	LeagueID    string
	PlatformIDs []string
	GameIDs     []string
	League      string
}

type HighlanderRecord struct {
	Wins       int
	Losses     int
	Ties       int
	Score      int
	Roster     string
	Tournament string
	Bracket    string
	ID         string
}

type Team struct {
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
	Bios         map[constants.Language]string
	// ForeignIDs
	Players  []int64
	Starters []int64
	Subs     []int64
}
