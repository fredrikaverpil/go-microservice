SHELL := /bin/bash

.PHONY: proto
proto: proto-tools buf-dep-update proto-lint buf-generate

# --- proto ---

# TODO: pin versions?
# BUF_VERSION=1.28.1
# API_LINTER_VERSION=0.1.0

# Install proto tools
.PHONY: proto-tools
proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/googleapis/api-linter/cmd/api-linter@latest
	go install go.einride.tech/aip/cmd/protoc-gen-go-aip@latest

  # these are defined in buf.gen.yaml and their versions in buf.lock: 
  # go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Install go tools
.PHONY: go-tools
go-tools:
	# AIP convenience functions, such as pagination, resourcename etc.
	# go get -u go.einride.tech/aip

	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install openapi tools
.PHONY: openapi-tools
openapi-tools:
	go install github.com/daveshanley/vacuum@latest

.PHONY: buf-dep-update
buf-dep-update:
	cd proto && buf dep update

.PHONY: proto-lint
proto-lint: buf-lint api-lint

.PHONY: openapi-lint
openapi-lint: openapiv2-lint openapiv3-lint

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

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run --config .golangci.yaml

.PHONY: openapiv2-lint
openapiv2-lint:
	go run github.com/daveshanley/vacuum@latest lint -d proto/gen/openapiv2/**/*.json

.PHONY: openapiv3-lint
openapiv3-lint:
	go run github.com/daveshanley/vacuum@latest lint -d proto/gen/openapiv3/*.yaml

.PHONY: buf-generate
buf-generate:
	cd proto && buf lint && buf generate

.PHONY: run-server
run-server:
	GO_ENV=development go run cmd/server/main.go

.PHONY: run-server-prod
run-server-prod:
	GO_ENV=production go run cmd/server/main.go

