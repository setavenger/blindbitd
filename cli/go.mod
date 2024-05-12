module github.com/setavenger/blindbitd/cli

go 1.21

toolchain go1.21.9

require (
	github.com/setavenger/blindbitd v0.0.0-20240504101325-39a1e38ef9e3
	github.com/spf13/cobra v1.8.0
	golang.org/x/crypto v0.22.0
	golang.org/x/text v0.14.0
	google.golang.org/grpc v1.63.2
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/term v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/setavenger/blindbitd => ../
