# Use an official Golang runtime as a parent image
FROM golang:1.23-alpine AS builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy the local package files to the container's workspace.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

RUN apk add -U --no-cache ca-certificates

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server cmd/main.go

# Use a Docker multi-stage build to create a lean production image.
# Start from scratch container for the output stage
FROM alpine:latest


WORKDIR /root/

# Copy the binary from the builder stage to the production image
COPY --from=builder /app/server .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Run the server executable
CMD ["./server"]