FROM golang AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server .

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/api-server /app/api-server

# Copy env file (optional, for local docker testing)
COPY .env /app/.env

# Expose API port
EXPOSE 8080

# Run API-server
ENTRYPOINT ["/app/api-server"]
