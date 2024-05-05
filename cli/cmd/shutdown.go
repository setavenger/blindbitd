package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"google.golang.org/grpc"
	"log"

	"github.com/spf13/cobra"
)

// shutdownCmd represents the shutdown command
var shutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shuts down the daemon",
	Long:  `Daemon has to be unlocked. This command shuts down the daemon.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		response, err := client.Shutdown(context.Background(), &pb.Empty{})
		if err != nil {
			fmt.Println(err)
			return
		}
		if response.Success {
			fmt.Println("Daemon is shutting down")
			return
		}
		log.Fatalf("Error: %v\n", response.Error)
	},
}

func init() {
	RootCmd.AddCommand(shutdownCmd)
	// todo add graceful option if necessary
}
