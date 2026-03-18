FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /balance ./cmd/balance

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /balance /balance
ENTRYPOINT ["/balance"]