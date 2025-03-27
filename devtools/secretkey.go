package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/setavenger/go-bip352"
	"github.com/spf13/cobra"
)

var (
	secretKey1Hex string
	secretKey2Hex string

	secretKeyCmd = &cobra.Command{
		Use:   "secretkey",
		Short: "operations with secret keys",
		Long:  ``,
	}

	addSecretKeysCmd = &cobra.Command{
		Use:   "add",
		Short: "adds two secret keys",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			secretKey1Bytes, err := hex.DecodeString(secretKey1Hex)
			if err != nil {
				log.Fatalln(err)
			}
			secretKey2Bytes, err := hex.DecodeString(secretKey2Hex)
			if err != nil {
				log.Fatalln(err)
			}

			secretKey1 := bip352.ConvertToFixedLength32(secretKey1Bytes)
			secretKey2 := bip352.ConvertToFixedLength32(secretKey2Bytes)

			fmt.Printf("%x\n", bip352.AddPrivateKeys(secretKey1, secretKey2))
		},
	}

	negateSecretKeyCmd = &cobra.Command{
		Use:   "negate",
		Short: "negates a secret key",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			secretKey1Bytes, err := hex.DecodeString(secretKey1Hex)
			if err != nil {
				log.Fatalln(err)
			}

			secretKey := bip352.ConvertToFixedLength32(secretKey1Bytes)

			fmt.Printf("%x\n", bip352.NegateSecretKey(secretKey))
		},
	}
)

func init() {
	RootCmd.AddCommand(secretKeyCmd)
	secretKeyCmd.AddCommand(addSecretKeysCmd)
	secretKeyCmd.AddCommand(negateSecretKeyCmd)

	secretKeyCmd.PersistentFlags().StringVar(&secretKey1Hex, "seckey1", "", "Set the 1st secret key for the addition")
	secretKeyCmd.PersistentFlags().StringVar(&secretKey2Hex, "seckey2", "", "Set the 2nd secret key for the addition")

	negateSecretKeyCmd.PersistentFlags().StringVar(&secretKey1Hex, "seckey", "", "Set the secret key to be negated")

}
