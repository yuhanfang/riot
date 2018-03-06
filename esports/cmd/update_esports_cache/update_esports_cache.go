// Command update_esports_cache creates and updates a cache of esports game
// data on Google Cloud.
//
// Flags:
//	project: Google Cloud project ID
//  dataset: BigQuery dataset ID
//  table: BigQuery table ID
//  entity: Datastore entity ID
//  namespace: Datastore entity namespace, or empty for default
//  league: league slug identifier. Use any of the slugs returned by Slug() in
//		esports/leagues/leagues.go. For example, "na-lcs" is the NA LCS league.
//  tournament: tournament title. Use any of the constants in
//		esports/tournaments/tournaments.go. For example, "na_2018_spring" is the
//		spring split of 2018 in the NA region.

package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"github.com/yuhanfang/riot/esports"
	"github.com/yuhanfang/riot/esports/cache"
	"github.com/yuhanfang/riot/esports/league"
	"github.com/yuhanfang/riot/esports/tournaments"
)

var (
	project    = flag.String("project", "", "existing Google Cloud project ID")
	dataset    = flag.String("dataset", "", "Google BigQuery dataset ID that will be created if missing")
	table      = flag.String("table", "", "Google BigQuery table ID that will be created if missing")
	entity     = flag.String("entity", "", "Google Datastore entity ID used to track cache status")
	namespace  = flag.String("namespace", "", "Google Datastore namespace ID to track cache status")
	leagueSlug = flag.String("league", league.LeagueNALCS.Slug(), "league slug to cache")
	tournament = flag.String("tournament", tournaments.NA2018Spring, "tournament title to cache")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	es := esports.NewClient(http.DefaultClient)
	bq, err := bigquery.NewClient(ctx, *project)
	if err != nil {
		log.Fatal(err)
	}
	ds, err := datastore.NewClient(ctx, *project)
	if err != nil {
		log.Fatal(err)
	}
	data := bq.Dataset(*dataset)
	tab := data.Table(*table)
	up := tab.Uploader()

	client := cache.New(es, up, ds, *entity)
	client.EntityNamespace = *namespace
	client.MaybeCreateDatasetAndTable(ctx, tab)

	leag := league.SlugToLeague(*leagueSlug)
	if leag == league.LeagueInvalid {
		log.Fatal("unsupported league: ", *leagueSlug)
	}
	err = client.UploadNewGames(ctx, leag, *tournament)
	if err != nil {
		log.Fatal(err)
	}
}
