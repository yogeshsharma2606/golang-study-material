@echo off
REM Generate Go code from greet.proto (requires protoc on PATH).
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/greet/v1/greet.proto
