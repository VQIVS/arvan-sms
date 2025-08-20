.PHONY: build run run-dev run-api run-consumer test clean swagger docker-build docker-deps-up docker-deps-down docker-run-api docker-run-consumer docker-run docker-stop docker-logs-api docker-logs-consumer docker-logs-postgres docker-logs-rabbitmq docker-clean lint lint-fix lint-detailed check install-tools security

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

install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

lint:
	golangci-lint run

lint-detailed:
	golangci-lint run --enable-all --disable=gochecknoglobals,gochecknoinits,godot,gomnd,gomodguard,goerr113,wrapcheck,exhaustruct,ireturn,varnamelen,nosnakecase

lint-fix:
	golangci-lint run --fix
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -local sms-dispatcher -w .; \
	else \
		echo "goimports not found. Run 'make install-tools' to install it."; \
	fi

security:
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Run 'make install-tools' to install it."; \
	fi

check: fmt vet lint test
	@echo "All checks passed!"

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

# Docker commands
docker-build:
	docker build -t sms-service -f build/Dockerfile .

docker-deps-up:
	@echo "Starting PostgreSQL..."
	docker run -d --name sms-postgres \
		-e POSTGRES_DB=sms \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=changeme \
		-p 5432:5432 \
		-v sms_postgres_data:/var/lib/postgresql/data \
		postgres:15-alpine || echo "PostgreSQL container already exists"
	@echo "Starting RabbitMQ..."
	docker run -d --name sms-rabbitmq \
		-e RABBITMQ_DEFAULT_USER=guest \
		-e RABBITMQ_DEFAULT_PASS=guest \
		-p 5672:5672 \
		-p 15672:15672 \
		-v sms_rabbitmq_data:/var/lib/rabbitmq \
		rabbitmq:3-management-alpine || echo "RabbitMQ container already exists"
	@echo "Waiting for services to be ready..."
	@sleep 10

docker-deps-down:
	@echo "Stopping and removing dependencies..."
	docker stop sms-postgres sms-rabbitmq || true
	docker rm sms-postgres sms-rabbitmq || true

docker-run-api:
	docker run -d --name sms-api \
		--link sms-postgres:postgres \
		--link sms-rabbitmq:rabbitmq \
		-p 8080:8080 \
		-v $(PWD)/config.json:/app/config.json:ro \
		sms-service

docker-run-consumer:
	docker run -d --name sms-consumer \
		--link sms-postgres:postgres \
		--link sms-rabbitmq:rabbitmq \
		-v $(PWD)/config.json:/app/config.json:ro \
		sms-service ./sms-consumer

docker-run: docker-deps-up docker-build
	@echo "Starting SMS API..."
	$(MAKE) docker-run-api || echo "API container already exists"
	@echo "Starting SMS Consumer..."
	$(MAKE) docker-run-consumer || echo "Consumer container already exists"
	@echo "All services started!"
	@echo "API: http://localhost:8080"
	@echo "RabbitMQ Management: http://localhost:15672 (guest/guest)"

docker-stop:
	@echo "Stopping application containers..."
	docker stop sms-api sms-consumer || true
	docker rm sms-api sms-consumer || true

docker-logs-api:
	docker logs -f sms-api

docker-logs-consumer:
	docker logs -f sms-consumer

docker-logs-postgres:
	docker logs -f sms-postgres

docker-logs-rabbitmq:
	docker logs -f sms-rabbitmq

docker-clean: docker-stop docker-deps-down
	@echo "Cleaning up Docker resources..."
	docker system prune -f
	@echo "Cleanup completed!"