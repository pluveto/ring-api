# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ring-api ./cmd/ring-api

# Runtime stage
FROM alpine:3.18
RUN apk add --no-cache alsa-utils
COPY --from=builder /app/ring-api /usr/local/bin/
COPY audio/ring.wav /var/lib/ring-api/ring.wav

EXPOSE 8080
CMD ["ring-api"]
