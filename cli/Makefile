.PHONY: all build install clean test fmt vet run

BINARY_NAME=gel
CMD_PATH=./cmd/gel

all: build

build:
	go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "Build successful! Binary created: $(BINARY_NAME)"

install:
	go install $(CMD_PATH)
	@echo "Installed $(BINARY_NAME) to Go bin path"

run:
	go run $(CMD_PATH)

clean:
	go clean
	rm -f $(BINARY_NAME)
	@echo "Cleaned up binary"

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
