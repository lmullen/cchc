# Start from the latest golang base image
FROM golang:latest AS compiler

# Set the working directory inside the container
WORKDIR /app

# Copy dependencies prior to building so that this layer is cached unless
# specified dependencies change
COPY go.mod go.sum /app/
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY itemmd /app

# Build the Go app, making sure it is a static binary with no debugging symbols
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags="-w -s" -o app

# Create non-root user information
RUN echo "cchc:x:65534:65534:CCHC:/:" > /etc_passwd

# Start over with a completely empty image
FROM scratch

# Include the certificates since we have to access an HTTPS API
COPY --from=compiler /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy over just the static binary to root
COPY --from=compiler /app/app /app

# Copy over non-root user information
COPY --from=0 /etc_passwd /etc/passwd

# Run as non-root user in container
USER cchc

# Command to run the executable
CMD ["/app"]
