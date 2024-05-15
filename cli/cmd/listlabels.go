package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"text/tabwriter"
	"time"
  "os"
)

var (
	labelsListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all labels",
		Long:  `This command shows all labels. M is the counter of the labels.`,
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
			labels, err := client.ListLabels(ctx, &pb.Empty{})
			if err != nil {
				log.Fatalf("could not retrieve labels: %v\n", err)
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			_, err = fmt.Fprintln(writer, "M\tAddress\tComment")
			if err != nil {
				log.Fatalln(err)
			}
			for _, label := range labels.Labels {
				_, err = fmt.Fprintf(writer, "%d\t%s\t%s\n", label.M, label.Address, label.Comment)
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
)


func init () {}

