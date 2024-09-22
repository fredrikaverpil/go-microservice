SHELL := /bin/bash

.PHONY: proto
proto: proto-tools buf-dep-update proto-lint buf-generate

# --- proto ---

# TODO: pin versions?
# Tool versions
# BUF_VERSION=1.28.1
# PROTOC_GEN_GO_VERSION=1.31.0
# PROTOC_GEN_GO_GRPC_VERSION=1.3.0
# PROTOC_GEN_GRPC_GATEWAY_VERSION=2.18.0
# PROTOC_GEN_OPENAPIV2_VERSION=2.18.0
# API_LINTER_VERSION=0.1.0

# Install tools
.PHONY: proto-tools
proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/googleapis/api-linter/cmd/api-linter@latest

.PHONY: buf-dep-update
buf-dep-update:
	cd proto && buf dep update

.PHONY: proto-lint
proto-lint: buf-lint api-lint

.PHONY: buf-lint
buf-lint:
	cd proto && buf lint

# Lint proto files using google's AIP linter
.PHONY: api-lint
api-lint:
	cd proto && \
  buf build -o descriptor-set.pb && \
	api-linter --descriptor-set-in=descriptor-set.pb \
		--config=api-linter.yaml \
		--output-format=yaml \
		--set-exit-status \
		gomicroservice/v1/*.proto

.PHONY: buf-generate
buf-generate:
	cd proto && buf lint && buf generate

