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
	leagues, err := client.GetLeaguesByID(ctx, 9)
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
	fmt.Printf("https://acs.leagueoflegends.com/v1/stats/game/%s/%s?gameHash=%s\n", realm, match.Games[game.ID].GameID, game.GameHash)
	fmt.Println(err)
}
