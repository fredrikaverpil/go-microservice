# go-microservice

## Captain's log

### Buf and protos setup

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
