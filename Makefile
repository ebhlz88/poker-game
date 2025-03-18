build:
	@go build -o bin/poker-game
run: build
	@./bin/poker-game
test:
	go test -v ./...