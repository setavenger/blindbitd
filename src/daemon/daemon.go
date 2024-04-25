package daemon

import (
	"encoding/hex"
	"fmt"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/networking"
	"github.com/setavenger/blindbitd/src/pb"
	"github.com/setavenger/gobip352"
)

type Daemon struct {
	Status pb.Status
	*src.Wallet
	*networking.Client
}

var exampleLabelComments = [5]string{"Hello", "Donations for project", "Family and Friends", "Deal 1", "Deal 2"}

func (d *Daemon) Run() {
	d.Status = pb.Status_STATUS_SCANNING
	scanBytes, _ := hex.DecodeString("78e7fd7d2b7a2c1456709d147021a122d2dccaafeada040cc1002083e2833b09")
	spendBytes, _ := hex.DecodeString("c88567742d5019d7ccc81f6e82cef8ef01997a6a3761cc9166036b580549539b")

	exampleScriptToWatch1, _ := hex.DecodeString("3fb78e99650db117f054546c0a99de47c8ab72f9db618f35510c8d960bfcced1")
	exampleScriptToWatch2, _ := hex.DecodeString("2a68ac94e5d66ad352876826bb3924118df4ca2854655ee0bc1a4512dffc7f80")

	d.Wallet.LoadWalletFromKeys(gobip352.ConvertToFixedLength32(scanBytes), gobip352.ConvertToFixedLength32(spendBytes))
	address, err := d.GenerateAddress()
	if err != nil {
		panic(err)
	}
	_, err = d.Wallet.GenerateChangeLabel()
	if err != nil {
		panic(err)
	}

	fmt.Println(address)
	for _, labelComment := range exampleLabelComments {
		address, err = d.GenerateNewLabel(labelComment)
		if err != nil {
			panic(err)
		}
		fmt.Println(address)
	}

	//fmt.Println(d.Wallet.Addresses)
	//fmt.Printf("%+v\n", d.LabelsMapping)

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

	var ownedUTXOs []src.OwnedUTXO
	ownedUTXOs, err = d.syncBlock(235)
	if err != nil {
		panic(err)
	}

	d.Wallet.UTXOs = append(d.Wallet.UTXOs, ownedUTXOs...)

	ownedUTXOs, err = d.syncBlock(236)
	if err != nil {
		panic(err)
	}

	d.Wallet.UTXOs = append(d.Wallet.UTXOs, ownedUTXOs...)

	if len(d.Wallet.UTXOs) > 0 {
		//fmt.Printf("%+v\n", d.Wallet.UTXOs)
		fmt.Println()
		for _, utxo := range d.Wallet.UTXOs {
			baseString := fmt.Sprintf("%x:%04d %016d", utxo.Txid, utxo.Vout, utxo.Amount)
			if utxo.Label == nil {
				fmt.Printf("%s - base\n", baseString)
			} else {
				fmt.Printf("%s - label-%d: %s\n", baseString, d.Wallet.LabelsMapping[utxo.Label.PubKey].M, d.Wallet.LabelsMapping[utxo.Label.PubKey].Comment)
			}
		}
	} // 0000002499982400
	//  2100000000000000
	fmt.Println()

	signedTx, err := d.Wallet.SendToRecipients([]*src.Recipient{
		{
			Address:    "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xqh6yy9ncjnqhqxazct0fzh98w7lpkm5fvlepqec2yy0sxlq4j6ccc9c679n",
			Amount:     int64(d.Wallet.UTXOs[0].Amount / 2),
			Annotation: map[string]any{"label": "just casually paying myself"},
		},
		{
			// this the 5th label
			Address:    "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xqml0tkdw0vxkg3yqkxfyxgqfa9s0znxagejzpmuljcpwa3700mjaqw8cvja",
			Amount:     int64(d.Wallet.UTXOs[0].Amount / 4),
			Annotation: map[string]any{"label": "paying myself on a label"},
		},
	}, 10_000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", signedTx)
	d.Status = pb.Status_STATUS_RUNNING
}
