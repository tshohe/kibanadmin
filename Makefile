GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: setup test build
setup:
	$(GOGET) github.com/bitly/go-simplejson
build:
	$(GOBUILD) -o bin/kibana-refresh -v
# test:
# 	$(GOTEST) -v ./test/
clean:
	$(GOCLEAN)
	rm -f bin/*
