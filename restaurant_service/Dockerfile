FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o restaurant ./cmd/app

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/restaurant .

EXPOSE 3000

CMD ["./restaurant"]