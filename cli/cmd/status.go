package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/pb"
	"google.golang.org/grpc"

	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of the daemon",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		statusResponse, err := client.Status(context.Background(), &pb.Empty{})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Status: %v\n", statusResponse.Status)
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
