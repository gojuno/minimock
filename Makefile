export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GOFLAGS := -mod=mod

all: install test lint clean

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock/v3.Tester -o ./tests
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_mock.go

./bin:
	mkdir ./bin

lint: ./bin/golangci-lint
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./tests/...

install:
	go mod download
	go install ./cmd/minimock

./bin/golangci-lint: ./bin
	@cd tools && go get github.com/golangci/golangci-lint/cmd/golangci-lint

./bin/goreleaser: ./bin
	@cd tools && go install github.com/goreleaser/goreleaser

./bin/minimock:
	go build ./cmd/minimock -o ./bin/minimock

clean:
	[ -e ./tests/formatter_mock.go.test_origin ] && mv -f ./tests/formatter_mock.go.test_origin ./tests/formatter_mock.go
	[ -e ./tests/tester_mock_test.go.test_origin ] && mv -f ./tests/tester_mock_test.go.test_origin ./tests/tester_mock_test.go
	rm -Rf bin/ dist/

test_save_origin:
	[ -e ./tests/formatter_mock.go.test_origin ] || cp ./tests/formatter_mock.go ./tests/formatter_mock.go.test_origin
	[ -e ./tests/tester_mock_test.go.test_origin ] || cp ./tests/tester_mock_test.go ./tests/tester_mock_test.go.test_origin

test: test_save_origin generate
	diff ./tests/formatter_mock.go ./tests/formatter_mock.go.test_origin
	diff ./tests/tester_mock_test.go ./tests/tester_mock_test.go.test_origin
	go test -race ./...

release: ./bin/goreleaser
	./bin/goreleaser release

build: ./bin/goreleaser
	./bin/goreleaser build --snapshot --rm-dist
