package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/setavenger/go-bip352"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var (
	address    string
	privKeyHex string
	txidHex    string
	vout       uint32
	mainnet    bool
	isTaproot  bool

	createOutputCmd = &cobra.Command{
		Use:   "create-output",
		Short: "Create a Silent Payments output based on simple inputs",
		Long: `This is a simple cli command that can create a silent payment x-only output. 
It only works with one input and one output.
The required arguments are a target address, 
the inputs private key, txid and vout.
If the sending output is a taproot output 
the --taproot flag HAS to be set.`,
		Run: func(cmd *cobra.Command, args []string) {
			// convert to bytes
			privKeyBytes, err := hex.DecodeString(privKeyHex)
			if err != nil {
				log.Fatalf("error: %s", err)
			}
			txidBytes, err := hex.DecodeString(txidHex)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			txid := bip352.ConvertToFixedLength32(txidBytes)
			privKey := bip352.ConvertToFixedLength32(privKeyBytes)

			// api requires array
			recipient := []*bip352.Recipient{{SilentPaymentAddress: address}}

			vins := []*bip352.Vin{{
				Txid:      txid,
				Vout:      vout,
				Taproot:   isTaproot,
				SecretKey: &privKey,
			}}
			err = bip352.SenderCreateOutputs(recipient, vins, mainnet, false)
			if err != nil {
				log.Fatalf("error: %s", err)
			}

			fmt.Printf("output: %x\n", recipient[0].Output)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.

func init() {
	RootCmd.AddCommand(createOutputCmd)
	createOutputCmd.PersistentFlags().StringVar(&address, "addr", "", "Set the recipients address")
	createOutputCmd.PersistentFlags().StringVar(&privKeyHex, "secret", "", "Set the inputs secret key")
	createOutputCmd.PersistentFlags().StringVar(&txidHex, "txid", "", "Set the inputs txid")
	createOutputCmd.PersistentFlags().Uint32Var(&vout, "vout", 0, "Set the inputs vout")
	createOutputCmd.PersistentFlags().BoolVar(&isTaproot, "taproot", false, "needs to be set if the sending private key is a taproot key")
	createOutputCmd.PersistentFlags().BoolVar(&mainnet, "main", false, "if flag is set everything is parsed for mainnet (not recommended)")

	// required flags
	err := cobra.MarkFlagRequired(createOutputCmd.PersistentFlags(), "addr")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(createOutputCmd.PersistentFlags(), "secret")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(createOutputCmd.PersistentFlags(), "txid")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(createOutputCmd.PersistentFlags(), "vout")
	if err != nil {
		log.Fatalln(err)
	}

}
