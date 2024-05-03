generate:
	protoc --proto_path=src/BlindBit-Protos src/BlindBit-Protos/*.proto --go_out=src --go-grpc_out=src