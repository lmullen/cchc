# Start from the latest golang base image
FROM rocker/r-ver:3.6.1

# Set the working directory inside the container
WORKDIR /predictor

# Create non-root user information
RUN echo "cchc:x:65534:65534:CCHC:/:" > /etc_passwd

# Install build dependencies for R packages
RUN apt-get update && apt-get install zlib1g-dev

# Install R packages
RUN install2.r --error --skipinstalled Matrix
RUN install2.r --error --skipinstalled broom
RUN install2.r --error --skipinstalled dplyr
RUN install2.r --error --skipinstalled fs
RUN install2.r --error --skipinstalled futile.logger
RUN install2.r --error --skipinstalled optparse
RUN install2.r --error --skipinstalled parsnip
RUN install2.r --error --skipinstalled readr
RUN install2.r --error --skipinstalled recipes
RUN install2.r --error --skipinstalled sessioninfo
RUN install2.r --error --skipinstalled data.table
RUN install2.r --error --skipinstalled text2vec
RUN install2.r --error --skipinstalled tokenizers

# Copy R scripts
COPY predictor/bin /predictor

# Create non-root user information
RUN echo "cchc:x:65534:65534:CCHC:/:" > /etc/passwd

# Run as non-root user in container
USER cchc

# Command to run the executable
CMD ["bash"]
