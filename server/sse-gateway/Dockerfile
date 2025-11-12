FROM golang AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o sse-gateway main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/sse-gateway .

EXPOSE 9002 9003 9004

CMD ["./sse-gateway"]
