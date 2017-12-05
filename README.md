Documentation: https://godoc.org/github.com/yuhanfang/riot

This project aims to provide a batteries-included toolkit for all publicly
available Riot data. The main features include:

  - Threadsafe, rate-limited API client. See `examples/example_apiclient`
  - Centralized rate-limiting service for multi-server configurations. See
    https://github.com/yuhanfang/riot/wiki/Rate-Limit-Service
  - Cached client built on top of a Google Cloud backend. See
    `examples/example_cachedclient`
  - Static data client backed by data dragons. See `examples/example_staticdata`
  - Uploader to structure API data in BigQuery for easy analysis. See
    `examples/example_bigquery_aggregator`
  - Competitive esports data API. See `examples/example_esports`

Please see CONTRIBUTING if you are interested in contributing features or
bugfixes.

Feel free to raise an issue with label "question" or DM me at
[@iso646](https://www.twitter.com/iso646) if you have any questions.

[![Build
Status](https://travis-ci.org/yuhanfang/riot.svg?branch=master)](https://travis-ci.org/yuhanfang/riot)
