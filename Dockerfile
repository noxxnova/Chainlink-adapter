FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o adapter .

# Use a small alpine image for the final container
FROM alpine:3.16

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/adapter .

# Expose the default port
EXPOSE 8080

# Run the adapter
CMD ["./adapter"]