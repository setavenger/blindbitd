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

// getchainCmd represents the chain command
var getchainCmd = &cobra.Command{
	Use:   "getchain",
	Short: "Gets the chain on which the daemon is running",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		chain, err := client.GetChain(context.Background(), &pb.Empty{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Daemon running on:", chain.Chain)
	},
}

func init() {
	RootCmd.AddCommand(getchainCmd)
}
