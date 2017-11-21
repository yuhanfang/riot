// Demonstrates how to use a BigQuery aggregator to aggregate match data from
// the Riot API into a queryable table.
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"github.com/yuhanfang/riot/analytics/data_aggregation/bigquery_aggregator"
	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

const game = 2644987649

func main() {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	key := os.Getenv("RIOT_APIKEY")
	ctx := context.Background()
	ds, err := datastore.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}
	bq, err := bigquery.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}
	agg := bigquery_aggregator.New("TestAggregator", ds, bq)
	limiter := ratelimit.NewLimiter()
	httpClient := http.DefaultClient
	client := apiclient.New(key, httpClient, limiter)

	match, err := client.GetMatch(ctx, region.NA1, game)
	if err != nil {
		log.Fatal(err)
	}
	ok, err := agg.SaveMatch(ctx, "TestAggregatorDataset", "matches", region.NA1, match)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		log.Println("match stored successfully")
	} else {
		log.Println("match already exists")
	}
}
