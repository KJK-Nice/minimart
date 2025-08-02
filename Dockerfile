# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

#Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application.
# -o /app/server: specifies the output file name and location.
# -ldflags"-w -s": strips debug information, reducing the binary size.
# CGO_ENABLED=0: creates a statically-linked binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags="-w -s" ./cmd/server

# Stage 2: Create the final, minimal image
FROM scratch

# Set the working directory
WORKDIR /app

# Copy the comiled binary from the builder stage
COPY --from=builder /app/server .

# Copy the configuration file
COPY config.yaml .

# Expose the port the app runs on
EXPOSE 3000

# The command to run when the container starts
CMD ["/app/server"]
