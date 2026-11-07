.PHONY: build test run fmt

build:
	go build -o bin/jurassic-balloon-term ./cmd/jurassic-balloon-term

test:
	go test ./...

run:
	go run ./cmd/jurassic-balloon-term

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

