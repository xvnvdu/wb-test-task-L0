FROM golang:1.23.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o /orders-service ./cmd/server

FROM alpine:latest
WORKDIR /app 
COPY --from=builder /orders-service .
COPY web ./web

EXPOSE 8080
CMD ["./orders-service"]
