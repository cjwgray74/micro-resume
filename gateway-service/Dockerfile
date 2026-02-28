FROM golang:1.25 AS builder

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download || true

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/gateway-service .

ENV PORT=8080
EXPOSE 8080

CMD ["./gateway-service"]