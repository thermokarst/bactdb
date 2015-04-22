#!/bin/sh

SECRET=secret \
PGDATABASE=bactdbtest \
DOMAINS="http://localhost:4200" \
bash -c 'go run *.go migrate --drop'
