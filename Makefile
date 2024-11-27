.PHONY: server watch translations test cover

# Path to your server's main file
SERVER_PATH = cmd/server/main.go

# Default target to run the server
server:
	go run $(SERVER_PATH)

# Target to watch for changes in .go files and re-run the server
watch:
	find . -name '*.go' | entr -r make server

translations:
	go generate translations/translations.go

test:
	go test ./...

cover:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func coverage.out
	go tool cover -html=coverage.out -o coverage.html