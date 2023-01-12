FROM golang:1.19.4-bullseye as builder

WORKDIR /app

COPY . .
RUN make build

FROM python:3.11.1-slim-bullseye

RUN apt-get update -q && \
  apt-get install -y --no-install-recommends apt-utils && \
  apt-get install -yqq \
  netcat \
  build-essential && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY ["docker-entrypoint.sh", "requirements.txt", "alembic.ini", "config.py", "models.py", "database.py", "./"]
COPY alembic ./alembic

RUN pip install -r requirements.txt

COPY --from=builder /app/build/aggregator /app/aggregator
COPY --from=builder /app/build/compactor /app/compactor

ENV TZ America/Vancouver
ENV PORT 8080
EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
