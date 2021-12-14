# Start from the latest golang base image
FROM golang:latest AS compiler

# Set the working directory inside the container
WORKDIR /cchc/predictor/aggregator

# Copy dependencies prior to building so that this layer is cached unless
# specified dependencies change
COPY go.mod go.sum /cchc/
RUN go mod download

# Copy Go code into the app
COPY common /cchc/common
COPY predictor/aggregator /cchc/predictor/aggregator

# Build the Go app, making sure it is a static binary with no debugging symbols
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags="-w -s" -o aggregator

# Start from the same R version used for APB on Argo HPC
FROM rocker/tidyverse:3.5.2

# Set the working directory inside the container
WORKDIR /predictor

# Install build dependencies for R packages
RUN apt-get update && apt-get install zlib1g-dev

# Install R packages
RUN install2.r --ncpus=-1 --error --skipinstalled Matrix broom dplyr fs futile.logger optparse parsnip readr recipes sessioninfo data.table text2vec tokenizers stringr

# Copy R scripts
COPY predictor/bin /predictor
COPY predictor/test /predictor/test

# Copy over just the static binary to root
COPY --from=compiler /cchc/predictor/aggregator/aggregator /aggregator

# Create non-root user information
# RUN echo "cchc:x:65534:65534:CCHC:/:" >> /etc/passwd

# Run as non-root user in container
# USER cchc

# Command to run the executable
CMD ["/aggregator"]
