# FROM golang:1.25-alpine AS builder

# WORKDIR /app

# # Copy only go.mod (go.sum may not exist)
# COPY go.mod ./
# RUN go mod download || true

# # Copy the rest of the source code
# COPY . .

# # Build the binary
# RUN CGO_ENABLED=0 GOOS=linux go build -o resume-service main.go

# FROM alpine:3.18

# WORKDIR /app

# COPY --from=builder /app/resume-service .

# ENV PORT=8081
# EXPOSE 8081

# CMD ["./resume-service"]

# -----------------------------
# Build stage (Go 1.22)
# -----------------------------
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the pdf-service binary
RUN CGO_ENABLED=0 GOOS=linux go build -o pdf-service main.go


# -----------------------------
# Runtime stage
# -----------------------------
FROM debian:bullseye

WORKDIR /app

# Install Chromium for chromedp
RUN apt-get update && apt-get install -y chromium && apt-get clean

# Copy binary and templates
COPY --from=builder /app/pdf-service .
COPY --from=builder /app/templates ./templates

ENV PORT=8082
EXPOSE 8082

CMD ["./pdf-service"]