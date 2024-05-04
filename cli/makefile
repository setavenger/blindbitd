generate:
	protoc --proto_path=lib/BlindBit-Protos lib/BlindBit-Protos/*.proto --go_out=lib --go-grpc_out=lib

# at the moment we only have automatic docs for the cobra cli, expand once something similar exists for the daemon
# needs to be built before such that latest RootCmd is in the build
gen_docs:
	go build -o bin/gen_doc ./cli/scripts && bin/gen_doc