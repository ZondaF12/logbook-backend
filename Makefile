build:
	@go build -o bin/logbook-backend cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/logbook-backend
