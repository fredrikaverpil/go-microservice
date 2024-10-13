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
	go install github.com/einride/protoc-gen-go-aip-test@latest

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

	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.2

# Install openapi tools
.PHONY: openapi-tools
openapi-tools:
	# go install github.com/daveshanley/vacuum@latest
	go install github.com/getkin/kin-openapi/cmd/validate@latest

.PHONY: buf-dep-update
buf-dep-update:
	cd proto && buf dep update

.PHONY: proto-lint
proto-lint: buf-lint api-lint

.PHONY: openapi-lint
openapi-lint: openapiv3-lint openapiv2-lint

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
	golangci-lint run --config .golangci.yml

.PHONY: openapiv2-lint
openapiv2-lint:
	# go run github.com/daveshanley/vacuum@latest lint -d proto/gen/openapiv2/gomicroservice/**/*.json
	find proto/gen/openapiv2/gomicroservice -name "*.json" -exec go run github.com/getkin/kin-openapi/cmd/validate@latest -- {} \; && echo "OpenAPI v2 validation successful"

.PHONY: openapiv3-lint
openapiv3-lint:
	# go run github.com/daveshanley/vacuum@latest lint -d proto/gen/openapiv3/*.yaml
	go run github.com/getkin/kin-openapi/cmd/validate@latest -- proto/gen/openapiv3/openapi.yaml && echo "OpenAPI v3 validation successful"

.PHONY: buf-generate
buf-generate:
	cd proto && buf lint && buf generate

.PHONY: tests
tests:
	go test -v ./...

.PHONY: run-server
run-server:
	GO_ENV=development go run cmd/server/main.go

.PHONY: run-server-prod
run-server-prod:
	GO_ENV=production go run cmd/server/main.go

.PHONY: openapi-docs-serve
openapi-docs-serve:
	docker run --rm -p 8080:80 \
		-e SPEC_URL=/api/openapi.yaml \
		-v $(PWD)/proto/gen/openapiv3:/usr/share/nginx/html/api \
		redocly/redoc

.PHONY: openapi-docs-static
openapi-docs-static:
	docker run --rm \
		-v $(PWD)/proto/gen/openapiv3:/tmp/spec \
		-v $(PWD)/docs:/tmp/out \
		redocly/cli build-docs /tmp/spec/openapi.yaml -o /tmp/out/index.html

.PHONY: openapi-download-swagger-ui-files
openapi-download-swagger-ui-files:
    # Clean
	rm -rf cmd/openapi/swagger-ui/*

	# Download and extract Swagger UI files
	curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.21.0.tar.gz | tar xz

	# Copy just the dist files to your swagger-ui directory
	cp -r swagger-ui-5.21.0/dist/* cmd/openapi/swagger-ui/

	# Clean up the downloaded archive
	rm -rf swagger-ui-5.21.0

    # Copy yaml file to swagger-ui directory
	cp proto/gen/openapiv3/openapi.yaml cmd/openapi/swagger-ui/openapi.yaml

    # Change the default URL to your OpenAPI spec
	sed -i 's#url: "https://petstore.swagger.io/v2/swagger.json",#url: "/api/openapi.yaml",#g' cmd/openapi/swagger-ui/swagger-initializer.js
