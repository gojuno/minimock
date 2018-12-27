all: install test lint

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock.Tester -o ./tests
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_mock.go

lint:
	#golangci-lint run ./...
	gometalinter ./... -I minimock -e vendor

install:
	go install ./cmd/minimock

test: generate
	go test -race ./...
