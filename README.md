# Ring API

A production-ready audio playback service with concurrency control and cross-platform support.

## System Requirements

- Go 1.21+
- Linux: ALSA tools (`aplay`)
    - Arch Linux: `yay -S alsa-utils`
- Windows: FFmpeg in PATH (`ffplay.exe`)

## Quick Start

```bash
# Using Go
AUDIO_PATH=/path/to/audio.wav go run cmd/ring-api/main.go

# Using Docker
docker run -p 8080:8080 -e AUDIO_PATH=/audio.wav pluveto/ring-api
```

## API Endpoint

**GET /api/ring**
- Success: 202 Accepted
- Conflicts: 429 Too Many Requests
- Errors: 4xx/5xx with error message

Example request:

```bash
curl -X GET http://localhost:8080/api/ring
```

## Configuration

| Environment Variable | Default                  | Description         |
|-----------------------|--------------------------|---------------------|
| AUDIO_PATH            | OS-specific path         | Audio file location |
| HTTP_PORT             | 8080                     | Service port        |

## Testing

```bash
# Run all tests
go test -v ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Platform-specific testing
GOOS=linux go test -v ./internal/player
GOOS=windows go test -v ./internal/player
```

## Build

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o ring-api-linux cmd/ring-api/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o ring-api.exe cmd/ring-api/main.go
```

## Deployment

```dockerfile
# Build image
docker build -t ring-api .

# Run container
docker run -d -p 8080:8080 \
  -e AUDIO_PATH=/audio/ring.wav \
  -v /host/audio:/audio \
  ring-api
```

## License

MIT License. See [LICENSE](LICENSE) for details.
```

