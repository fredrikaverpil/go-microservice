SHELL := /bin/bash

# TODO: pin versions?
# Tool versions
# BUF_VERSION=1.28.1
# PROTOC_GEN_GO_VERSION=1.31.0
# PROTOC_GEN_GO_GRPC_VERSION=1.3.0
# PROTOC_GEN_GRPC_GATEWAY_VERSION=2.18.0
# PROTOC_GEN_OPENAPIV2_VERSION=2.18.0

# Install tools
.PHONY: tools
tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Generate protobuf code
.PHONY: proto
proto:
	cd proto && buf lint && buf generate

.PHONY: lint
lint:
	cd proto && buf lint
