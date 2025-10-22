# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o engine ./cmd/engine

# Runtime stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/engine .

EXPOSE 8080 50051
CMD ["./engine"]