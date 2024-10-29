# Use the official Loki base image
FROM docker.io/grafana/loki:2.9.2

# Set the working directory to /etc/loki
WORKDIR /etc/loki

# Copy your custom config file into the container
COPY ./config.yaml /etc/loki/config/
COPY ./runtime-config.yaml /etc/loki/runtime-config/runtime-config.yaml

# Setup loki user permissions for working directory
USER root
RUN mkdir -p /var/loki/retention && chown -R loki:loki /var/loki
USER loki

ENTRYPOINT ["/usr/bin/loki"]
