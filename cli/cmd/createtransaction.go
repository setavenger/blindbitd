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

// createtransactionCmd represents the createtransaction command
var (
	address    string
	amount     uint64
	feeRate    int64
	annotation string

	broadcast bool

	createtransactionCmd = &cobra.Command{
		Use:   "createtransaction",
		Short: "Construct a transaction",
		Long: "This command can be used to create a transaction sending to exactly one address.\n" +
			"Use `createtransactionmany` if you want to send to several recipients.\n" +
			"By default this will output the raw transaction hex.\n" +
			"Setting the `--broadcast` flag will automatically broadcast the transaction.\n" +
			"Then the command will output the txid of the created transaction.",
		Run: func(cmd *cobra.Command, args []string) {
			if address == "" {
				log.Fatalln("address is required, got:", address)
			}
			if amount == 0 {
				log.Fatalln("amount is required, got:", amount)
			}
			if feeRate == 0 {
				log.Fatalln("feeRate is required, got:", amount)
			}

			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			transactionParams := &pb.CreateTransactionRequest{
				Recipients: []*pb.TransactionRecipient{
					{
						Address:    address,
						Amount:     amount,
						Annotation: annotation,
					},
				},
				FeeRate: feeRate,
			}

			if broadcast {
				txid, err := client.CreateTransactionAndBroadcast(context.Background(), transactionParams)
				if err != nil {
					log.Fatalln("Error:", err)
				}
				fmt.Printf("txid: %s\n", txid.Txid)
			} else {
				transaction, err := client.CreateTransaction(context.Background(), transactionParams)
				if err != nil {
					log.Fatalln("Error:", err)
				}

				fmt.Printf("rawTx: %x\n", transaction.RawTx)

			}
		},
	}
)

func init() {
	RootCmd.AddCommand(createtransactionCmd)

	createtransactionCmd.PersistentFlags().StringVar(&address, "addr", "", "address you want to send to")
	createtransactionCmd.PersistentFlags().Uint64Var(&amount, "amt", 0, "amount you want to send to the address in satoshis [1 BTC = 100,000,000 sats]")
	createtransactionCmd.PersistentFlags().Int64Var(&feeRate, "sat_per_byte", 0, "set the fee rate (in sats/vByte) for the transaction. Has to be an integer")
	//createtransactionCmd.PersistentFlags().StringVar(&annotation, "annotation", "", "add an annotation the recipient")  // todo not used in a meaningful way in daemon yet
	createtransactionCmd.PersistentFlags().BoolVar(&broadcast, "broadcast", false, "broadcasts the transaction directly")

	// required flags
	err := cobra.MarkFlagRequired(createtransactionCmd.PersistentFlags(), "addr")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(createtransactionCmd.PersistentFlags(), "amt")
	if err != nil {
		log.Fatalln(err)
	}
	err = cobra.MarkFlagRequired(createtransactionCmd.PersistentFlags(), "sat_per_byte")
	if err != nil {
		log.Fatalln(err)
	}
}
