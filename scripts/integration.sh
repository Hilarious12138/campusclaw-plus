#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
DOCKER=${DOCKER_BIN:-docker}
config_file="$PWD/.env.integration-$$"
trap 'rm -f "$config_file"' EXIT INT TERM
umask 077
"$DOCKER" run --rm -v "$PWD/scripts/test-env.mjs:/test-env.mjs:ro" node:24-alpine@sha256:ebfe2f90462722a7a4de65e91990e97fe0d401c70e0e762c5b53302f905ec1c1 node /test-env.mjs >"$config_file"
GRADE_ENV_FILE="$config_file" sh scripts/grade.sh
