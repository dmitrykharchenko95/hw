BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o ./bin/calendar -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o ./bin/scheduler ./cmd/scheduler
	go build -v -o ./bin/sender ./cmd/sender

run: build
	./bin/calendar -config ./configs/calendar_config.json
	./bin/scheduler -config ./configs/scheduler_config.json
	./bin/sender -config ./configs/sender_config.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	./bin/calendar version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	mkdir -p internal/server/pb

	protoc \
		--proto_path=api/ \
        --go_out=. \
        --go-grpc_out=. \
        api/*.proto

.PHONY: build run build-img run-img version test lint
