# Define directories
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
COVERAGE_DIR := $(BUILD_DIR)/coverage

all: dirs test gosvcsample

.PHONY: dirs
dirs:
	mkdir -p $(BIN_DIR)
	mkdir -p $(COVERAGE_DIR)

test: dirs
	go test -v ./... -coverprofile=$(COVERAGE_DIR)/coverage.out
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

gosvcsample: dirs
	go build -o $(BIN_DIR)/gosvcsample ./cmd/gosvcsample/main.go

clean:
	rm -rf $(BUILD_DIR)
