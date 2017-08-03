PKG=github.com/gojuno/minimock

all: test

generate:
	go run ./cmd/minimock/minimock.go -f ${PKG} -i Tester -o ./tests/tester_mock_test.go -p tests
	go run ./cmd/minimock/minimock.go -f ${PKG}/tests -i Stringer -o ./tests/stringer_mock.go -p tests

test: generate
	go test ${PKG}/tests
