#!/bin/bash

set -e
set -u

function create_user_and_database() {
  local database=$1
  echo "  Creating user and database '$database'"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
      CREATE DATABASE "$database" OWNER $POSTGRES_USER;
EOSQL
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
  echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
  for database in $(echo "$POSTGRES_MULTIPLE_DATABASES" | tr ',' ' '); do
    create_user_and_database "$database"
  done
  echo "Multiple databases created"
fi
