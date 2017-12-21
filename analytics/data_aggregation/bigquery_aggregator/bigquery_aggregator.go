// Package bigquery_aggregator uploads data to BigQuery. State is stored in
// Datastore to prevent duplicate uploads. The aggregator is threadsafe, using
// Datastore transactions to manage ownership of each uploaded item.
//
// This API is unstable.
package bigquery_aggregator

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"github.com/nu7hatch/gouuid"
	"github.com/yuhanfang/riot/analytics/data_aggregation"
	"github.com/yuhanfang/riot/constants/region"
)

// Aggregator aggregates data into a BigQuery table. It is illegal to construct
// an instance directly. Use the New constructor instead.
type Aggregator struct {
	ns      string
	dataset string
	table   string
	ds      *datastore.Client
	bq      *bigquery.Client
}

// New returns an aggregator configured with the given namespace and datastore
// backend. The datastore is used as a lockservice for preventing duplicate
// rows. If namespace is empty, then the global namespace will be used.
func New(namespace, dataset, table string, ds *datastore.Client, bq *bigquery.Client) *Aggregator {
	return &Aggregator{
		ns:      namespace,
		dataset: dataset,
		table:   table,
		ds:      ds,
		bq:      bq,
	}
}

// lockvalue is a token stored in Datastore that is used as a proof of
// ownership. If a client supplies an identical UUID as the remote token,
// it knows that it owns the key.
type lockvalue struct {
	Token string
}

// acquireLock returns a string token if the key is acquired by this call, or
// empty if it is already owned by somebody else.
func (a *Aggregator) acquireLock(ctx context.Context, key *datastore.Key) (string, error) {
	uid, err := uuid.NewV4()
	tok := uid.String()
	if err != nil {
		return "", err
	}
	tx, err := a.ds.NewTransaction(ctx)
	if err != nil {
		return "", err
	}

	var current lockvalue
	// Continue until commit succeeds or we have an error.
	for {
		err = tx.Get(key, &current)
		// Already owned.
		if err == nil {
			tx.Rollback()
			return "", nil
		}

		// Actual error.
		if err != datastore.ErrNoSuchEntity {
			tx.Rollback()
			return "", err
		}

		current.Token = tok
		_, err = tx.Put(key, &current)
		// Actual error.
		if err != nil {
			tx.Rollback()
			return "", err
		}

		_, err = tx.Commit()
		// Owned by us.
		if err == nil {
			return tok, nil
		}
		// Actual error.
		if err != datastore.ErrConcurrentTransaction {
			return "", err
		}
	}
}

// releaseLock returns true if the key was released by this call, or false if
// the key is already released.
func (a *Aggregator) releaseLock(ctx context.Context, key *datastore.Key, tok string) (bool, error) {
	tx, err := a.ds.NewTransaction(ctx)
	if err != nil {
		return false, err
	}

	// Continue until commit succeeds or we have an error.
	var current lockvalue
	for {
		err = tx.Get(key, &current)
		// Not owned.
		if err == datastore.ErrNoSuchEntity {
			tx.Rollback()
			return false, nil
		}
		// Actual error.
		if err != nil {
			tx.Rollback()
			return false, err
		}
		// Not owned by us.
		if current.Token != tok {
			tx.Rollback()
			return false, nil
		}

		err = tx.Delete(key)
		// Actual error.
		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = tx.Commit()
		// Released by us.
		if err == nil {
			return true, nil
		}
		// Actual error.
		if err != datastore.ErrConcurrentTransaction {
			return false, err
		}
	}
}

// markDone returns true if the key was marked done by this call, or false if
// the key is not owned.
func (a *Aggregator) markDone(ctx context.Context, key *datastore.Key, tok string) (bool, error) {
	tx, err := a.ds.NewTransaction(ctx)
	if err != nil {
		return false, err
	}

	// Continue until commit succeeds or we have an error.
	var current lockvalue
	for {
		err = tx.Get(key, &current)
		// Not owned.
		if err == datastore.ErrNoSuchEntity {
			tx.Rollback()
			return false, nil
		}
		// Actual error.
		if err != nil {
			tx.Rollback()
			return false, err
		}
		// Not owned by us.
		if current.Token != tok {
			tx.Rollback()
			return false, nil
		}

		current.Token = "done"
		_, err = tx.Put(key, &current)
		// Actual error.
		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = tx.Commit()
		// Released by us.
		if err == nil {
			return true, nil
		}
		// Actual error.
		if err != datastore.ErrConcurrentTransaction {
			return false, err
		}
	}
}

// MatchExists returns true if the match ID for the given region is already
// stored.
func (a *Aggregator) MatchExists(ctx context.Context, r region.Region, id int64) (bool, error) {
	var current lockvalue
	key := a.gameKey(r, id)
	err := a.ds.Get(ctx, key, &current)
	if err == nil {
		return true, nil
	}
	if err == datastore.ErrNoSuchEntity {
		return false, nil
	}
	return false, err
}

func (a *Aggregator) gameKey(r region.Region, id int64) *datastore.Key {
	return &datastore.Key{
		Kind:      fmt.Sprintf("aggregator-save-match-%s:%s:%s", r, a.dataset, a.table),
		ID:        id,
		Namespace: a.ns,
	}
}

// uploadKeyValue wraps a match that corresponds to a given key.
type uploadKeyValue struct {
	Key   *datastore.Key
	Value data_aggregation.Match
}

// SaveMatches stores the match from the given region. Returns whether the
// function stored the match. If not, and there was no error, then the match
// was already cached.
func (a *Aggregator) SaveMatches(ctx context.Context, matches []data_aggregation.Match) error {
	ds := a.bq.Dataset(a.dataset)
	ds.Create(ctx, nil)
	tab := ds.Table(a.table)
	schema, err := bigquery.InferSchema(data_aggregation.Match{})
	if err != nil {
		return err
	}
	tab.Create(ctx, &bigquery.TableMetadata{Schema: schema})
	u := tab.Uploader()

	var (
		acquired   = make(map[*datastore.Key]string)
		toUpload   []*bigquery.StructSaver
		acquireErr error
	)

	for _, match := range matches {
		key := a.gameKey(match.Region, match.ID)
		tok, err := a.acquireLock(ctx, key)
		if err == nil && tok != "" {
			toUpload = append(toUpload, &bigquery.StructSaver{
				Struct:   match,
				Schema:   schema,
				InsertID: fmt.Sprintf("%s:%d:%s:%s", key.Kind, key.ID, key.Name, key.Namespace),
			})
			acquired[key] = tok
		}
		if err != nil && acquireErr == nil {
			acquireErr = err
		}
	}

	putErr := u.Put(ctx, toUpload)

	if putErr == nil {
		// Mark each acquired key as done.
		for key, tok := range acquired {
			for {
				_, err := a.markDone(ctx, key, tok)
				if err == nil {
					break
				}
				select {
				case <-ctx.Done():
					return err
				default:
				}
			}
		}
		return acquireErr
	}

	for key, tok := range acquired {
		// If the match can't be saved, then the lock must be released so that others
		// can try again. If the context is done, the lock release could fail on
		// datastore problems. We should address this possibility in the future.
		for {
			_, err := a.releaseLock(ctx, key, tok)
			if err == nil {
				break
			}
			select {
			case <-ctx.Done():
				return err
			default:
			}
		}
	}
	return putErr
}
