package main

import (
       "bytes"
	"context"	
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/pretty"	
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {

	valid_modes := strings.Join(index.Modes(), ",")
	desc_modes := fmt.Sprintf("The mode to use importing data. Valid modes are: %s.", valid_modes)

	var mode = flag.String("mode", "repo", desc_modes)
	var json_pretty = flag.Bool("pretty", false, "Generate pretty-printed JSON.")
	
	flag.Parse()

	mu := new(sync.Mutex)

	totals := make(map[string]int64)
	stats := make(map[string]int64)

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		is_wof, err := uri.IsWOFFile(path)

		if err != nil {
			return err
		}

		if !is_wof {
			return nil
		}

		info, err := os.Stat(path)

		if err != nil {
			return err
		}

		sz_b := info.Size()
		sz_kb := float64(sz_b) / 1024.

		var k string

		if sz_kb <= 10.0 {
			k = "0-10k"
		} else if sz_kb <= 100.0 {
			k = "10-100k"
		} else if sz_kb <= 500.0 {
			k = "100-500k"
		} else if sz_kb <= 1024.0 {
			k = "500k-1m"
		} else if sz_kb <= 10024.0 {
			k = "1-10m"
		} else if sz_kb <= 100024.0 {
			k = "10-100m"
		} else {
			k = "100m-BIG"
		}

		mu.Lock()
		defer mu.Unlock()

		_, ok := stats[k]

		if ok {
			stats[k] += 1
		} else {
			stats[k] = 1
		}

		_, ok = totals["files"]

		if ok {
			totals["files"] += 1
		} else {
			totals["files"] = 1
		}

		_, ok = totals["bytes"]

		if ok {
			totals["bytes"] += sz_b
		} else {
			totals["bytes"] = sz_b
		}

		return nil
	}

	idx, err := index.NewIndexer(*mode, cb)

	if err != nil {
		log.Fatal(err)
	}

	sources := flag.Args()
	err = idx.IndexPaths(sources)

	if err != nil {
		log.Fatal(err)
	}

	report := make(map[string]interface{})
	report["totals"] = totals
	report["stats"] = stats

	body, err := json.Marshal(report)

	if err != nil {
		log.Fatal(err)
	}

	if *json_pretty {
		body = pretty.Pretty(body)
	}

	wr := os.Stdout	// in advance of custom options
	
	r := bytes.NewReader(body)
	io.Copy(wr, r)

	os.Exit(0)
}
