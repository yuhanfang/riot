// Package cache saves data from the esports API into BigQuery.
//
// Most users should use the update_esports_cache command instead of calling
// this API directly.
package cache

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/esports"
	"github.com/yuhanfang/riot/esports/league"
)

// Cache represents a copy of the esports API data stored on Google Cloud. It
// is illegal to construct an instance directly. Use New to return a valid
// instance.
type Cache struct {
	// EntityNamespace is the namespace used when saving datastore entities. By
	// default, this is constructed with the empty string, which represents the
	// default namespace. The parameter can be set immediately after calling New.
	EntityNamespace string

	es         *esports.Client
	bq         *bigquery.Uploader
	ds         *datastore.Client
	doneEntity string
}

// New returns a Cache. doneEntity is the datastore entity name used to track
// whether a game was cached.
func New(es *esports.Client, bq *bigquery.Uploader, ds *datastore.Client, doneEntity string) *Cache {
	return &Cache{
		es:         es,
		bq:         bq,
		ds:         ds,
		doneEntity: doneEntity,
	}
}

// MaybeCreateDatasetAndTable attempts to create the dataset containing the
// given table, as well as the table itself.
func (c Cache) MaybeCreateDatasetAndTable(ctx context.Context, table *bigquery.Table) error {
	schema, err := bigquery.InferSchema(Game{})
	if err != nil {
		return err
	}
	client, err := bigquery.NewClient(ctx, table.ProjectID)
	if err != nil {
		return err
	}
	ds := client.Dataset(table.DatasetID)

	// Silently fail if dataset exists, or if creation fails, since on true
	// error, the table creation will also fail.
	ds.Create(ctx, nil)

	return table.Create(ctx, &bigquery.TableMetadata{
		Schema: schema,
	})
}

// UploadNewGames queries for all games in the given league and tournament
// title and uploads new games to BigQuery using the configured uploader. The
// table must exist. Call MaybeCreateDatasetAndTable to attempt the table creation.
func (c Cache) UploadNewGames(ctx context.Context, league league.League, tournament string) error {
	schema, err := bigquery.InferSchema(Game{})
	if err != nil {
		return err
	}
	leagues, err := c.es.GetLeagues(ctx, league)
	if err != nil {
		return err
	}
	if len(leagues.Leagues) != 1 {
		return fmt.Errorf("expected exactly one league for %s; got %d", league, len(leagues.Leagues))
	}
	leag := leagues.Leagues[0]
	t, ok := findTournament(leagues.HighlanderTournaments, tournament)
	if !ok {
		return fmt.Errorf("league %s does not have tournament %s", league, tournament)
	}
	for _, bracket := range t.Brackets {
		for _, match := range bracket.Matches {
			newGames, err := c.newGames(ctx, leag, t, bracket, match)
			if err != nil {
				return err
			}
			if len(newGames) == 0 {
				continue
			}
			mapping, err := c.es.GetHighlanderMatchDetails(ctx, t.ID, match.ID)
			if err != nil {
				return err
			}
			for _, game := range newGames {
				g, err := c.fetchGame(ctx, bracket, match, game, mapping)
				if err != nil {
					return err
				}
				g.LeagueName = leag.Name
				g.TournamentTitle = tournament
				err = c.putGame(ctx, schema, g, c.gameKey(leag, t, bracket, match, game))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c Cache) putGame(ctx context.Context, schema bigquery.Schema, g *Game, key datastore.Key) error {
	err := c.bq.Put(ctx, &bigquery.StructSaver{
		Struct:   g,
		Schema:   schema,
		InsertID: key.Name,
	})
	if err != nil {
		return err
	}
	_, err = c.ds.Put(ctx, &key, &struct{}{})
	return err
}

// Game is a serialized version of a specific game.
type Game struct {
	// LeagueName is the league in which the game took place.
	LeagueName string

	// TournamentTitle is the title of the tournament where the game took place.
	TournamentTitle string

	// BracketName is the name of the bracket where the game took place.
	BracketName string

	// MatchName is the name of the match where the game took place. For example,
	// a match may consists of a best of five between two teams.
	MatchName string

	// Game is the summary view of the game.
	Game *apiclient.Match

	// GameTimeline is the timeline view of the game.
	GameTimeline *apiclient.MatchTimeline
}

func (c Cache) fetchGame(ctx context.Context, bracket esports.Leagues_HighlanderTournament_Bracket, match esports.Leagues_HighlanderTournament_Bracket_Match, game esports.Leagues_HighlanderTournament_Bracket_Match_Game, mapping *esports.HighlanderMatchDetails) (*Game, error) {
	gameID := game.GameID
	region := game.GameRealm
	gameHash, ok := findGameHash(mapping, game.ID)
	if !ok {
		return nil, fmt.Errorf("could not find game ID %s in game hash mapping", game.ID)
	}
	stats, err := c.es.GetGameStats(ctx, region, gameID.Int64(), gameHash)
	if err != nil {
		return nil, err
	}
	timeline, err := c.es.GetGameTimeline(ctx, region, gameID.Int64(), gameHash)
	if err != nil {
		return nil, err
	}
	return &Game{
		BracketName:  bracket.Name,
		MatchName:    match.Name,
		Game:         stats,
		GameTimeline: timeline,
	}, nil
}

func (c Cache) gameKey(league esports.Leagues_League, tournament esports.Leagues_HighlanderTournament, bracket esports.Leagues_HighlanderTournament_Bracket, match esports.Leagues_HighlanderTournament_Bracket_Match, game esports.Leagues_HighlanderTournament_Bracket_Match_Game) datastore.Key {
	return datastore.Key{
		Kind:      c.doneEntity,
		Namespace: c.EntityNamespace,
		Name:      strconv.FormatInt(league.ID, 10) + ":" + tournament.ID + ":" + bracket.ID + ":" + match.ID + ":" + game.ID,
	}
}

func (c Cache) newGames(ctx context.Context, league esports.Leagues_League, tournament esports.Leagues_HighlanderTournament, bracket esports.Leagues_HighlanderTournament_Bracket, match esports.Leagues_HighlanderTournament_Bracket_Match) ([]esports.Leagues_HighlanderTournament_Bracket_Match_Game, error) {
	var games []esports.Leagues_HighlanderTournament_Bracket_Match_Game
	for _, game := range match.Games {
		if game.GameID == 0 {
			continue
		}
		key := c.gameKey(league, tournament, bracket, match, game)
		q := datastore.NewQuery(c.doneEntity).Namespace(c.EntityNamespace).Filter("__key__ =", &key).KeysOnly()
		got, err := c.ds.GetAll(ctx, q, nil)
		if err != nil {
			return nil, err
		}
		if len(got) == 0 {
			games = append(games, game)
		}
	}
	return games, nil
}

func findGameHash(mapping *esports.HighlanderMatchDetails, id string) (string, bool) {
	for _, m := range mapping.GameIDMappings {
		if m.ID == id {
			return m.GameHash, true
		}
	}
	return "", false
}

func findTournament(tournaments []esports.Leagues_HighlanderTournament, title string) (esports.Leagues_HighlanderTournament, bool) {
	for _, t := range tournaments {
		if t.Title == title {
			return t, true
		}
	}
	return esports.Leagues_HighlanderTournament{}, false
}
