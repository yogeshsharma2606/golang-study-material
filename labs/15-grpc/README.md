# Lab 15: gRPC

A minimal gRPC service and client defined with Protocol Buffers.

## Objectives

- Author a `.proto` service definition.
- Generate Go stubs with `protoc` and the Go plugins.
- Run a gRPC server and call it from a client.

## Setup

Generated Go files are committed under pi/greet/v1/ so go run works without protoc. Re-run gen.bat after editing the .proto when you have protoc installed.

Install tools (once):

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Install [protoc](https://grpc.io/docs/protoc-installation/) and ensure it is on your `PATH`.

Generate code:

```bash
cd labs/15-grpc
go mod tidy
gen.bat
```

Terminal 1 — server:

```bash
go run ./cmd/server
```

Terminal 2 — client:

```bash
go run ./cmd/client Gopher
```

## Exercises

1. Add an RPC `SayHelloAgain` with a `count` field and implement it on the server.
2. Enable server reflection and explore with `grpcurl`.
3. Add a deadline to the client context and observe `DeadlineExceeded`.
4. Discuss when gRPC is preferable to REST for your team.

## What to take away

- Protobuf defines contracts; generated code keeps client and server in sync.
- gRPC uses HTTP/2, strong typing, and efficient payloads.
- Graceful stop (`GracefulStop`) matters for rolling deploys.

## Cleanup

Stop the server with `Ctrl+C`.

## Related Modules

- RPC and API design modules.
- Lab 14 (HTTP microservices) for comparison.

