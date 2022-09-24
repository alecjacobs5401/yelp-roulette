.PHONY: build
build:
	@go build -o bin/yelp-roulette ./cmd/cli

.PHONY: build-server
build-server:
	@go build -o bin/server .
