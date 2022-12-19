FROM golang:1.19.4-bullseye as builder

WORKDIR /app

COPY . .
RUN make build

FROM debian:bullseye-slim

RUN apt-get update -q && \
  apt-get install -y --no-install-recommends apt-utils && \
  apt-get install -yqq \
  netcat \
  build-essential && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=builder /app/build/main /app/main
COPY docker-entrypoint.sh docker-entrypoint.sh

ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
