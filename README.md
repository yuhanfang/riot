This project includes the following:
  - Threadsafe, rate-limited API client. See `examples/example_apiclient`
  - Centralized rate-limiting service. See below.
  - Cached client with a Google Cloud backend. See `examples/example_cachedclient`

In development are the following:
  - Champion select solver. Given a utility function for evaluating a team
    composition, solve for optimal pick/ban either as of the beginning of the
    game, or as of any node given any current state. See `analytics/champion_select/solver`

# Centralized rate limit service

A common problem for large-scale scraping of the Riot API is synchronization of
limits and quota across multiple clients. This package provides a ratelimit
server that centralizes quota management. Here is an example of how to set up
the client to point to the rate limiting service:

```go
server, err := url.Parse("http://your-server-here:1234")
limiter := client.New(http.DefaultClient, server)
client := apiclient.New("MY_API_KEY", http.DefaultClient, limiter)
```

Separately, you need to start the server listening on that port:

```bash
ratelimit_service --port=1234
```

The rate limit server is a simple RESTful service, and isn't limited to Go.
See the service documentation for details on the REST API.

# Project Goals

There are plenty of Riot API clients out there. I've even taken several stabs
at one myself already. So why make another one? Here is the experience I want to
build towards:

```bash
# What data do I want?
league_data --interactive --output_schema=schema.json

# Pull some data into a local CSV for small-scale testing. Use a local BoltDB
# to track what data has already been downloaded so that we don't repeat it.
league_data --input_schema=schema.json --max_mb=100 --csv_sink=test.csv --bolt_state=test.db 

# Continuously pull the items you want from the firehose. This time, throw it
# into Google BigQuery, using Google Datastore for state persistence.
league_data --continuous --schema=schema.json --bigquery_sink=my-dataset --datastore_sink=my-state
```

By abstracting away the whole data curation pipeline into a one-liner, I hope
this project can open the gates to the many data scientists out there who might
want to look at the game if there are no barriers to entry, but dread the
plumbing that is currently required to get something nice, structured, and
auto-updating.

Pull requests are always welcome.
