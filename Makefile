PKG=github.com/gojuno/minimock

all: test

generate:
	go run minimock.go -f ${PKG}/tests -i Tester -o ./tests/tester_mock_test.go -p tests
	go run minimock.go -f ${PKG}/tests -i Stringer -o ./tests/stringer_mock_test.go -p tests -testingType *TesterMock

test: generate
	go test ${PKG}/tests
