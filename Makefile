BINARY_NAME=unfurl
BINARY_PATH=./bin/$(BINARY_NAME)

build:
	@go build -o $(BINARY_PATH) cmd/main.go 

run: build
	@$(BINARY_PATH)

clean:
	@rm -rf ./bin 
