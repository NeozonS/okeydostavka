BINARY_NAME=testcenozavr
MAIN_PACKAGE=./cmd/server
GO_FLAGS=-ldflags="-s -w"
GO_ENV=CGO_ENABLED=0
DOCKER_IMAGE=okey-parser
DOCKER_TAG=latest
PORT=8080
OUTPUT_DIR=./output

.PHONY: build run clean docker-build docker-run help

all: build

build:
	$(GO_ENV) go build $(GO_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@mkdir -p $(OUTPUT_DIR)
	docker run --rm -p $(PORT):$(PORT) \
		--env-file .env \
		-v $(shell pwd)/$(OUTPUT_DIR):/data \
		$(DOCKER_IMAGE):$(DOCKER_TAG)