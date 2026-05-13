# Copyright 2026 itscheems
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# syntax=docker/dockerfile:1.9
#
# Build context must be the repo root, not deploy/, because we need the
# whole workspace (go.work + every module). From repo root:
#   docker build -f deploy/Dockerfile --build-arg SERVICE=app -t golang-template:dev .

# ────────────────────────────────
# Builder stage
# ────────────────────────────────
FROM golang:1.26.3-alpine AS builder
WORKDIR /src

# Pick which service under services/<SERVICE> to build.
ARG SERVICE=app

# Copy workspace + module manifests first to maximise layer cache hits.
COPY go.work go.work.su[m] ./
COPY services/${SERVICE}/go.mod services/${SERVICE}/go.su[m] services/${SERVICE}/
COPY libs/ libs/

# Fetch deps for the targeted module + its workspace neighbours.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/${SERVICE} && go mod download

# Now bring in the rest of the source for compilation.
COPY services/${SERVICE} services/${SERVICE}

ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/${SERVICE} && \
    CGO_ENABLED=0 GOOS=linux go build \
        -trimpath \
        -ldflags "-s -w \
            -X github.com/itscheems/golang-template/services/${SERVICE}/internal/version.Version=${VERSION} \
            -X github.com/itscheems/golang-template/services/${SERVICE}/internal/version.Commit=${COMMIT} \
            -X github.com/itscheems/golang-template/services/${SERVICE}/internal/version.Date=${BUILD_DATE}" \
        -o /out/app ./cmd/${SERVICE}

# ────────────────────────────────
# Runtime stage — distroless static, non-root
# ────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

COPY --from=builder /out/app /app/app

USER nonroot:nonroot
WORKDIR /app

# EXPOSE 8080
# HEALTHCHECK configured at orchestrator level (k8s livenessProbe, etc.)

ENTRYPOINT ["/app/app"]
