MODULE  := github.com/RomanKovalev007/topqueries-service
PROTO_DIR := api/proto
GEN_DIR   := gen/pb

.PHONY: proto build run test bench

proto:
	mkdir -p $(GEN_DIR)
	protoc \
		--go_out=. \
		--go_opt=module=$(MODULE) \
		--go-grpc_out=. \
		--go-grpc_opt=module=$(MODULE) \
		$(PROTO_DIR)/*.proto

build:
	go build ./cmd/...

run:
	go run ./cmd/main.go

test:
	go test ./...

bench:
	go test ./... -bench=. -benchmem -benchtime=5s -run=^$
