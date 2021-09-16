# Start from the same R version used for APB on Argo HPC
FROM rocker/tidyverse:3.5.2

# Set the working directory inside the container
WORKDIR /predictor

# Create non-root user information
RUN echo "cchc:x:65534:65534:CCHC:/:" > /etc_passwd

# Install build dependencies for R packages
RUN apt-get update && apt-get install zlib1g-dev

# Install R packages
RUN install2.r --ncpus=-1 --error --skipinstalled Matrix broom dplyr fs futile.logger optparse parsnip readr recipes sessioninfo data.table text2vec tokenizers stringr

# Copy R scripts
COPY predictor/bin /predictor

# Create non-root user information
RUN echo "cchc:x:65534:65534:CCHC:/:" > /etc/passwd

# Run as non-root user in container
USER cchc

# Command to run the executable
CMD ["bash"]
