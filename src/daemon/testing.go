package daemon

import (
	"encoding/hex"
	"fmt"
	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/gobip352"
)

func (d *Daemon) LoadTestData() error {
	d.Status = pb.Status_STATUS_STARTING
	scanBytes, _ := hex.DecodeString("78e7fd7d2b7a2c1456709d147021a122d2dccaafeada040cc1002083e2833b09")
	spendBytes, _ := hex.DecodeString("c88567742d5019d7ccc81f6e82cef8ef01997a6a3761cc9166036b580549539b")

	d.Wallet.LoadKeys(gobip352.ConvertToFixedLength32(scanBytes), gobip352.ConvertToFixedLength32(spendBytes))
	address, err := d.Wallet.GenerateAddress()
	if err != nil {
		panic(err)
	}

	fmt.Println(address)
	if len(d.Wallet.LabelsMapping) < 5 {
		for _, labelComment := range exampleLabelComments {
			var label *src.Label
			label, err = d.Wallet.GenerateNewLabel(labelComment)
			if err != nil {
				panic(err)
			}
			fmt.Println(label.Address)
		}
	}

	return nil
}

func (d *Daemon) RunTests() {
	if d.Wallet == nil {
		panic("wallet not set")
	}
	if d.ClientBlindBit == nil {
		panic("client not set")
	}

	d.Status = pb.Status_STATUS_SCANNING
	scanBytes, _ := hex.DecodeString("78e7fd7d2b7a2c1456709d147021a122d2dccaafeada040cc1002083e2833b09")
	spendBytes, _ := hex.DecodeString("c88567742d5019d7ccc81f6e82cef8ef01997a6a3761cc9166036b580549539b")

	d.Wallet.LoadKeys(gobip352.ConvertToFixedLength32(scanBytes), gobip352.ConvertToFixedLength32(spendBytes))
	address, err := d.Wallet.GenerateAddress()
	if err != nil {
		panic(err)
	}
	_, err = d.Wallet.GenerateChangeLabel()
	if err != nil {
		panic(err)
	}

	fmt.Println(address)
	for _, labelComment := range exampleLabelComments {
		var label *src.Label
		label, err = d.Wallet.GenerateNewLabel(labelComment)
		if err != nil {
			panic(err)
		}
		fmt.Println(label.Address)
	}

	// 0 sets fetch to tip?
	err = d.SyncToTip(0)
	if err != nil {
		panic(err)
	}

	var balance uint64

	if len(d.Wallet.UTXOs) > 0 {
		//fmt.Printf("%+v\n", d.Wallet.UTXOs)
		fmt.Println()
		for _, utxo := range d.Wallet.UTXOs {
			if utxo.State == src.StateSpent {
				continue
			}
			balance += utxo.Amount
			baseString := fmt.Sprintf("%x:%04d %016d - time: %d", utxo.Txid, utxo.Vout, utxo.Amount, utxo.Timestamp)
			if utxo.Label == nil {
				fmt.Printf("%s - base\n", baseString)
			} else {
				fmt.Printf("%s - label-%d: %s\n", baseString, d.Wallet.LabelsMapping[utxo.Label.PubKey].M, d.Wallet.LabelsMapping[utxo.Label.PubKey].Comment)
			}
		}
	}

	fmt.Println()

	fmt.Printf("Balance: %d\n", balance)
	signedTx, err := d.SendToRecipients([]*src.Recipient{
		{
			Address:    "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xqh6yy9ncjnqhqxazct0fzh98w7lpkm5fvlepqec2yy0sxlq4j6ccc9c679n",
			Amount:     int64(d.Wallet.UTXOs[0].Amount / 2),
			Annotation: "just casually paying myself",
		},
		{
			// this the 5th label
			Address:    "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xqml0tkdw0vxkg3yqkxfyxgqfa9s0znxagejzpmuljcpwa3700mjaqw8cvja",
			Amount:     int64(d.Wallet.UTXOs[0].Amount / 4),
			Annotation: "paying myself on a label",
		},
	}, 10_000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", signedTx)
	d.Status = pb.Status_STATUS_RUNNING
}
