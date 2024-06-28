package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	getScanOnlyDataCmd = &cobra.Command{
		Use:   "getscanonlydata",
		Short: "Shows the scan secret key and the spend public key which can then be passed to a scan only daemon.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := client.GetScanOnlyData(ctx, &pb.Empty{})
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("Secret key scan: %x\n", result.ScanSecretKey)
			fmt.Printf("Spend Public Key: %x\n", result.SpendPublicKey)
		},
	}
)

func init() {
	RootCmd.AddCommand(getScanOnlyDataCmd)
}
