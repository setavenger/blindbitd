package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
)

// getmnemonicCmd represents the chain command
var getmnemonicCmd = &cobra.Command{
	Use:   "getmnemonic",
	Short: "CAUTION: Shows the wallets mnemonic",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		resp, err := client.GetMnemonic(context.Background(), &pb.Empty{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Your mnemonic:", resp.Mnemonic)
	},
}

func init() {
	RootCmd.AddCommand(getmnemonicCmd)
	// todo add graceful option if necessary
}
