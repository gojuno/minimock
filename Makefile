all: install test lint clean

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock/pkg.Tester -o ./tests
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_mock.go

lint:
	gometalinter ./... -I minimock -e gopathwalk --disable=gotype --deadline=2m

install:
	go mod download
	go install ./cmd/minimock

clean:
	rm -Rf bin/ dist/

test: generate
	go test -race ./...
