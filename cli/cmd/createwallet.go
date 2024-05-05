package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// createwalletCmd represents the createwallet command
var (
	showMnemonic      bool
	useSeedPassphrase bool
	createwalletCmd   = &cobra.Command{
		Use:   "createwallet",
		Short: "Create a new wallet",
		Long: `Create a new wallet in the daemon. This should fail if the daemon already contains a wallet.

The encryption password is only to encrypt your wallet data (keys, utxos, etc.) on disk. It is not used for your seed. 
To add a passphrase to your seed set the --seedpass flag (not extensively tested yet)
`,
		Run: func(cmd *cobra.Command, args []string) {
			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			fmt.Println(
				"NOTE: The encryption password is only to encrypt your wallet data (keys, utxos, etc.) on disk." +
					"\nIt is not used for your seed. To add a passphrase to your seed set the --seedpass flag (not extensively tested yet)",
			)
			fmt.Print("Encryption password: ")
			passwordBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalln("Error reading password")
			}
			fmt.Println()

			var seedPassphrase string
			if useSeedPassphrase {
				fmt.Println("Enter your seed passphrase: ")
				seedPassphraseBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					log.Fatalln("Error reading seed passphrase")
				}
				seedPassphrase = string(seedPassphraseBytes)
				fmt.Println()
			}

			response, err := client.CreateNewWallet(context.Background(), &pb.NewWalletRequest{EncryptionPassword: string(passwordBytes), SeedPassphrase: seedPassphrase})
			if err != nil {
				fmt.Println(err)
				return
			}

			if response.Mnemonic == "" {
				log.Fatalln("Error: mnemonic was empty without throwing error. Check daemon logs.")
			}
			if showMnemonic {
				fmt.Println()
				fmt.Println("Mnemonic:", response.Mnemonic)
				fmt.Println("The 12/24 words above are the access to your funds.\nDon't publish or share them, you risk loosing your funds.")
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(createwalletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	createwalletCmd.PersistentFlags().BoolVar(&showMnemonic, "show", false, "show the wallet seed phrase after wallet creation")
	createwalletCmd.PersistentFlags().BoolVar(&useSeedPassphrase, "seedpass", false, "add a passphrase to the wallet seed")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createwalletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
