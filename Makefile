export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GOFLAGS := -mod=mod

all: install test lint clean

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock/v3.Tester -o ./tests
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_mock.go

#lint:
#	gometalinter ./tests/ -I minimock -e gopathwalk --disable=gotype --deadline=2m
#

./bin:
	mkdir ./bin

lint: install-tools
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./tests/...

install:
	go mod download
	go install ./cmd/minimock

# iterate over requirements from tools/tools.go and install them to ./bin
install-tools: ./bin
	@cd tools && go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install

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

release: install-tools
	./bin/goreleaser release

build: install-tools
	./bin/goreleaser build --snapshot --rm-dist
