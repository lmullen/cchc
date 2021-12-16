# Start from the latest golang base image
FROM golang:latest AS compiler

# Set the working directory inside the container
WORKDIR /cchc/itemmd

# Copy dependencies prior to building so that this layer is cached unless
# specified dependencies change
COPY go.mod go.sum /cchc/
RUN go mod download

# Copy the source from the current directory to the container
COPY common /cchc/common
COPY itemmd /cchc/itemmd

# Build the Go app, making sure it is a static binary with no debugging symbols
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags="-w -s" -o cchc-itemmd

# Create non-root user information
RUN echo "cchc:x:65534:65534:CCHC:/:" > /etc_passwd

# Start over with a completely empty image
FROM scratch

# Include the certificates since we have to access an HTTPS API
COPY --from=compiler /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy over just the static binary to root
COPY --from=compiler /cchc/itemmd/cchc-itemmd /cchc-itemmd

# Copy over non-root user information
COPY --from=0 /etc_passwd /etc/passwd

# Run as non-root user in container
USER cchc

# Command to run the executable
CMD ["/cchc-itemmd"]
