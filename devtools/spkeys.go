package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/src"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	spKeysCmd = &cobra.Command{
		Use:   "spkeys",
		Short: "get the scan and spend secret key for a mnemonic",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Print("Mnemonic: ")
			mnemonicBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln("Error reading mnemonic")
			}
			fmt.Println()

			if mainnet {
				src.ChainParams = &chaincfg.MainNetParams
			} else {
				// should not matter which non mainnet chain we use
				src.ChainParams = &chaincfg.SigNetParams
			}

			var passphrase string
			if seedPassPassphrase {
				fmt.Print("Seed passphrase: ")
				passBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					log.Fatalln("Error reading seed passphrase")
				}
				fmt.Println()
				fmt.Print("Confirm seed passphrase: ")
				passBytesConf, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					log.Fatalln("Error reading seed passphrase")
				}
				fmt.Println()

				if !bytes.Equal(passBytes, passBytesConf) {
					log.Fatalln("passphrases did not match")
				}
				passphrase = string(passBytes)

			}

			keys, err := src.KeysFromMnemonic(string(mnemonicBytes), passphrase)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("Scan secret key:   %x\n", keys.ScanSecretKey)
			fmt.Printf("Spend secret key:  %x\n", keys.SpendSecretKey)
		},
	}
)

func init() {
	RootCmd.AddCommand(spKeysCmd)

	spKeysCmd.PersistentFlags().BoolVar(&mainnet, "mainnet", false, "set flag if mainnet should be used for key derivation")
	spKeysCmd.PersistentFlags().BoolVar(&seedPassPassphrase, "passphrase", false, "set flag if you need to enter a seed passphrase")
}
