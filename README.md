# go-whosonfirst-stats

Go tools for generating usage stats for Who's On First

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### wof-stats-data

_THIS IS WORK IN PROGRESS_

```
./bin/wof-stats-data -mode repo /usr/local/whosonfirst-data/whosonfirst-data | jq
{
  "stats": {
    "alt_files": 312438,
    "borough": 247,
    "campus": 16081,
    "continent": 8,
    "count": 1044399,
    "country": 229,
    "county": 47029,
    "dependency": 40,
    "disputed": 91,
    "empire": 12,
    "is_ceased_true": 29908,
    "is_ceased_unknown": 702053,
    "is_current_false": 50888,
    "is_current_true": 348742,
    "is_current_unknown": 332331,
    "is_deprecated_false": 711042,
    "is_deprecated_true": 20919,
    "is_superseded_false": 714979,
    "is_superseded_true": 16982,
    "is_superseding_false": 717293,
    "is_superseding_true": 14668,
    "localadmin": 189742,
    "locality": 344351,
    "macrocounty": 375,
    "macrohood": 1193,
    "macroregion": 112,
    "marinearea": 305,
    "microhood": 2023,
    "neighbourhood": 124760,
    "ocean": 7,
    "planet": 1,
    "region": 4979,
    "timezone": 376
  }
}

./bin/wof-stats-data -mode repo /usr/local/whosonfirst-data/whosonfirst-data-venue-us-ca | jq
{
  "stats": {
    "count": 1519389,
    "is_ceased_false": 2,
    "is_ceased_true": 89,
    "is_ceased_unknown": 1519298,
    "is_current_false": 12710,
    "is_current_true": 23489,
    "is_current_unknown": 1483190,
    "is_deprecated_false": 1506753,
    "is_deprecated_true": 12636,
    "is_superseded_false": 1506711,
    "is_superseded_true": 12678,
    "is_superseding_false": 1507015,
    "is_superseding_true": 12374,
    "venue": 1519389
  }
}
```

Some existential flags like `is_superseded_false` and  `is_deprecated_false` do not need to be included in this output. I've included both dumps here so I have something to work with addressing the issue.

### wof-stats-du

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