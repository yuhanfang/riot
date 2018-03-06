package esports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/esports/tournaments"
)

func Example() {
	ctx := context.Background()
	c := NewClient(http.DefaultClient)
	leagues, err := c.GetLeagues(ctx, LeagueNALCS)
	if err != nil {
		panic(err)
	}
	var tournament Leagues_HighlanderTournament
	for _, t := range leagues.HighlanderTournaments {
		if t.Title == tournaments.NA2018Spring {
			tournament = t
			break
		}
	}

	var bracket Leagues_HighlanderTournament_Bracket
	for _, b := range tournament.Brackets {
		if b.Name == "regular_season" {
			bracket = b
		}
	}

	var match *apiclient.Match
	for _, m := range bracket.Matches {
		fmt.Println("m.ID:", m.ID)
		matchDetails, err := c.GetHighlanderMatchDetails(ctx, tournament.ID, m.ID)
		if err != nil {
			panic(err)
		}
		for _, g := range m.Games {
			gameID := g.GameID
			// Game hasn't taken place yet.
			if gameID == 0 {
				continue
			}
			region := g.GameRealm

			fmt.Println("gameID:", gameID)

			var gameHash string
			for _, mapping := range matchDetails.GameIDMappings {
				fmt.Println("mapping.ID:", mapping.ID)
				if mapping.ID == g.ID {
					gameHash = mapping.GameHash
					break
				}
			}

			fmt.Println("gameHash:", gameHash)

			match, err = c.GetGameStats(ctx, region, gameID.Int64(), gameHash)
			if err != nil {
				panic(err)
			}
			break
		}
		break
	}

	fmt.Printf("%+v", match.Participants[0])

	// fmt.Printf("%+v", tournament.Brackets)

	// Output:
}

func ExampleGameDetailsEndToEnd() {
	ctx := context.Background()
	c := NewClient(http.DefaultClient)
	leagues, err := c.GetLeagues(ctx, LeagueNALCS)
	if err != nil {
		panic(err)
	}
	tournament := leagues.HighlanderTournaments[2]
	fmt.Println("tournament:", tournament.Title)
	bracket := tournament.Brackets["25cf16fe-5fac-492e-aa18-887b5461e70f"]
	fmt.Println("bracket:", bracket.Name)
	match := bracket.Matches["afbb3bd6-48b7-40b6-bb57-b1bd48206c21"]
	fmt.Println("match:", match.Name)
	game := match.Games["e425aba9-9317-45d9-a457-360c59df1d81"]

	gameID := game.GameID
	region := game.GameRealm
	fmt.Println("gameID:", gameID)
	fmt.Println("region:", region)

	details, err := c.GetHighlanderMatchDetails(ctx, tournament.ID, match.ID)
	if err != nil {
		panic(err)
	}

	stats, err := c.GetGameStats(ctx, region, gameID.Int64(), details.GameIDMappings[0].GameHash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", stats)
}
