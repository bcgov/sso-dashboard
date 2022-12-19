#!/bin/bash

if [[ -n "$1" ]]; then
    db="$1"
else
    db="aggregation"
fi

current_user=$(whoami)

create_superuser() {
    echo "SELECT 'CREATE USER $1 SUPERUSER' WHERE NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = '$1')\gexec" | psql -d postgres
}

create_db() {
    echo "SELECT 'CREATE DATABASE $1' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$1')\gexec" | psql -d postgres
}

create_superuser $current_user
create_db $current_user
create_db $db
