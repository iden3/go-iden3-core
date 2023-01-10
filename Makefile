#!/bin/bash

test:
	go test -count=1  -timeout=60s ./...

lint:
	golangci-lint run