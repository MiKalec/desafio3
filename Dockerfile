# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ordersystem ./cmd/ordersystem

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ordersystem .

# Copy config file
COPY cmd/ordersystem/.env .

EXPOSE 8000 50051 8080

CMD ["./ordersystem"]
