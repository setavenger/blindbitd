package coinselector

import (
	"bytes"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/logging"
	"testing"
)

func init() {
	logging.LoadLoggersMock()
}

type TestCase struct {
	Comment string
	Given   struct {
		Utxos           src.UtxoCollection
		Recipients      []*src.Recipient
		FeeRate         uint32
		MinChangeAmount uint64
	}
	Expected struct {
		Change             uint64
		NumOfSelectedUTXOs int
		AbsolutFee         uint64
		Err                error
	}
}

// todo unify to use the same utxo set, increases readability
var testCases = []TestCase{
	{
		Comment: "Simple case",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				}, {
					Amount: 40_000,
				}, {
					Amount: 60_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "bc1qua7e852suw0p74e2lzxwmk2tw8fd2zuzexc866",
					Amount:  5_000,
				},
			},
			FeeRate:         1,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 14858, NumOfSelectedUTXOs: 1, AbsolutFee: 142, Err: nil},
	},
	{
		Comment: "fits with change amount",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				}, {
					Amount: 40_000,
				}, {
					Amount: 60_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "bc1qua7e852suw0p74e2lzxwmk2tw8fd2zuzexc866",
					Amount:  54_000,
				},
			},
			FeeRate:         1,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 5800, NumOfSelectedUTXOs: 2, AbsolutFee: 200, Err: nil},
	},
	{
		Comment: "Does it work with different fee rates",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				}, {
					Amount: 40_000,
				}, {
					Amount: 60_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "bc1qua7e852suw0p74e2lzxwmk2tw8fd2zuzexc866",
					Amount:  50_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 8007, NumOfSelectedUTXOs: 2, AbsolutFee: 1993, Err: nil},
	},
	{
		Comment: "fails because not enough funds",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "bc1qua7e852suw0p74e2lzxwmk2tw8fd2zuzexc866",
					Amount:  20_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 0, NumOfSelectedUTXOs: 0, AbsolutFee: 0, Err: src.ErrInsufficientFunds},
	},
	{
		Comment: "sp recipient",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  10_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 8460, NumOfSelectedUTXOs: 1, AbsolutFee: 1540, Err: nil},
	},
	{
		Comment: "two sp recipient",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 40_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  10_000,
				},
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  10_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 1000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 18030, NumOfSelectedUTXOs: 1, AbsolutFee: 1970, Err: nil},
	},
	{
		Comment: "two p2pkh recipients",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 40_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "16SuRwey62GASRV2V7rtYoakDUd5oAApv7",
					Amount:  10_000,
				},
				{
					Address: "16SuRwey62GASRV2V7rtYoakDUd5oAApv7",
					Amount:  10_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 1000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 18210, NumOfSelectedUTXOs: 1, AbsolutFee: 1790, Err: nil},
	},
	{
		Comment: "two p2pkh recipients one pre set taproot output pkScript",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 40_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "16SuRwey62GASRV2V7rtYoakDUd5oAApv7",
					Amount:  10_000,
				},
				{
					Address: "16SuRwey62GASRV2V7rtYoakDUd5oAApv7",
					Amount:  10_000,
				},
				{
					PkScript: bytes.Repeat([]byte{0x00}, 34),
					Amount:   10_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 1000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 7780, NumOfSelectedUTXOs: 1, AbsolutFee: 2220, Err: nil},
	},
	{
		Comment: "two p2pkh and one sp recipient",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 40_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "16SuRwey62GASRV2V7rtYoakDUd5oAApv7",
					Amount:  10_000,
				},
				{
					Address: "16SuRwey62GASRV2V7rtYoakDUd5oAApv7",
					Amount:  10_000,
				},
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  10_000,
				},
			},
			FeeRate:         1,
			MinChangeAmount: 1000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 9778, NumOfSelectedUTXOs: 1, AbsolutFee: 222, Err: nil},
	},
	{
		Comment: "should fail - fee rate too low",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  10_000,
				},
			},
			FeeRate:         0,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 0, NumOfSelectedUTXOs: 0, AbsolutFee: 0, Err: src.ErrInvalidFeeRate},
	},
	{
		Comment: "should give insufficient funds error can't reach min change amount although target value was allocated",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  15_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 0, NumOfSelectedUTXOs: 0, AbsolutFee: 0, Err: src.ErrInsufficientFunds},
	},
	{
		Comment: "should fail - one recipient has zero amount",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  15_000,
				},
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  15_000,
				},
				{
					Address: "tsp1qqfqnnv8czppwysafq3uwgwvsc638hc8rx3hscuddh0xa2yd746s7xq36vuz08htp29hyml4u9shtlvcvqxuhjzldxjwyfnxmamz3ft8mh5tzx0hu",
					Amount:  0,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 0, NumOfSelectedUTXOs: 0, AbsolutFee: 0, Err: src.ErrRecipientAmountIsZero},
	},
	{
		Comment: "two p2pkh recipients - one invalid",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint32
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 40_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "16SuRwey62GASRV2V7rtYoakDUd5oAApv7",
					Amount:  10_000,
				},
				{
					Address: "16SuRw00000000000000YoakDUd5oAApv7",
					Amount:  10_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 1000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			AbsolutFee         uint64
			Err                error
		}{Change: 0, NumOfSelectedUTXOs: 0, AbsolutFee: 0, Err: errors.New("decoded address is of unknown format")},
	},
}

func TestFeeRateCoinSelector_CoinSelect(t *testing.T) {
	src.ChainParams = &chaincfg.MainNetParams

	for i, testCase := range testCases {
		t.Logf("Test Case: %d - %s", i, testCase.Comment)
		cs := NewFeeRateCoinSelector(testCase.Given.Utxos, testCase.Given.MinChangeAmount, testCase.Given.Recipients)
		selectedCoins, change, err := cs.CoinSelect(testCase.Given.FeeRate)
		if testCase.Expected.Err != nil {
			if err == nil {
				t.Errorf("expected error but got nil")
				return
			}
			if err.Error() != testCase.Expected.Err.Error() {
				t.Errorf("expected error %v but got %v", testCase.Expected.Err, err)
				return
			}
		}

		if change != testCase.Expected.Change {
			t.Errorf("Error: change incorrect %d != %d", change, testCase.Expected.Change)
			return
		}

		if len(selectedCoins) != testCase.Expected.NumOfSelectedUTXOs {
			t.Errorf("Error: wrong number of coins selected %d != %d", len(selectedCoins), testCase.Expected.NumOfSelectedUTXOs)
			return
		}

		if testCase.Expected.Err != nil {
			continue
		}

		// validate amounts
		var sumSelectedAmounts uint64
		for _, coin := range selectedCoins {
			sumSelectedAmounts += coin.Amount
		}

		var sumRecipientsAmounts int64
		for _, recipient := range testCase.Given.Recipients {
			sumRecipientsAmounts += recipient.Amount
		}

		var sumTxOutputs uint64 = uint64(sumRecipientsAmounts) + change
		var fee = sumSelectedAmounts - sumTxOutputs
		if sumSelectedAmounts-sumTxOutputs != testCase.Expected.AbsolutFee {
			t.Errorf("Error amounts: inputs - %d; outputs - %d; fee - %d --- expected fee %d", sumSelectedAmounts, sumTxOutputs, fee, testCase.Expected.AbsolutFee)
			return
		}

	}
}
