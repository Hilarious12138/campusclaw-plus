#!/bin/sh
set -eu
if [ "${#MYSQL_ROOT_PASSWORD}" -lt 24 ] || [ "$MYSQL_ROOT_PASSWORD" = "$MYSQL_PASSWORD" ]; then
    printf '%s\n' '{"level":"ERROR","module":"config","event":"invalid_config","request_id":"system","config_key":"DB_ROOT_PASSWORD"}'
    exit 64
fi
if [ "${#MYSQL_PASSWORD}" -lt 24 ]; then
    printf '%s\n' '{"level":"ERROR","module":"config","event":"invalid_config","request_id":"system","config_key":"DB_PASSWORD"}'
    exit 64
fi
exec docker-entrypoint.sh "$@"
