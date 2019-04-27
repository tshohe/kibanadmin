.PHONY: build clean

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

build: clean
	export GO111MODULE=on
	$(GOBUILD) -o bin/kibana-refresh cmd/kibana-refresh/main.go

# test:
# 	$(GOTEST) -v ./test/

clean:
	rm -f bin/*
