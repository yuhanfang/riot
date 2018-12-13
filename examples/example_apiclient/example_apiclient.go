// This script demonstrates how to use the apiclient package.
//
// Before running, set the RIOT_APIKEY enviornment variable to your developer
// key from the Riot developer portal.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/champion"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

const (
	playerID = "x9k0laU59wtIYnd8zt1dZmtJ_wXl13bqjhTRRC8FPwTbYVA" // These are encrypted per the api key used
	name     = "waddlechirp"
	account  = "hJN7Yl1FSZLD4vGKUIAMVFI_IWqK7WmY6Lb9S2QGRSUes8U" // These are encrypted per the api key used
	game     = 2644987649
	league   = "6b5c7950-5260-11e7-8125-c81f66dbb56c"
	reg      = region.NA1
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
	key := os.Getenv("RIOT_APIKEY")
	httpClient := http.DefaultClient
	ctx := context.Background()
	limiter := ratelimit.NewLimiter()
	client := apiclient.New(key, httpClient, limiter)

	// Champion mastery

	fmt.Println("GetAllChampionMasteries")
	masteries, err := client.GetAllChampionMasteries(ctx, reg, playerID)
	prettyPrint(masteries, err)

	fmt.Println("GetChampionMastery")
	asheMastery, err := client.GetChampionMastery(ctx, reg, playerID, champion.Ashe)
	prettyPrint(asheMastery, err)

	fmt.Println("GetChampionMasteryScore")
	score, err := client.GetChampionMasteryScore(ctx, reg, playerID)
	prettyPrint(score, err)

	// Champions

	fmt.Println("GetChampions")
	champions, err := client.GetChampions(ctx, reg)
	prettyPrint(champions, err)

	fmt.Println("GetChampionByID")
	ashe, err := client.GetChampionByID(ctx, reg, champion.Ashe)
	prettyPrint(ashe, err)

	// League

	fmt.Println("GetChallengerLeague")
	challenger, err := client.GetChallengerLeague(ctx, reg, queue.RankedSolo5x5)
	prettyPrint(challenger, err)

	fmt.Println("GetMasterLeague")
	master, err := client.GetMasterLeague(ctx, reg, queue.RankedSolo5x5)
	prettyPrint(master, err)

	fmt.Println("GetAllLeaguePositionsForSummoner")
	positions, err := client.GetAllLeaguePositionsForSummoner(ctx, reg, playerID)
	prettyPrint(positions, err)

	fmt.Println("GetLeagueByID")
	myLeague, err := client.GetLeagueByID(ctx, reg, league)
	prettyPrint(myLeague, err)

	// Match

	fmt.Println("GetMatch")
	myMatch, err := client.GetMatch(ctx, reg, game)
	prettyPrint(myMatch, err)

	fmt.Println("GetMatchTimeline")
	timeline, err := client.GetMatchTimeline(ctx, reg, game)
	prettyPrint(timeline, err)

	fmt.Println("GetMatchlist")
	matchlist, err := client.GetMatchlist(ctx, reg, account, nil)
	prettyPrint(matchlist, err)

	fmt.Println("GetRecentMatchlist")
	recentMatchlist, err := client.GetRecentMatchlist(ctx, reg, account)
	prettyPrint(recentMatchlist, err)

	// Spectator

	fmt.Println("GetCurrentGameInfoBySummoner")
	current, err := client.GetCurrentGameInfoBySummoner(ctx, reg, playerID)
	prettyPrint(current, err)

	fmt.Println("GetFeaturedGames")
	featured, err := client.GetFeaturedGames(ctx, reg)
	prettyPrint(featured, err)

	// Summoner
	fmt.Println("GetByAccountID")
	summoner, err := client.GetByAccountID(ctx, reg, account)
	prettyPrint(summoner, err)

	fmt.Println("GetBySummonerName")
	summoner, err = client.GetBySummonerName(ctx, reg, name)
	prettyPrint(summoner, err)

	fmt.Println("GetBySummonerID")
	summoner, err = client.GetBySummonerID(ctx, reg, playerID)
	prettyPrint(summoner, err)

}
