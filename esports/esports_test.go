package esports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yuhanfang/riot/esports/league"
)

func ExampleGameDetailsEndToEnd() {
	ctx := context.Background()
	c := NewClient(http.DefaultClient)
	leagues, err := c.GetLeagues(ctx, league.LeagueNALCS)
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
