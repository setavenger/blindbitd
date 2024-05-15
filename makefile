generate:
	protoc --proto_path=BlindBit-Protos BlindBit-Protos/*.proto --go_out=. --go-grpc_out=.

# at the moment we only have automatic docs for the cobra cli, expand once something similar exists for the daemon
# needs to be built before such that latest RootCmd is in the build
gen_docs:
	go build -C ./cli -o ../bin/gen_doc ./scripts && bin/gen_doc


build:
	go build -o bin/blindbitd .
	go build -C cli -o ../bin/blindbit-cli .


build-cli:
	go build -C cli -o ../bin/blindbit-cli .

build-daemon:
	go build -o bin/blindbitd .

build-devtools:
	go build -C devtools -o ../bin/devtools .
