# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o workspace-controller main.go

# Final stage
FROM scratch

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /app/workspace-controller .

# Copy the system definition
COPY system/system-definition.yaml ./system/system-definition.yaml

# Set the entrypoint
ENTRYPOINT ["./workspace-controller"]

# Default command
CMD ["help"]
