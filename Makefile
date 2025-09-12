BINARY=wyag

BIN_DIR=bin

GO=go

all: build

build:
	$(GO) build -o $(BIN_DIR)/$(BINARY) .

install:
	$(GO) install ./cmd/wyag

run:
	$(GO) run ./cmd/wyag

clean:
	rm -rf $(BIN_DIR)

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(BINARY)-linux ./cmd/wyag

build-windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(BINARY).exe ./cmd/wyag
