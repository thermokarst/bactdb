#!/bin/sh

if [ -f .env ]; then
    source .env
fi

SECRET=secret \
PGDATABASE=bactdbtest \
DOMAINS="http://localhost:4200" \
ACCOUNT_KEYS=$ACCOUNT_KEYS \
bash -c 'go run *.go serve'
