# Use the official Loki base image
FROM docker.io/grafana/loki:2.9.2

# Set the working directory to /etc/loki
WORKDIR /etc/loki

# Copy your custom config file into the container
COPY ./config.yaml /etc/loki/config/
COPY ./runtime-config.yaml /etc/loki/runtime-config/runtime-config.yaml

# Give loki user permissions for config folder
USER root
RUN mkdir -p /var/loki/retention

# Ensure the container starts with Loki using your custom config file
ENTRYPOINT ["/usr/bin/loki"]
