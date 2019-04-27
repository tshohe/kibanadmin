GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: clean setup build
setup:
	$(GOGET) github.com/bitly/go-simplejson
build:
	export GO111MODULE=on
	$(GOBUILD) -o bin/kibana-refresh -v
# test:
# 	$(GOTEST) -v ./test/
clean:
	$(GOCLEAN)
	rm -f bin/*
