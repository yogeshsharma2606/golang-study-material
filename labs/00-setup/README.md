# Lab 00: Setup Verification

## Objectives

- Install Go and verify your toolchain (go version, go env).
- Initialize a module with go mod init and run a minimal program.
- Confirm you can use go run, go build, and go test from the lab directory.

## Setup

1. Install Go from [https://go.dev/dl/](https://go.dev/dl/) (1.22+ recommended).
2. Open a terminal in labs/00-setup/hello/.
3. This repo includes go.mod; to recreate: go mod init github.com/golang-study/00-setup.

## Exercises

1. Run go version and note your Go version.
2. Run go run . and confirm output is Hello, Go!.
3. Run go build -o hello ., run the binary, then remove it.
4. Run go env GOPATH GOMODCACHE and read go help environment for one variable you did not know.

## What to take away

- package main with unc main() is the executable entry point.
- go.mod names the module; go run . builds and runs the main package in the current directory.

## Cleanup

Delete any binary you built (hello or hello.exe).

## Related Modules

- [Go fundamentals](../../modules/01-go-fundamentals.md)
- [Tooling, CI/CD, containers](../../modules/28-tooling-cicd-containers.md)