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

// syncheightCmd represents the syncheight command
var syncheightCmd = &cobra.Command{
	Use:   "syncheight",
	Short: "Get the last sync height",
	Long:  `Daemon has to be unlocked. Shows to which height the daemon is synced`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)

		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		syncHeightResponse, err := client.SyncHeight(context.Background(), &pb.Empty{})
		if err != nil {
			log.Fatalf("Error getting sync height: %v", err)
			return
		}
		fmt.Printf("Daemon synced to: %s\n", lib.ConvertIntToThousandString(int(syncHeightResponse.Height)))
	},
}

func init() {
	RootCmd.AddCommand(syncheightCmd)
}
