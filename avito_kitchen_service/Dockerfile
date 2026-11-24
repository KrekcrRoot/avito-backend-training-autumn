FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kitchen ./cmd/app

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/kitchen .

EXPOSE 3001

CMD ["./kitchen"]