package src

import (
	"encoding/hex"
	"fmt"
	"github.com/setavenger/gobip352"
)

type Daemon struct {
	*Wallet
	*Client
}

func (d *Daemon) Run() {

	scanBytes, _ := hex.DecodeString("78e7fd7d2b7a2c1456709d147021a122d2dccaafeada040cc1002083e2833b09")
	spendBytes, _ := hex.DecodeString("c88567742d5019d7ccc81f6e82cef8ef01997a6a3761cc9166036b580549539b")

	exampleScriptToWatch1, _ := hex.DecodeString("3fb78e99650db117f054546c0a99de47c8ab72f9db618f35510c8d960bfcced1")
	exampleScriptToWatch2, _ := hex.DecodeString("2a68ac94e5d66ad352876826bb3924118df4ca2854655ee0bc1a4512dffc7f80")

	d.Wallet.LoadWalletFromKeys(gobip352.ConvertToFixedLength32(scanBytes), gobip352.ConvertToFixedLength32(spendBytes))

	fmt.Println(d.Wallet.Addresses)
	d.Wallet.PubKeysToWatch = [][32]byte{
		gobip352.ConvertToFixedLength32(exampleScriptToWatch1),
		gobip352.ConvertToFixedLength32(exampleScriptToWatch2),
	}
	if d.Wallet == nil {
		panic("wallet not set")
	}
	if d.Client == nil {
		panic("client not set")
	}

	ownedUTXOs, err := d.syncBlock(232)
	if err != nil {
		panic(err)
	}

	if len(ownedUTXOs) > 0 {
		fmt.Printf("%+v\n", ownedUTXOs[0])
	}

	d.Wallet.UTXOs = append(d.Wallet.UTXOs, ownedUTXOs...)

	rawTx, err := d.Wallet.SimpleSendToRecipient(Recipient{
		Address:    "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xqh6yy9ncjnqhqxazct0fzh98w7lpkm5fvlepqec2yy0sxlq4j6ccc9c679n",
		Amount:     3_000_000_000,
		Annotation: "helloooo world",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%x\n", rawTx)

}
