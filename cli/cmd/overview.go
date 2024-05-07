package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/pb"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"

	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Get an overview over your wallet",
	Long: `Displays the status of the daemon, 
the height to which the daemon is synced, 
the chain on which the daemon is running,
the current balance, and the first 5 addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := lib.NewClient(socketPath)
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		}(conn)

		var wg sync.WaitGroup

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		status, err := client.Status(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalln(err)
		}

		// if the daemon is not locked or has no wallet
		var daemonAvailable bool

		switch status.Status {
		case pb.Status_STATUS_NO_WALLET, pb.Status_STATUS_LOCKED:
		default:
			daemonAvailable = true
		}

		if !daemonAvailable {
			// we don't get the other data and end here
			fmt.Println("--- Blindbit Daemon Overview ---")
			fmt.Printf("Daemon Status:  %v\n", status)
			return
		}
		wg.Add(4)

		var balance uint64
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			utxoCollectionResp, err := client.ListUTXOs(ctx, &pb.Empty{})
			if err != nil {
				log.Fatal(err)
			}
			for _, utxo := range utxoCollectionResp.Utxos {
				if utxo.UtxoState == pb.UTXOState_UNSPENT {
					balance += utxo.Amount
				}
			}
		}()

		var addressesResp *pb.AddressesCollection
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			addressesResp, err = client.ListAddresses(ctx, &pb.Empty{})
			if err != nil {
				log.Fatal(err)
			}
		}()

		var syncHeightResp *pb.SyncHeightResponse
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			syncHeightResp, err = client.SyncHeight(ctx, &pb.Empty{})
			if err != nil {
				log.Fatal(err)
			}
		}()

		var chainResp *pb.Chain
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			chainResp, err = client.GetChain(ctx, &pb.Empty{})
			if err != nil {
				log.Fatal(err)
			}
		}()
		wg.Wait()

		fmt.Println("--- Blindbit Daemon Overview ---")
		fmt.Printf("Daemon Status:  %v\n", status.Status)
		fmt.Printf("Sync Height:    %s\n", lib.ConvertIntToThousandString(int(syncHeightResp.Height)))
		fmt.Printf("Running Chain:  %s\n", chainResp.Chain)
		fmt.Printf("Balance:        %s\n", lib.ConvertIntToThousandString(int(balance)))

		addresses := addressesResp.GetAddresses()
		fmt.Println("\nRegistered Addresses:")
		for i, addr := range addresses {
			if i >= 5 {
				break
			}
			fmt.Printf("%d. %s - %s\n", i+1, addr.Address, addr.Comment)
		}
	},
}

func init() {
	RootCmd.AddCommand(overviewCmd)
}
