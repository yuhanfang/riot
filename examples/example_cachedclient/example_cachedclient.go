// Demonstrates how to set up a cached client. A cached client avoids hitting
// the Riot servers if data stored in the configured backend is fresh enough.
// This example uses Google Datastore as the cache backend.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/cachedclient"
	"github.com/yuhanfang/riot/cachedclient/google"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

const (
	playerID = 84289964
	name     = "waddlechirp"
	account  = 237254272
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
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	httpClient := http.DefaultClient
	ctx := context.Background()
	limiter := ratelimit.NewLimiter()
	c := apiclient.New(key, httpClient, limiter)
	ds, err := google.NewDatastore(ctx, project, "TestCachedClient")
	if err != nil {
		log.Fatal(err)
	}
	client := cachedclient.New(c, ds)
	fmt.Println("GetChallengerLeague")
	challenger, err := client.GetChallengerLeague(ctx, reg, queue.RankedSolo5x5)
	prettyPrint(challenger, err)

	fmt.Println("GetChallengerLeague")
	challenger, err = client.GetChallengerLeague(ctx, reg, queue.RankedSolo5x5)
	prettyPrint(challenger, err)
}
