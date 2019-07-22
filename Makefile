NAME=beemo

.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -ldflags="-s -w" -i -o ./cmd/${NAME}

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean -i && rm cmd/beemo
