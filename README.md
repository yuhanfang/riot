# Riot
[![Build
Status](https://travis-ci.org/yuhanfang/riot.svg?branch=master)](https://travis-ci.org/yuhanfang/riot)
[![](https://godoc.org/github.com/yuhanfang/riot?status.svg)](http://godoc.org/github.com/yuhanfang/riot)

This project aims to provide a batteries-included toolkit for all publicly
available Riot data written in Go.

# Features
  - Threadsafe, rate-limited API client. See `examples/example_apiclient`
  - Centralized rate-limiting service for multi-server configurations. See
    https://github.com/yuhanfang/riot/wiki/Rate-Limit-Service
  - Cached client built on top of a Google Cloud backend. See
    `examples/example_cachedclient`
  - Static data client backed by data dragons. See `examples/example_staticdata`
  - Uploader to structure API data in BigQuery for easy analysis. See
    `examples/example_bigquery_aggregator`
  - Competitive esports data API. See `examples/example_esports`

# Dependencies
Riot uses Go Modules to manage dependencies. Installing and updating packages are just like you would without using a dependency manager.

```
go get [package]
```



See the wiki entry [Modules](https://github.com/golang/go/wiki/Module) for more on Modules

# Contributing
Please see CONTRIBUTING if you are interested in contributing features or
bugfixes.

Feel free to raise an issue with label "question" or DM me at
[@yuhan_fang](https://www.twitter.com/yuhan_fang) if you have any questions.
