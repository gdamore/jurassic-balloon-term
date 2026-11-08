.PHONY: build web test run fmt

build:
	go build -o bin/jurassic-balloon-term ./cmd/jurassic-balloon-term

web:
	./scripts/build-web.sh

test:
	go test ./...

run:
	go run ./cmd/jurassic-balloon-term

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

