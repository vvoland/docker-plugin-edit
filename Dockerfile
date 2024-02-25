FROM --platform=$BUILDPLATFORM golang:1.21-bullseye AS base

FROM base as build

WORKDIR /go/src

COPY go.mod go.sum ./
RUN --mount=type=bind,src=./go.mod,target=./go.mod \
    --mount=type=bind,src=./go.sum,target=./go.sum \
    go mod download

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=bind,target=/go/src \
    --mount=type=cache,target=/root/.cache/go-build \
        CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /binary ./cmd/edit


FROM base AS gopls
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
        GOBIN=/build/ GO111MODULE=on go install "golang.org/x/tools/gopls@latest" \
     && /build/gopls version


FROM base as devcontainer
COPY --from=gopls /build/gopls /usr/local/bin/gopls


FROM scratch AS binary
COPY --from=build /binary /docker-edit
