package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/setavenger/go-bip352"

	"github.com/spf13/cobra"
)

type VinArgInput struct {
	PrivKeyHex string `json:"priv_key_hex"`
	TxidHex    string `json:"txid_hex"`
	Vout       uint32 `json:"vout"`
	IsTaproot  bool   `json:"is_taproot"`
}

// RootCmd represents the base command when called without any subcommands
var (
	address2 string
	vinArgs  []string

	createOutputsCmd = &cobra.Command{
		Use:   "create-outputs",
		Short: "Create a Silent Payments output based on multiple inputs",
		Long: `This is a simple cli command that can create a silent payment x-only output. 
It works with several inputs and one output address.
The required arguments are a target address, 
    and the vins need to be in string json format (can be several vins)
If the sending output is a taproot output, is_taproot field needs to be set to true.`,
		Run: func(cmd *cobra.Command, args []string) {
			// convert to bytes

			// api requires array
			recipient := []*bip352.Recipient{{SilentPaymentAddress: address2}}

			vins, err := parseVinArgs(vinArgs)
			if err != nil {
				log.Fatalf("err parsing vin args: %v\n", err)
			}

			err = bip352.SenderCreateOutputs(recipient, vins, mainnet, false)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			fmt.Printf("output: %x\n", recipient[0].Output)
		},
	}
)

func parseVinArgs(vinArgs []string) (vins []*bip352.Vin, err error) {
	for i := range vinArgs {
		var vin *bip352.Vin
		vin, err = parseVinArg(vinArgs[i])
		if err != nil {
			fmt.Printf("unable to parse vinArg: %v\n", err)
			return nil, err
		}
		vins = append(vins, vin)

	}
	return
}

func parseVinArg(s string) (vin *bip352.Vin, err error) {
	var input VinArgInput
	err = json.Unmarshal([]byte(s), &input)
	if err != nil {
		fmt.Printf("unable to unmarshal VinArgInput: %v\n", err)
		return
	}

	privKeyBytes, err := hex.DecodeString(input.PrivKeyHex)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	txidBytes, err := hex.DecodeString(input.TxidHex)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	txid := bip352.ConvertToFixedLength32(txidBytes)
	privKey := bip352.ConvertToFixedLength32(privKeyBytes)

	return &bip352.Vin{
		Txid:      txid,
		Vout:      input.Vout,
		SecretKey: &privKey,
		Taproot:   input.IsTaproot,
	}, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.

func init() {
	RootCmd.AddCommand(createOutputsCmd)
	createOutputsCmd.PersistentFlags().StringVar(&address2, "addr", "", "Set the recipients address")
	createOutputsCmd.PersistentFlags().StringArrayVar(&vinArgs, "vins", []string{}, "vins as json representation with the fields: priv_key_hex, txid_hex, vout, is_taproot")
	createOutputsCmd.PersistentFlags().BoolVar(&mainnet, "mainnet", false, "if flag is set everything is parsed for mainnet (not recommended)")

	// required flags
	err := cobra.MarkFlagRequired(createOutputsCmd.PersistentFlags(), "addr")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(createOutputsCmd.PersistentFlags(), "vins")
	if err != nil {
		log.Fatalln(err)
	}
}
