// This script demonstrates how to use the esports package.
//
// The example is subject to change, since the esports package is currently
// unstable.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuhanfang/riot/esports"
	"github.com/yuhanfang/riot/esports/league"
)

func prettyPrint(res interface{}, err error) {
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}
	js, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}
	fmt.Println(string(js))
}

func main() {
	ctx := context.Background()
	client := esports.NewClient(http.DefaultClient)

	// Select a tournament, and pick a match in that tournament.
	leagues, err := client.GetLeagues(ctx, league.LeagueWorlds)
	tournament := leagues.HighlanderTournaments[0].ID
	var bracket esports.Leagues_HighlanderTournament_Bracket
	for _, b := range leagues.HighlanderTournaments[0].Brackets {
		bracket = b
		break
	}
	var (
		matchID string
		match   esports.Leagues_HighlanderTournament_Bracket_Match
	)
	for mi, m := range bracket.Matches {
		matchID = mi
		match = m
		break
	}

	// The match consists of possibly multiple games. Fetch the game IDs and
	// associated hashes. Pick the first game.
	details, err := client.GetHighlanderMatchDetails(ctx, tournament, matchID)
	game := details.GameIDMappings[0]

	// The realm is embedded in the match information. Note that the "ID" is used
	// to join Games with the match details.
	realm := match.Games[game.ID].GameRealm

	// Note that ACS uses GameID instead of ID. This is the URL with the game
	// stats.
	stats, err := client.GetGameStats(ctx, realm, match.Games[game.ID].GameID.Int64(), game.GameHash)
	prettyPrint(stats, err)

	timeline, err := client.GetGameTimeline(ctx, realm, match.Games[game.ID].GameID.Int64(), game.GameHash)
	prettyPrint(timeline, err)
}
