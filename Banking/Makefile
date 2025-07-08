build:
	@go build -o bin/gobank
run:
	@./bin/gobank
test:
	#@go test-v/...
	@go test -v ./...
	@go test -v ./... -cover
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out
