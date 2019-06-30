# Indexer
An Indexer reads through rows in the stats databases and creates stats tables.
For example counting up kills and deaths. We use indexers so that we don't have to query potentially
1000s of rows for every request.

## Configuration

create or add to the existing collector settings.json
You can set the path to this file with the environment variable `KAGSTATS_CONFIG`
```
{
	"databaseConnect": "root:password@tcp(127.0.0.1:3306)/kagstats"
	"indexer": {
		"batchSize": 100,
		"interval": "30s"
	}
}
```

The above are default values for batchSize and interval.

Alternatively environment variables can be used to configure the indexer.
`INDEXER_DB` - the database connection string
`INDEXER_BATCHSIZE` - the number of rows to process per an interval
`INDEXER_INTERVAL` - the time between processing

## Creating New Indexers

You can create a new indexer using this package. The table creation, row tracking and
processing framework is all handled by the package function Process. All you need to do
is implement KillsIndexer and call Run. See the basic indexer for a full example.
