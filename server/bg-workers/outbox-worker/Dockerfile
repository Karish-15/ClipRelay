FROM golang AS builder

WORKDIR /app

# Copy only go.mod first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o outbox-worker main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/outbox-worker .
COPY .env .

CMD ["./outbox-worker"]