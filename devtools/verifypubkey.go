package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/setavenger/go-bip352"
	"github.com/spf13/cobra"
)

var (
	secretKeyHex string
	publicKeyHex string
	negate       bool

	verifyPubKeyCmd = &cobra.Command{
		Use:   "verify-pubkey",
		Short: "Verifies a secret key generates the desired public key",
		Long:  `Checks whether a given secret key belongs to a given public key`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("secret key", secretKeyHex)
			fmt.Println("public key", publicKeyHex)

			secretKeyBytes, err := hex.DecodeString(secretKeyHex)
			if err != nil {
				log.Fatalln(err)
			}
			publicKeyBytes, err := hex.DecodeString(publicKeyHex)
			if err != nil {
				log.Fatalln(err)
			}

			if len(secretKeyBytes) != 32 {
				log.Fatalf("secret key wrong length. expected 32 got %d\n", len(secretKeyBytes))
			}
			//message := []byte("message")
			//aux := []byte("random auxiliary data")

			//Hashing message and auxiliary data
			//msgHash := sha256.Sum256(message)
			//auxHash := sha256.Sum256(aux)
			if negate {
				secret := bip352.NegateSecretKey(bip352.ConvertToFixedLength32(secretKeyBytes))
				secretKeyBytes = secret[:]
			}
			generatedSecretKey, generatedPublicKey := btcec.PrivKeyFromBytes(secretKeyBytes[:])

			_ = generatedSecretKey

			compressedBytes := generatedPublicKey.SerializeCompressed()

			fmt.Printf("%x\n", compressedBytes)
			fmt.Printf("%x\n", publicKeyBytes)

			if bytes.Equal(publicKeyBytes, compressedBytes[1:]) {
				fmt.Println("x-only public keys match")
			} else {
				fmt.Println("x-only public keys do not match")
				return
			}
			// prepend even parity so we know the parity matches as well
			publicKeyBytes = append([]byte{0x02}, publicKeyBytes...)
			if bytes.Equal(publicKeyBytes, compressedBytes) {
				fmt.Println("Parity matches as well")
			} else {
				fmt.Println("Parity does not match")
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(verifyPubKeyCmd)

	verifyPubKeyCmd.PersistentFlags().StringVar(&secretKeyHex, "seckey", "", "Set the secret key which generates the signature (hex)")
	verifyPubKeyCmd.PersistentFlags().StringVar(&publicKeyHex, "pubkey", "", "Set the publicKey to verify against (32 byte hex)")
	verifyPubKeyCmd.PersistentFlags().BoolVar(&negate, "negate", false, "set this flag if the private key should be negated before comparison")
}
