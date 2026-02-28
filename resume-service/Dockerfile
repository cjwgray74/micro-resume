FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy only go.mod (go.sum may not exist)
COPY go.mod ./
RUN go mod download || true

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o resume-service main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/resume-service .

ENV PORT=8081
EXPOSE 8081

CMD ["./resume-service"]