# go-whosonfirst-stats

Go tools for generating statistics for Who's On First documents.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.10](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### wof-stats-data

```
./bin/wof-stats-data -h
Usage of ./bin/wof-stats-data:
  -format string
    	Write stats in this format. Valid formats are: json, markdown. (default "json")
  -mode string
    	The mode to use importing data. Valid modes are: directory,feature,feature-collection,files,geojson-ls,meta,path,repo,sqlite. (default "repo")
  -out string
    	Write stats to this path. If empty write stats to STDOUT.
```

For example:

```
./bin/wof-stats-data /usr/local/whosonfirst-data/whosonfirst-data-venue-us-ca | jq
{
  "stats": {
    "count": 1519389,
    "is_ceased": 89,
    "is_current": 23489,
    "is_current_false": 12710,
    "is_deprecated": 12636,
    "is_superseded": 12678,
    "is_superseding": 12374,
    "venue": 1519389
  }
}
```

### wof-stats-du

```
./bin/wof-stats-du -h
Usage of ./bin/wof-stats-du:
  -mode string
    	The mode to use importing data. Valid modes are: directory,feature,feature-collection,files,geojson-ls,meta,path,repo,sqlite. (default "repo")
```

For example:

```
./bin/wof-stats-du -mode repo /usr/local/data/whosonfirst-data* | python -mjson.tool
{   
    "stats": {
        "0-10k": 26336534,
        "1-10m": 428,
        "10-100k": 170186,
        "10-100m": 30,
        "100-500k": 14563,
        "100m-BIG": 1,
        "500k-1m": 827
    },
    "totals": {
        "bytes": 58514102879,
        "files": 26522569
    }
}
```