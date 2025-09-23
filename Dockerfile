FROM golang:1.24-alpine AS builder

WORKDIR /app

# https://github.com/confluentinc/confluent-kafka-go/blob/master/README.md#librdkafka
RUN apk add --no-cache git ca-certificates build-base

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# https://github.com/confluentinc/confluent-kafka-go/blob/master/README.md#static-builds-on-linux
RUN go build -tags musl -o main ./cmd/notification/main.go

FROM alpine:3.18

RUN apk add --no-cache ca-certificates bash netcat-openbsd

WORKDIR /app

COPY --from=builder /app/main .

# COPY bin/wait-for-it.sh ./bin/wait-for-it.sh
# RUN chmod +x ./bin/wait-for-it.sh

COPY .env .
RUN mkdir logs/

EXPOSE 8080

CMD ["./main"]
# CMD ["./bin/wait-for-it.sh", "localhost", "9092", "./main"]