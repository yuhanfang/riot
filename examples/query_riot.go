package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

func main() {
	key := os.Getenv("RIOT_APIKEY")
	httpClient := http.DefaultClient
	limiter := ratelimit.NewLimiter()
	client := apiclient.New(key, httpClient, limiter)
	ctx := context.Background()

	res, err := client.GetMatchlist(ctx, region.NA1, 237254272, nil)
	if err != nil {
		log.Fatal(err)

	}
	js, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(js))
}
