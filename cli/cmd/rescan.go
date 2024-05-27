package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// rescanCmd represents the chain command
var (
	height int64

	rescanCmd = &cobra.Command{
		Use:   "rescan",
		Short: "calling this triggers a rescan of the chain from height",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			resp, err := client.ForceRescanFromHeight(context.Background(), &pb.RescanRequest{Height: height})
			if err != nil {
				log.Fatal(err)
			}

			if resp.Success {
				fmt.Println("Rescan triggered successfully")
				return
			} else {
				fmt.Printf("Failed with error: %s", resp.Error)
				return
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(rescanCmd)

	rescanCmd.PersistentFlags().Int64Var(&height, "height", -1, "set the height from which the wallet should scan. -1 syncs to tip from last scan height")

	err := cobra.MarkFlagRequired(rescanCmd.PersistentFlags(), "height")
	if err != nil {
		log.Fatal(err)
	}
}
