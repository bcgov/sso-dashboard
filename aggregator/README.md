# SSO Keycloak event log aggregator

A Go HTTP server to host the Keycloak event log aggregation system that receives the log metadata from `Grafana Promtail` and upserts the client-level event statistical data.
In order to avoid the custom codebase parsing the requests, it relies on `Grafana Loki`'s public function `ParseRequest` and takes the advantage of code simplicity and reliability.

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
