#!/usr/bin/env bash
#
# run golangci-lint

golangci-lint -E misspell,gocyclo,dupl,gofmt,golint,unconvert,goimports,depguard,gocritic,interfacer run --disable-all

# consider adding funlen
