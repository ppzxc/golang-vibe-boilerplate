# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# Run stage
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/server /server
COPY --from=builder /app/internal/adapter/postgresrepo/migrations /migrations

EXPOSE 8080

ENTRYPOINT ["/server"]
