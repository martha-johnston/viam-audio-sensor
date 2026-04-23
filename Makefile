export CGO_ENABLED := 1

BIN_OUTPUT_PATH = bin
UNAME_S ?= $(shell uname -s)

build:
	@rm -f $(BIN_OUTPUT_PATH)/audio-sensor
	@go build -o $(BIN_OUTPUT_PATH)/audio-sensor main.go

module.tar.gz: build
	@rm -f $(BIN_OUTPUT_PATH)/module.tar.gz
	@tar czf $(BIN_OUTPUT_PATH)/module.tar.gz $(BIN_OUTPUT_PATH)/audio-sensor

setup:
	@go mod tidy

clean:
	@rm -rf $(BIN_OUTPUT_PATH)

format:
	@gofmt -w -s .
