# Build the manager binary
FROM golang:1.24.2 AS builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY .netrc /root/.netrc

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN GOPRIVATE=github.com/cloudogu/doguctl,github.com/cloudogu/cesapp go mod download

# Copy the go source
COPY cmd cmd
COPY internal internal

# Copy .git files as the build process builds the current commit id into the binary via ldflags
COPY .git .git

# Copy build files
COPY build build
COPY Makefile Makefile

# Build
RUN go mod vendor
RUN make compile-generic

FROM gcr.io/distroless/static:nonroot
LABEL maintainer="hello@cloudogu.com" \
      NAME="dogu-additional-mounts-init" \
      VERSION="0.1.2"

WORKDIR /

COPY --from=builder /workspace/target/dogu-additional-mounts-init /dogu-additional-mounts-init

ENTRYPOINT ["/dogu-additional-mounts-init"]
