This project includes the following:
  - Threadsafe, rate-limited API client. See `examples/example_apiclient`
  - Centralized rate-limiting service. See below.
  - Cached client with a Google Cloud backend. See `examples/example_cachedclient`

In development are the following:
  - Champion select solver. Given a utility function for evaluating a team
    composition, solve for optimal pick/ban either as of the beginning of the
    game, or as of any node given any current state. See `analytics/champion_select/solver`

If any of this is interesting and you want to help out, I'm always happy to chat! https://discordapp.com/invite/8XPtaFB

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

