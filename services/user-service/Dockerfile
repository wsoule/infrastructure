# start from a minimal Go image
FROM golang:1.25-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go module files first (for caching dependencies)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# build the Go binary
RUN go build -o user-service


# -- Runtime stage --
FROM alpine:3.20

# Create working dir
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/user-service .

# Expose the port
EXPOSE 8081

# Run the binary
CMD ["./user-service"]
