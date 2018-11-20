CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test ! -d src; then mkdir src; fi
	mkdir -p src
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-index"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-uri"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-cli"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/pretty"
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/whosonfirst/warning src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/tidwall/gjson src/github.com/tidwall/
	mv src/github.com/whosonfirst/go-whosonfirst-geojson-v2/vendor/github.com/tidwall/match src/github.com/tidwall/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go

bin: 	self
	@GOPATH=$(GOPATH) go build -o bin/wof-stats-du cmd/wof-stats-du.go
	@GOPATH=$(GOPATH) go build -o bin/wof-stats-counts cmd/wof-stats-counts.go

# this is left here as a reference but will otherwise fail because of the sqlite 
# dependency (in go-whosonfirst-index) which is OS-specific (20181119/thisisaaronland)

dist-build:
	OS=darwin make dist-os
	OS=windows make dist-os
	OS=linux make dist-os

dist-os:
	mkdir -p dist/$(OS)
	GOOS=$(OS) GOPATH=$(GOPATH) GOARCH=386 go build -o dist/$(OS)/wof-stats-du cmd/wof-stats-du.go
	GOOS=$(OS) GOPATH=$(GOPATH) GOARCH=386 go build -o dist/$(OS)/wof-stats-count cmd/wof-stats-count.go

rmdist:
	if test -d dist; then rm -rf dist; fi
