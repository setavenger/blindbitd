package main

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/src"
	"github.com/spf13/cobra"
)

var (
	newMnemonicCmd = &cobra.Command{
		Use:   "newmnemonic",
		Short: "Generates a new mnemonic.",
		Long:  `Generate a new mnemonic.`,
		Run: func(cmd *cobra.Command, args []string) {

			src.ChainParams = &chaincfg.MainNetParams

			keys, err := src.CreateNewKeys("")
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Println(keys.Mnemonic)
		},
	}
)

func init() {
	RootCmd.AddCommand(newMnemonicCmd)
}
