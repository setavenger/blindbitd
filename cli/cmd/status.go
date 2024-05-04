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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
