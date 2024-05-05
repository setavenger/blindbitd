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

// listaddressesCmd represents the listaddresses command
var listaddressesCmd = &cobra.Command{
	Use:   "listaddresses",
	Short: "Lists all addresses belonging to the user",
	Long: `Daemon has to be unlocked. Lists all addresses belonging to the user. 
If a user has not created any labels this will always 
return the standard silent payment address. 
The daemon returns the addresses in the order the labels were created. 
1. the base address then the labels with increasing m 1,2,... 
The change label address will not be returned.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		addresses, err := client.ListAddresses(context.Background(), &pb.Empty{})
		if err != nil {
			log.Fatalln(err)
		}
		for _, address := range addresses.Addresses {
			fmt.Printf("%s - %s\n", address.Address, address.Comment)
		}
	},
}

func init() {
	RootCmd.AddCommand(listaddressesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listaddressesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listaddressesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
