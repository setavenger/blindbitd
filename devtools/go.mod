module blindbitd/debug-tools

go 1.21.9 // todo unify with the other modules

require (
	github.com/setavenger/go-bip352 v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/btcsuite/btcd v0.23.5-0.20231215221805-96c9fd8078fd // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.5 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

replace github.com/setavenger/go-bip352 => ../../go-bip352
