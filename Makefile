.PHONY: build run run-dev run-api run-consumer test clean swagger docker-build

build:
	go build -o ./bin/api ./cmd/api
	go build -o ./bin/consumer ./cmd/consumer

test:
	go test -v ./...

clean:
	rm -rf ./bin

swagger:
	swag init -g ./cmd/api/main.go -o ./docs

fmt:
	go fmt ./...

vet:
	go vet ./...

deps:
	go mod download

update-deps:
	go get -u ./...
	go mod tidy

run-api:
	go run ./cmd/api/main.go

run-consumer:
	go run ./cmd/consumer/main.go

run-dev:
	$(MAKE) build && $(MAKE) swagger && ($(MAKE) run-api & $(MAKE) run-consumer)