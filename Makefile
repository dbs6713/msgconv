# Author: Don B. Stringham <donbstringham@gmail.com>

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=$(GOCMD) fmt

# Application parameters
AUTHOR=donbstringham
BUILD=`git rev-list --count --first-parent HEAD`
BUILD_TIME=`date +%FT%T%z`
DOMAIN=github.com
LDFLAGS=-ldflags "-s -w -X $(DOMAIN)/$(NAME)/core.Build=$(BUILD) -X $(DOMAIN)/$(NAME)/core.BuildTime=$(BUILD_TIME) -X $(DOMAIN)/$(NAME)/core.Version=$(VERSION_NUM) -X $(DOMAIN)/$(NAME)/core.Name=$(NAME)"
NAME=httpmessageconverter
TAG=$(VERSION_NUM)+$(BUILD)
VERSION_FILE=.version
VERSION_FILES=$(sort $(wildcard VERSION*))
VERSION_NUM=$$(cat $(VERSION_FILE))
VERSION_NUM_PREV=$$(cat $(VERSION_FILE))

# List building
ALL_LIST=./...
DIR_LIST=`go list ./... | grep -v /vendor/`

# Targets
all: clean test

clean:
	@go clean ./...

gkgo:
	@$(GOGINKGO)

test:
	@time go test $(DIR_LIST) -cover

ver:
	@echo $(TAG)
