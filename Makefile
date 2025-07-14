
BINARY_NAME=zetten

ARGS=""

.PHONY: all build run install clean help

all: build

build:
	@echo "🔨 Building..."
	go build -o bin/$(BINARY_NAME) .

run: build
	@echo "🚀 Running..."
	./bin/$(BINARY_NAME) $(ARGS)

install:
	@echo "📦 install $(GOBIN)..."
	go install .

clean:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_NAME)

test:
	go test ./... -v