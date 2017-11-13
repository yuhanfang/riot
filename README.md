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
