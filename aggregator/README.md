# SSO Keycloak event log aggregator & compactor

## Aggregator

A Go HTTP server to host the Keycloak event log aggregation system that receives the log metadata from `Grafana Promtail` and upserts the client-level event statistical data.
In order to avoid the custom codebase parsing the requests, it relies on `Grafana Loki`'s public function `ParseRequest` and takes the advantage of code simplicity and reliability.

## Compactor

A lightweight Go server running scheduled jobs. There are two cronjobs it controls:

1. Deleting old client events.
2. Collecting client session counts.

## Environment Variables

- `DB_HOSTNAME`: the name of host to connect to.
- `DB_PORT`: : the port number of host to connect to.
- `DB_DATABASE`: the database name.
- `DB_USERNAME`: the user name to connect as.
- `DB_PASSWORD`: the password to be used for password authentication.
- `RETENTION_PERIOD`: the duration of time to keep the aggregated data.
  - please see [Postgres Interval Input](https://www.postgresql.org/docs/current/datatype-datetime.html#DATATYPE-INTERVAL-INPUT) for the unit convention.
- `RC_WEBHOOK`: The url for the rocketchat webhook to use when notifying from the compactor
- `DEV_KEYCLOAK_URL`: The development keycloak base URL
- `DEV_KEYCLOAK_CLIENT_ID`: The development keycloak client id
- `DEV_KEYCLOAK_USERNAME`: The development keycloak username
- `DEV_KEYCLOAK_PASSWORD`: The development keycloak passowrd
- `TEST_KEYCLOAK_URL`: The test keycloak base URL
- `TEST_KEYCLOAK_CLIENT_ID`: The test keycloak client id
- `TEST_KEYCLOAK_USERNAME`: The test keycloak username
- `TEST_KEYCLOAK_PASSWORD`: The test keycloak passowrd
- `PROD_KEYCLOAK_URL`: The prod keycloak base URL
- `PROD_KEYCLOAK_CLIENT_ID`: The prod keycloak client id
- `PROD_KEYCLOAK_USERNAME`: The prod keycloak username
- `PROD_KEYCLOAK_PASSWORD`: The prod keycloak passowrd

## Local development setup

1. start the local `Postgres` database.

   ```sh
   make db
   ```

1. if working on the `aggregator` codebase, run:

   ```sh
   make dev
   ```

   - the entrypoint for `aggregator` is `main.go`.

1. if working on the `compactor` codebase, run:

   ```sh
   make compactor
   ```

   - the entrypoint for `compactor` is `compactor.go`.

## Database migration

It uses `Alembic`, which is a lightweight database migration tool, to migrate database schema. The step to create a migration script is:

1. install the required tools.

   ```sh
   make install
   ```

1. run the local database with the up-to-date schema.

   ```sh
   make db
   ```

1. make changes in `models.py`.
1. create an alembic revision.

   ```sh
   alembic revision --autogenerate -m "<message>"
   ```

1. apply the changes in the local database.

   ```sh
   alembic upgrade head
   ```

## Deployment

The docker images of the SSO Keycloak aggregator is built and published via GitHub Actions and located in GitHub Packages.

- see [SSO Aggregator package](https://github.com/bcgov/sso-dashboard/pkgs/container/sso-aggregator) to check the status of the published docker images.
- see [publish-aggregator.yml](../.github/workflows/publish-aggregator.yml) to check the CD pipeline of the Go server.
