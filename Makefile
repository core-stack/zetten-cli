
BINARY_CLI_NAME=zetten-cli
BINARY_SERVICE_NAME=zetten-service

ARGS=""

.PHONY: all build run install clean help

all: build

build-cli:
	@echo "🔨 Building..."
	go build -o bin/$(BINARY_CLI_NAME) ./cmd/cli

run-cli: build
	@echo "🚀 Running..."
	./bin/$(BINARY_CLI_NAME) $(ARGS)

install-cli:
	@echo "📦 install $(GOBIN)..."
	go install .

clean-cli:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_CLI_NAME)

test:
	go test ./... -v


build-service:
	@echo "🔨 Building..."
	go build -o bin/$(BINARY_SERVICE_NAME) ./cmd/service

run-service: build
	@echo "🚀 Running..."
	./bin/$(BINARY_SERVICE_NAME) $(ARGS)

install-service:
	@echo "📦 install $(GOBIN)..."
	go install .

clean:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_SERVICE_NAME)