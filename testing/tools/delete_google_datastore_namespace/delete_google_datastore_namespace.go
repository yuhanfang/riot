// Deletes a namespace from Google datastore.
//
// Batch deletes all keys that have the given namespace.
//
// Usage:
//   delete_google_datastore_namespace --namespace MyNamespace
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"cloud.google.com/go/datastore"
)

var (
	namespace = flag.String("namespace", "", "namespace to delete")
)

func main() {
	flag.Parse()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	ctx := context.Background()
	ds, err := datastore.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}
	query := datastore.NewQuery("").Namespace(*namespace).KeysOnly()
	keys, err := ds.GetAll(ctx, query, nil)
	batchSize := 500
	for i := 0; i < len(keys); i += batchSize {
		max := i + batchSize
		if max > len(keys) {
			max = len(keys)
		}
		err = ds.DeleteMulti(ctx, keys[i:max])
		if err != nil {
			log.Fatal(err)
		}
	}
}
