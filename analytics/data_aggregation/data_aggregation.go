// Package data_aggregation pulls data from a Riot API client and stores it in
// a centralized table.
package data_aggregation

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/queue"
	"github.com/yuhanfang/riot/constants/region"
)

// Match is the serialized format of match data. Match or Timeline may be nil
// if the data is unavailable from the API. Matches are unique by ID and
// Region.
type Match struct {
	ID       int64
	Region   region.Region
	Match    *apiclient.Match
	Timeline *apiclient.MatchTimeline
}

// Sink stores match data in an aggregate source.
type Sink interface {
	// MatchExists returns true if the match ID for the given region is already
	// stored in the sink.
	MatchExists(ctx context.Context, r region.Region, id int64) (bool, error)

	// SaveMatches saves the matches into the sink. It is not an error to save a
	// match that has already been saved. Sink must not duplicate the match, but
	// may either ignore or update the saved match.
	SaveMatches(ctx context.Context, m []Match) error
}

// Aggregator aggregates match data from a client into a sink. It is illegal to
// construct an instance directly. Use NewAggregator to return a valid
// instance.
type Aggregator struct {
	client apiclient.Client
	sink   Sink
}

// NewAggregator returns an Aggregator configured to query the Riot API with
// the given client, and aggregate results into the given Sink.
func NewAggregator(c apiclient.Client, s Sink) *Aggregator {
	return &Aggregator{
		client: c,
		sink:   s,
	}
}

// AggregateChallengerLeagueMatches aggregates all games from accounts
// currently in challenger league in the given region in queue. Only games
// since the given begin time are considered. The zero time value indicates
// that all games should be considered.
func (a Aggregator) AggregateChallengerLeagueMatches(ctx context.Context, r region.Region, q queue.Queue, since time.Time) error {
	league, err := a.client.GetChallengerLeague(ctx, r, q)
	if err != nil {
		return err
	}
	accountIDs := a.getAccountIDsInLeague(ctx, r, league)
	matches := a.GetMatchIDsForAccounts(ctx, r, q, since, accountIDs)
	a.UploadMatches(ctx, r, matches)
	return nil
}

func (a Aggregator) GetMatchIDsForAccounts(ctx context.Context, r region.Region, q queue.Queue, since time.Time, accountIDs []int64) map[int64]struct{} {
	// Query recent matches for each account.
	opts := apiclient.GetMatchlistOptions{
		Queue:     []queue.Queue{q},
		BeginTime: &since,
	}
	matches := make(map[int64]struct{})
	matchIDs := make(chan int64)
	wg := sync.WaitGroup{}
	done := make(chan bool)

	for _, account := range accountIDs {
		wg.Add(1)
		account := account
		// Process each account concurrently.
		go func() {
			defer wg.Done()
			matchlist, err := a.client.GetMatchlist(ctx, r, account, &opts)
			if err == nil {
				var (
					subwg   sync.WaitGroup
					subdone = make(chan bool)
				)
				// Process each account match concurrently.
				for _, m := range matchlist.Matches {
					subwg.Add(1)
					m := m
					go func() {
						defer subwg.Done()
						exists, err := a.sink.MatchExists(ctx, r, m.GameID)
						if err != nil {
							log.Printf("MatchExists failed for region %s game %d: %v", r, m.GameID, err)
						}
						if err == nil && !exists {
							select {
							case matchIDs <- m.GameID:
							case <-ctx.Done():
								return
							}
						}
					}()

					// Block until cancelled or all matches processed.
					go func() {
						subwg.Wait()
						subdone <- true
					}()
					select {
					case <-ctx.Done():
					case <-subdone:
					}
				}
			} else {
				log.Printf("GetMatchlist failed for region %s account %d: %v", r, account, err)
			}
		}()
	}

	// Block until cancelled or all matches processed.
	go func() {
		wg.Wait()
		done <- true
	}()
	more := true
	for more {
		select {
		case <-ctx.Done():
			more = false
		case match := <-matchIDs:
			matches[match] = struct{}{}
		case <-done:
			more = false
		}
	}

	return matches
}

func (a Aggregator) UploadMatches(ctx context.Context, r region.Region, matches map[int64]struct{}) {
	// Retrieve match and timeline data for each match.
	wg := sync.WaitGroup{}
	for m := range matches {
		match, err := a.client.GetMatch(ctx, r, m)
		if err != nil && err != apiclient.ErrDataNotFound {
			log.Printf("GetMatch failed for region %s game %d: %v", r, m, err)
			continue
		}
		timeline, err := a.client.GetMatchTimeline(ctx, r, m)
		if err != nil && err != apiclient.ErrDataNotFound {
			log.Printf("GetMatchTimeline failed for region %s game %d: %v", r, m, err)
			continue
		}
		upload := Match{
			ID:       m,
			Region:   r,
			Match:    match,
			Timeline: timeline,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := a.sink.SaveMatches(ctx, []Match{upload})
			if err != nil {
				log.Printf("SaveMatches failed for region %s game %d: %v", upload.Region, upload.ID, err)
			}
		}()
	}
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}
}
func (a Aggregator) getAccountIDsInLeague(ctx context.Context, r region.Region, league *apiclient.LeagueList) []int64 {
	var (
		accounts   = make(chan int64)
		accountIDs []int64
		wg         sync.WaitGroup
		done       = make(chan bool)
	)
	for _, entry := range league.Entries {
		entry := entry
		wg.Add(1)
		go func() {
			defer wg.Done()
			parsed, err := strconv.ParseInt(entry.PlayerOrTeamID, 10, 64)
			if err != nil {
				log.Printf("ParseInt failed for region %s player %s in league %s: %v", r, entry.PlayerOrTeamID, league.LeagueID, err)
				return
			}
			summoner, err := a.client.GetBySummonerID(ctx, r, parsed)
			if err != nil {
				log.Printf("GetBySummonerID failed for region %s summoner %d: %v", r, parsed, err)
				return
			}
			select {
			case <-ctx.Done():
			case accounts <- summoner.AccountID:
			}
		}()
	}

	go func() {
		wg.Wait()
		done <- true
	}()

	more := true
	for more {
		select {
		case <-ctx.Done():
			more = false
		case id := <-accounts:
			accountIDs = append(accountIDs, id)
		case <-done:
			more = false
		}
	}
	return accountIDs
}
