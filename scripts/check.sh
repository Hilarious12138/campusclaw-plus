#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
DOCKER=${DOCKER_BIN:-docker}
python3 scripts/check.py
"$DOCKER" build -t campusclaw-plus-check -f deploy/check.Dockerfile .
"$DOCKER" run --rm -v "$PWD/backend:/work" -w /work -v campusclaw-plus-go-cache:/go/pkg -v campusclaw-plus-go-build:/root/.cache/go-build campusclaw-plus-check sh -c 'test -z "$(gofmt -l .)" && CGO_ENABLED=1 go test -race ./... && go vet ./...'
"$DOCKER" run --rm -v "$PWD/frontend:/work" -v campusclaw-plus-check-node:/work/node_modules -w /work node:24-alpine@sha256:ebfe2f90462722a7a4de65e91990e97fe0d401c70e0e762c5b53302f905ec1c1 sh -c 'npm ci && npm test && npm run build'
OPENSPEC_TELEMETRY=0 openspec validate --all --strict
