# gRPC authentication service

# Usage

Clone repo:
```shell
git clone https://github.com/dev-yeva/auth-service
cd auth-service
go -C ./auth mod download
```

Use task to run auth service and client:
1. Install [task](https://taskfile.dev/docs/installation)
2. Run server
```shell
task run-server
```
3. Run client
```shell
task run-client
```