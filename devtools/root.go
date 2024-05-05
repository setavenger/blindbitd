package main

import (
	"encoding/hex"
	"fmt"
	"github.com/setavenger/go-bip352"
	"log"
	"os"

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

	RootCmd = &cobra.Command{
		Use:   "blindbit-cli",
		Short: "A simple cli application to compute an silent payment output",
		Long: `This is a simple cli tool that can create a silent payment x-only output. 
It only works with one input and one output.
The required arguments are a target address, 
the inputs private key, txid and vout.`,
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
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&address, "addr", "", "Set the recipients address")
	RootCmd.PersistentFlags().StringVar(&privKeyHex, "secret", "", "Set the inputs secret key")
	RootCmd.PersistentFlags().StringVar(&txidHex, "txid", "", "Set the inputs txid")
	RootCmd.PersistentFlags().Uint32Var(&vout, "vout", 0, "Set the inputs vout")
	RootCmd.PersistentFlags().BoolVar(&isTaproot, "taproot", false, "needs to be set if the sending private key is a taproot key")
	RootCmd.PersistentFlags().BoolVar(&mainnet, "main", false, "if flag is set everything is parsed for mainnet (not recommended)")

	// required flags
	err := cobra.MarkFlagRequired(RootCmd.PersistentFlags(), "addr")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(RootCmd.PersistentFlags(), "secret")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(RootCmd.PersistentFlags(), "txid")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(RootCmd.PersistentFlags(), "vout")
	if err != nil {
		log.Fatalln(err)
	}

}
