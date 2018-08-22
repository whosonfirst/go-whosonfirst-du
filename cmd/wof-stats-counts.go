package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/pretty"	
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/warning"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func main() {

	valid_modes := strings.Join(index.Modes(), ",")
	desc_modes := fmt.Sprintf("The mode to use importing data. Valid modes are: %s.", valid_modes)

	var mode = flag.String("mode", "repo", desc_modes)
	var format = flag.String("format", "json", "Write stats in this format. Valid formats are: json, markdown.")
	var out = flag.String("out", "", "Write stats to this path. If empty write stats to STDOUT.")

	var json_pretty = flag.Bool("pretty", false, "Generate pretty-printed JSON.")

	flag.Parse()

	wr := os.Stdout

	if *out != "" {

		path, err := filepath.Abs(*out)

		if err != nil {
			log.Fatal(err)
		}

		fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatal(err)
		}

		wr = fh
	}

	mu := new(sync.Mutex)
	stats := make(map[string]int64)

	incr := func(key string) {

		mu.Lock()
		defer mu.Unlock()

		count, ok := stats[key]

		if !ok {
			count = 0
		}

		stats[key] = count + 1
	}

	incr_existential := func(key string, str_flag string) {

		switch key {

		case "is_current":

			switch str_flag {
			case "-1":
				return
			case "0":
				key = fmt.Sprintf("%s_false", key)
			default:
				key = fmt.Sprintf("%s", key)
			}

		default:

			if str_flag != "1" {
				return
			}
		}

		incr(key)
	}

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

		is_alt, err := uri.IsAltFile(path)

		if err != nil {
			return err
		}

		incr("count")

		if is_alt {
			incr("is_alt")
			return nil
		}

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil && !warning.IsWarning(err) {
			return err
		}

		pt := f.Placetype()
		incr(pt)

		is_current, err := whosonfirst.IsCurrent(f)

		if err != nil {
			return err
		}

		incr_existential("is_current", is_current.StringFlag())

		is_deprecated, err := whosonfirst.IsDeprecated(f)

		if err != nil {
			return err
		}

		incr_existential("is_deprecated", is_deprecated.StringFlag())

		is_ceased, err := whosonfirst.IsCeased(f)

		if err != nil {
			return err
		}

		incr_existential("is_ceased", is_ceased.StringFlag())

		is_superseded, err := whosonfirst.IsSuperseded(f)

		if err != nil {
			return err
		}

		incr_existential("is_superseded", is_superseded.StringFlag())

		is_superseding, err := whosonfirst.IsSuperseding(f)

		if err != nil {
			return err
		}

		incr_existential("is_superseding", is_superseding.StringFlag())

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
	report["stats"] = stats

	// please put this is a package or something...
	// (20180817/thisisaaronland)

	switch *format {

	case "markdown":

		for c, i := range report {

			details := i.(map[string]int64)
			keys := make([]string, 0)

			for k, _ := range details {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			wr.Write([]byte(fmt.Sprintf("## %s\n\n", c)))

			for _, k := range keys {

				wr.Write([]byte(fmt.Sprintf("* **%s** %d\n", k, details[k])))
			}

			wr.Write([]byte("\n"))
		}

	default:

		body, err := json.Marshal(report)
		
		if err != nil {
			log.Fatal(err)
		}

		if *json_pretty {
			body = pretty.Pretty(body)
		}
		
		r := bytes.NewReader(body)
		io.Copy(wr, r)
	}

	os.Exit(0)
}
