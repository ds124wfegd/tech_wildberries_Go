# Using the official Go image
FROM golang:1.21-alpine AS builder

# Installing the working directory
WORKDIR /app

# Copying the dependency files
COPY go.mod go.sum ./

# Loading dependencies
RUN go mod download

# Copying the source code
COPY . .

# Building the app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# We use a minimal image for the final container
FROM alpine:latest

# Installing ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copying the compiled application
COPY --from=builder /app/main .

# Launching the app
CMD ["./main"] 