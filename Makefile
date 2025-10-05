.PHONY: test test-verbose test-coverage test-integration run build clean

test:
	go test ./...

test-verbose:
	go test -v ./...

test-integration:
	INTEGRATION_TESTS=1 go test -v ./internal/clients/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

run:
	go run main.go

build:
	go build -o sportsagent

clean:
	rm -f sportsagent coverage.out coverage.html
