#!/bin/bash

gosu postgres postgres --single <<- EOSQL
CREATE USER "$PGUSER" WITH SUPERUSER PASSWORD '$PGPASSWORD';
CREATE DATABASE $PGUSER;
CREATE DATABASE $PGDB;
EOSQL

{ echo; echo "host all \"$PGUSER\" 0.0.0.0/0 md5"; } >> "$PGDATA"/pg_hba.conf

