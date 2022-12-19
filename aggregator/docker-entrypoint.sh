#!/bin/bash -e

if ! command -v nc &>/dev/null; then
    echo "nc not found"
    exit
fi

while ! nc -z ${DB_HOSTNAME} ${DB_PORT:-5432}; do
    echo "wait for ${DB_HOSTNAME} ${DB_PORT:-5432}"
    sleep 0.5
done

alembic upgrade head

echo "DB migration finished."

./app/main
