ARTIFACT_ID=dogu-additional-mounts-init
VERSION=0.1.2
MAKEFILES_VERSION=10.6.0
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
GOTAG?=1.26
LINT_VERSION?=v2.9.0
MOCKERY_IGNORED=vendor,build,docs,generatedv

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

include build/make/variables.mk
include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
GO_BUILD_FLAGS+="./cmd/copier"
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/mocks.mk
include build/make/release.mk
