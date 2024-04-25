generate:
	protoc --proto_path=src/proto src/proto/*.proto --go_out=src --go-grpc_out=src