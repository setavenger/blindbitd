package cmd

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"google.golang.org/grpc"
)

var (
	spendPubKey  string
	birthHeight2 uint64 // todo how do we handle reoccurring var names
	labelCount2  uint32

	setupScanOnlyCmd = &cobra.Command{
		Use:   "setupscanonly",
		Short: "Call this command to set up a daemon in scan only mode.",
		Long: `The daemon has to be set to scan only.
    You can pass --birthheight and --labelcount to fine tune the scanning.`,
		Run: func(cmd *cobra.Command, args []string) {
			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			pubKey, err := hex.DecodeString(spendPubKey)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(
				"NOTE: The encryption password is only to encrypt your wallet data (keys, utxos, etc.) on disk.")
			fmt.Print("Encryption password: ")
			passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln("Error reading password")
			}
			fmt.Println()
			fmt.Print("Confirm password: ")
			passworConfirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln("Error reading password")
			}
			fmt.Println()
			if !bytes.Equal(passwordBytes, passworConfirmBytes) {
				log.Fatalln("Passwords do not match")
			}

			fmt.Println(
				"NOTE: enter your scan secret key in hex format")
			fmt.Print("Secret Key scan: ")
			secretKeyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln("Error reading secret key scan")
			}
			fmt.Println()
			fmt.Print("Confirm secret key scan: ")
			secretKeyConfirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln("Error reading secret key scan")
			}
			fmt.Println()
			if !bytes.Equal(secretKeyBytes, secretKeyConfirmBytes) {
				log.Fatalln("secret keys do not match")
			}

			secretKey, err := hex.DecodeString(string(secretKeyBytes))
			if err != nil {
				log.Fatalln(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			response, err := client.SetupScanOnly(ctx, &pb.ScanOnlySetupRequest{
				EncryptionPassword: string(passwordBytes),
				ScanSecretKey:      secretKey,
				SpendPublicKey:     pubKey,
				BirthHeight:        &birthHeight2,
				LabelCount:         &labelCount2,
			})
			if err != nil {
				log.Fatalln(err)
			}
			if response.Success {
				fmt.Println("Success")
			} else {
				fmt.Println("Failed:", response.Error)
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(setupScanOnlyCmd)

	setupScanOnlyCmd.PersistentFlags().StringVar(&spendPubKey, "spendpub", "", "spend public key for scanning")
	setupScanOnlyCmd.PersistentFlags().Uint64Var(&birthHeight2, "birthheight", 0, "set the birthheight, the daemon will start scanning from that height")
	setupScanOnlyCmd.PersistentFlags().Uint32Var(&labelCount2, "labelcount", 0, "set the number of labels that should be pre-generated, such that all payments are found")
}
