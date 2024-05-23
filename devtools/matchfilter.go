package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/gcs"
	"github.com/btcsuite/btcutil/gcs/builder"
	"github.com/setavenger/go-bip352"
	"github.com/spf13/cobra"
)

var (
	filterNData  string
	valueToMatch string
	key          string

	format string = "hex"

	matchFilterCmd = &cobra.Command{
		Use:   "match-filter",
		Short: "match a value against a GCS filter",
		Long: `This command allows you to find out whether a value exists in a filter.
    The input format has to be hex for both the filter and the value.`,
		Run: func(cmd *cobra.Command, args []string) {

			filterData, err := hex.DecodeString(filterNData)
			if err != nil {
				log.Fatalln(err)
			}
			keyBytes, err := hex.DecodeString(key)
			if err != nil {
				log.Fatalln(err)
			}

			c, err := chainhash.NewHash(bip352.ReverseBytes(keyBytes))
			if err != nil {
				log.Fatalln(err)
			}

			value, err := hex.DecodeString(valueToMatch)
			if err != nil {
				log.Fatalln(err)
			}

			var valuesToMatch = [][]byte{value}

			filter, err := gcs.FromNBytes(builder.DefaultP, builder.DefaultM, filterData)
			if err != nil {
				log.Fatalln(err)
			}

			key := builder.DeriveKey(c)

			isMatch, err := filter.HashMatchAny(key, valuesToMatch)
			if err != nil {
				log.Fatalln(err)
			}

			if isMatch {
				fmt.Println("value matched")
			} else {
				fmt.Println("value not found")
			}

		},
	}
)

func init() {
	RootCmd.AddCommand(matchFilterCmd)
	matchFilterCmd.PersistentFlags().StringVar(&filterNData, "ndata", "", "Set the data of the filter needs to be nbytes (hex)")
	matchFilterCmd.PersistentFlags().StringVar(&valueToMatch, "value", "", "Set the value that should be matched (hex)")
	matchFilterCmd.PersistentFlags().StringVar(&key, "key", "", "Set the key for the filter the key is normally a blockhash of which the first 16 bytes are used (hex)")

	// required flags
	err := cobra.MarkFlagRequired(matchFilterCmd.PersistentFlags(), "ndata")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(matchFilterCmd.PersistentFlags(), "value")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(matchFilterCmd.PersistentFlags(), "key")
	if err != nil {
		log.Fatalln(err)
	}
}
