# go-microservice

## Quickstart

```bash
# install deps
make
# run server
make run-server
```

### Calling endpoints

#### gRPC APIs with grpcurl

```bash
# List all available services
grpcurl -plaintext localhost:50051 list

# Create a user with server-assigned ID
grpcurl -plaintext -d '{"user":{"display_name":"John Doe","email":"john@example.com"}}' \
  localhost:50051 gomicroservice.v1.UserService/CreateUser

# Create a user with client-assigned ID
grpcurl -plaintext -d '{"user":{"name":"users/user123","display_name":"John Doe","email":"john@example.com"}}' \
  localhost:50051 gomicroservice.v1.UserService/CreateUser

# Get a user
grpcurl -plaintext -d '{"name":"users/user123"}' \
  localhost:50051 gomicroservice.v1.UserService/GetUser

# List users
grpcurl -plaintext -d '{"page_size":10}' \
  localhost:50051 gomicroservice.v1.UserService/ListUsers

# Update a user
grpcurl -plaintext -d '{"user":{"name":"users/user123","display_name":"John Smith"}}' \
  localhost:50051 gomicroservice.v1.UserService/UpdateUser

# Delete a user
grpcurl -plaintext -d '{"name":"users/user123"}' \
  localhost:50051 gomicroservice.v1.UserService/DeleteUser
```

#### REST APIs with curl

```bash
# Create a user with server-assigned ID
curl -X POST -H "Content-Type: application/json" \
  -d '{"user":{"display_name":"John Doe","email":"john@example.com"}}' \
  http://localhost:8080/v1/users

# Create a user with client-assigned ID
curl -X POST -H "Content-Type: application/json" \
  -d '{"user":{"name":"users/user123","display_name":"John Doe","email":"john@example.com"}}' \
  http://localhost:8080/v1/users

# Get a user
curl http://localhost:8080/v1/users/user123

# List users
curl http://localhost:8080/v1/users

# Update a user
curl -X PATCH -H "Content-Type: application/json" \
  -d '{"user":{"display_name":"John Smith"},"update_mask":{"paths":["display_name"]}}' \
  http://localhost:8080/v1/users/user123

# Delete a user
curl -X DELETE http://localhost:8080/v1/users/user123
```

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

### To do

#### Buf config

- Add
  [protoc-gen-go-aip-test](https://github.com/einride/protoc-gen-go-aip-test)
- Add [aip-cli-go](https://github.com/einride/aip-cli-go)
