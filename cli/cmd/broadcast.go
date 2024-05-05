package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/setavenger/blindbitd/pb"
	"google.golang.org/grpc"
	"log"

	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/spf13/cobra"
)

// broadcastCmd represents the status command
var (
	rawTx string

	broadcastCmd = &cobra.Command{
		Use:   "broadcast",
		Short: "broadcast a raw transaction",
		Long: `This command allows you to broadcast any valid transaction 
to the wider bitcoin network.`,
		Run: func(cmd *cobra.Command, args []string) {
			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			txBytes, err := hex.DecodeString(rawTx)
			if err != nil {
				log.Fatalln(err)
			}

			statusResponse, err := client.BroadcastRawTx(context.Background(), &pb.RawTransaction{RawTx: txBytes})
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("Txid: %s\n", statusResponse.Txid)
		},
	}
)

func init() {
	RootCmd.AddCommand(broadcastCmd)

	broadcastCmd.PersistentFlags().StringVar(&rawTx, "rawtx", "", "transaction to broadcast in hex format ")
	err := cobra.MarkFlagRequired(broadcastCmd.PersistentFlags(), "rawtx")
	if err != nil {
		log.Fatalln(err)
	}
}
