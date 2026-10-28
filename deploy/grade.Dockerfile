FROM golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb AS toolchain
RUN apk add --no-cache build-base
FROM node:24-alpine@sha256:ebfe2f90462722a7a4de65e91990e97fe0d401c70e0e762c5b53302f905ec1c1 AS node
FROM toolchain
COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
RUN ln -s ../lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && ln -s ../lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx
RUN timeout 600 apk add --no-cache chromium font-noto-cjk git python3
ENV PATH="/usr/local/go/bin:/go/bin:$PATH" GOPATH=/go OPENSPEC_TELEMETRY=0 GIT_OPTIONAL_LOCKS=0
RUN npm install --global @fission-ai/openspec@1.13.0 && go install golang.org/x/vuln/cmd/govulncheck@v1.9.0
WORKDIR /work
COPY backend/go.mod backend/go.sum /work/backend/
RUN cd backend && go mod download
COPY frontend/package.json frontend/package-lock.json /work/frontend/
RUN cd frontend && npm ci
COPY . /work
RUN git config --global --add safe.directory /work
ENTRYPOINT ["node", "/work/scripts/grader.mjs"]
