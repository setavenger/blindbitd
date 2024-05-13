package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"os"
	"text/tabwriter"
)

// listaddressesCmd represents the listaddresses command
var listaddressesCmd = &cobra.Command{
	Use:   "listaddresses",
	Short: "Lists all addresses belonging to the user",
	Long: `Daemon has to be unlocked. Lists all addresses belonging to the user. 
If a user has not created any labels this will always 
return the standard (no label applied) Silent Payments address. 
The daemon returns the addresses in the order the labels were created. 
First the base address then the labels with increasing m 1,2,... 
The standard address will always be shown as "standard".
The change label address will not be returned.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(conn)

		addresses, err := client.ListAddresses(context.Background(), &pb.Empty{})
		if err != nil {
			log.Fatalln(err)
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, err = fmt.Fprintln(writer, "Address\tComment")
		if err != nil {
			log.Fatalln(err)
		}
		for _, address := range addresses.Addresses {
			_, err = fmt.Fprintf(writer, "%s\t%s\n", address.Address, address.Comment)
			if err != nil {
				log.Fatalln(err)
			}
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(listaddressesCmd)
}
