export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GOFLAGS := -mod=mod

all: install test lint

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock/v3.Tester -o ./tests
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_mock.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericInout -o ./tests/generic/generic_inout.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericOut -o ./tests/generic/generic_out.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericIn -o ./tests/generic/generic_in.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericSpecific -o ./tests/generic/generic_specific.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericSimpleUnion -o ./tests/generic/generic_simple_union.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericComplexUnion -o ./tests/generic/generic_complex_union.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericInlineUnion -o ./tests/generic/generic_inline_union.go

./bin:
	mkdir ./bin

lint: ./bin/golangci-lint
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./tests/...

install:
	go mod download
	go install ./cmd/minimock

./bin/golangci-lint: ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.51.1

./bin/goreleaser: ./bin
	go install -modfile tools/go.mod github.com/goreleaser/goreleaser

./bin/minimock:
	go build ./cmd/minimock -o ./bin/minimock

test:
	go test $(go list ./... | grep -v /snapshots)

release: ./bin/goreleaser
	./bin/goreleaser release

build: ./bin/goreleaser
	./bin/goreleaser build --snapshot --rm-dist
