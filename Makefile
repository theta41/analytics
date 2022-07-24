.PHONY: proto swagger

all: proto swagger test build

test:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out
	rm cover.out

build:
	go mod download
	go build -o main cmd/main.go

swagger:
	swag init -g internal/http/api.go

proto:
	protoc -I proto/ --go_out=. --go-grpc_out=. proto/events.proto
	protoc -I proto/ --go_out=. --go-grpc_out=. proto/auth.proto