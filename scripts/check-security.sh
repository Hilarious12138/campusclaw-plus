#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
DOCKER=${DOCKER_BIN:-docker}
python3 scripts/security_scan.py --history
python3 -m unittest discover -s scripts -p test_security_scan.py
sh scripts/audit.sh
"$DOCKER" run --rm -v "$PWD/backend:/work" -w /work -v campusclaw-plus-go-cache:/go/pkg golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb go test ./...
