# Build stage
FROM golang:alpine AS builder

# Install git
RUN apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /go/src/app

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Fetch dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o /go/bin/app -v .

# Final stage
FROM alpine:latest

# Install CA certificates, required for making HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/bin/app /app

# Copy the SQLite database file
COPY mainDB.db .

# Command to run the executable
ENTRYPOINT ["/app"]

LABEL Name=mattodo Version=0.0.1
EXPOSE 9003
