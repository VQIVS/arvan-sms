FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./sms-api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./sms-consumer ./cmd/consumer/main.go

FROM alpine:latest as deploy

RUN apk update && apk add --no-cache tzdata ca-certificates && \
    addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

ENV TZ=Asia/Tehran

WORKDIR /app

COPY --from=builder /app/sms-api ./sms-api
COPY --from=builder /app/sms-consumer ./sms-consumer

COPY config.yaml ./config.yaml

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080
CMD ["./sms-api"]