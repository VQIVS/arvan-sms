FROM golang:1.24-slim as builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./sms-api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./sms-consumer cmd/consumer/main.go

FROM alpine as deploy

RUN apk add --no-cache tzdata ca-certificates
ENV TZ=Asia/Tehran

WORKDIR /app

COPY --from=builder /app/sms-api ./sms-api
COPY --from=builder /app/sms-consumer ./sms-consumer
COPY config.json ./config.json

EXPOSE 8080

CMD ["./sms-api"]
