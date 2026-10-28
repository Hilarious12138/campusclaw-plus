#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
DOCKER=${DOCKER_BIN:-docker}
"$DOCKER" run --rm -v "$PWD/frontend:/work" -v campusclaw-plus-audit-node:/work/node_modules -w /work node:24-alpine@sha256:ebfe2f90462722a7a4de65e91990e97fe0d401c70e0e762c5b53302f905ec1c1 sh -c 'npm ci && npm audit --audit-level=high'
"$DOCKER" run --rm -v "$PWD/backend:/work" -w /work -v campusclaw-plus-go-cache:/go/pkg golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb sh -c 'go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...'
