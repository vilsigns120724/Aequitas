FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o aequitasd ./cmd/aequitasd/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/aequitasd .
COPY --from=builder /app/genesis.json .
EXPOSE 4001
CMD ["./aequitasd"]
