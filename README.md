# go-microservice

## Captain's log

### Buf and protos setup

#### Write protos

- Add proto(s) into `proto/gomicroservice/v1`.
- Read the
  [buf tutorial](https://buf.build/docs/tutorials/getting-started-with-buf-cli#update-directory-path-and-build-module).
- `cd protos/gomicroservice/v1 && buf config init`
- Run `buf dep update` to download/update the dependencies and create the
  `buf.lock` file.
- Add the `modules` key to `buf.yaml`, which dictates where proto code is
  generated.
- Run `buf lint` to make sure the proto is valid. Add dependencies and modify
  proto file as to satisfy the linter and AIP rules.

#### Generate go code

- Install dependencies:

  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest  # non-gRPC go code
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest # gRPC go code
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
  ```

- Add plugins and configure them in `buf.gen.yaml`.
- Run `buf generate` to run codegen.
- Add api-linter (for AIP). Add `go install ...` command.
- Run linter and add the `client.proto` import for api methods, fix various
  issues.

#### To do

- Add AIP tests.
