package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/setavenger/blindbitd/cli/lib"
	"github.com/setavenger/blindbitd/pb"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

// createtransactionCmd represents the createtransaction command
var (
	addresses   []string
	amounts     []int64
	feeRate     int64
	annotations []string

	broadcast           bool
	notMarkSpent        bool
	useSpentUnconfirmed bool

	createtransactionCmd = &cobra.Command{
		Use:   "createtransaction",
		Short: "Construct a transaction",
		Long: "This command can be used to create a transaction sending to several addresses if needed.\n" +
			"By default this will output the raw transaction hex.\n" +
			"Setting the `--broadcast` flag will automatically broadcast the transaction.\n" +
			"Then the command will output the txid of the created transaction.\n" +
			"UTXOs used in a transaction are automatically marked as spent_unconfirmed.\n" +
			"Use --notmarkspent to not do this.\n" +
			"Use --usespent to include spent_unconfirmed UTXOs in transaction creation.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(addresses) < 1 {
				log.Fatalln("needs at least one address")
			}
			if len(addresses) != len(amounts) {
				log.Fatalf("different number of addresses (%d) and amounts (%d)", len(addresses), len(amounts))
			}
			if len(annotations) > 0 && len(addresses) != len(annotations) {
				log.Fatalf("number annotations (%d) does not match addresses (%d). When using annotations the number of annotations has to be the same as addresses/amounts. Use `--note \"\"` for recipients without annotations.", len(annotations), len(addresses))
			}
			if feeRate == 0 {
				log.Fatalln("feeRate is required, got:", feeRate)
			}

			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			var recipients []*pb.TransactionRecipient
			for i, addr := range addresses {
				recipient := &pb.TransactionRecipient{
					Address: addr,
					Amount:  uint64(amounts[i]),
				}
				if len(annotations) > 0 {
					// we checked the lengths above already
					recipient.Annotation = annotations[i]
				}

				recipients = append(recipients, recipient)
			}

			transactionParams := &pb.CreateTransactionRequest{
				Recipients:          recipients,
				FeeRate:             feeRate,
				MarkSpent:           !notMarkSpent,
				UseSpentUnconfirmed: useSpentUnconfirmed,
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

	createtransactionCmd.PersistentFlags().StringSliceVar(&addresses, "addr", nil, "address you want to send to")
	createtransactionCmd.PersistentFlags().Int64SliceVar(&amounts, "amt", nil, "amount you want to send to the address in satoshis [1 BTC = 100,000,000 sats]")
	createtransactionCmd.PersistentFlags().Int64Var(&feeRate, "sat_per_byte", 0, "set the fee rate (in sats/vByte) for the transaction. Has to be an integer")
	createtransactionCmd.PersistentFlags().StringSliceVar(&annotations, "note", nil, "add annotation to recipient")
	//createtransactionCmd.PersistentFlags().StringVar(&annotation, "annotation", "", "add an annotation the recipient")  // todo not used in a meaningful way in daemon yet
	createtransactionCmd.PersistentFlags().BoolVar(&broadcast, "broadcast", false, "broadcasts the transaction directly")
	createtransactionCmd.PersistentFlags().BoolVar(&notMarkSpent, "notmarkspent", false, "not mark utxos of the transaction as spent_unconfirmed")
	createtransactionCmd.PersistentFlags().BoolVar(&useSpentUnconfirmed, "usespent", false, "include utxos with state spent_unconfirmed")

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
