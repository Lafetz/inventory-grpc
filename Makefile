run-inventory:
	go run cmd/inventory/main.go
proto-inventory:
	protoc --go_out=. --go_opt=paths=source_relative \
  	  --go-grpc_out=. --go-grpc_opt paths=source_relative \
	  proto/inventory.proto
.PHONY: proto 