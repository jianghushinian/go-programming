# Build the cronx binary
FROM golang:1.27 AS builder

ARG TARGETOS
ARG TARGETARCH

ARG VERSION=unknown
ARG COMMIT=""
ARG TREESTATE=""
ARG BUILDTIME=""

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build \
      -ldflags "-X cronx/pkg/version.Version=${VERSION} \
                -X cronx/pkg/version.GitCommit=${COMMIT} \
                -X cronx/pkg/version.GitTreeState=${TREESTATE} \
                -X cronx/pkg/version.BuildTime=${BUILDTIME}" \
      -o ./cronx ./cmd/main.go

# Use distroless as minimal base image to package the cronx binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/cronx .
USER 65532:65532

ENTRYPOINT ["/cronx"]
