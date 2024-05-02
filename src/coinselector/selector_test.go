package coinselector

import (
	"bytes"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/gobip352"
	"testing"
)
import (
	"github.com/setavenger/blindbitd/src"
)

var (
	txid = gobip352.ConvertToFixedLength32(bytes.Repeat([]byte{0x00}, 32))
)

type TestCase struct {
	Comment string
	Given   struct {
		Utxos           src.UtxoCollection
		Recipients      []*src.Recipient
		FeeRate         uint64
		MinChangeAmount uint64
	}
	Expected struct {
		Change             uint64
		NumOfSelectedUTXOs int
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
			FeeRate         uint64
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
					Address: "bcrt1qs7k3knzplrzegwja45suv30584t4ryaf5xwf7f",
					Amount:  5_000,
				},
			},
			FeeRate:         1,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			Err                error
		}{Change: 14901, NumOfSelectedUTXOs: 1, Err: nil},
	},
	{
		Comment: "fits with change amount",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint64
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
					Address: "bcrt1qs7k3knzplrzegwja45suv30584t4ryaf5xwf7f",
					Amount:  54_000,
				},
			},
			FeeRate:         1,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			Err                error
		}{Change: 5901, NumOfSelectedUTXOs: 2, Err: nil},
	},
	{
		Comment: "Does it work with different fee rates",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint64
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
					Address: "bcrt1qs7k3knzplrzegwja45suv30584t4ryaf5xwf7f",
					Amount:  50_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			Err                error
		}{Change: 9010, NumOfSelectedUTXOs: 2, Err: nil},
	},
	{
		Comment: "fails because not enough funds",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint64
			MinChangeAmount uint64
		}{
			Utxos: src.UtxoCollection{
				{
					Amount: 20_000,
				},
			},
			Recipients: []*src.Recipient{
				{
					Address: "bcrt1qs7k3knzplrzegwja45suv30584t4ryaf5xwf7f",
					Amount:  20_000,
				},
			},
			FeeRate:         10,
			MinChangeAmount: 5000,
		},
		Expected: struct {
			Change             uint64
			NumOfSelectedUTXOs int
			Err                error
		}{Change: 0, NumOfSelectedUTXOs: 0, Err: src.ErrInsufficientFunds},
	},
	{
		Comment: "sp recipient",
		Given: struct {
			Utxos           src.UtxoCollection
			Recipients      []*src.Recipient
			FeeRate         uint64
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
			Err                error
		}{Change: 8890, NumOfSelectedUTXOs: 1, Err: nil},
	},
}

func TestFeeRateCoinSelector_CoinSelect(t *testing.T) {
	src.ChainParams = &chaincfg.RegressionNetParams

	for i, testCase := range testCases {
		t.Logf("Test Case: %d - %s", i, testCase.Comment)
		cs := NewFeeRateCoinSelector(testCase.Given.Utxos, testCase.Given.MinChangeAmount, testCase.Given.Recipients)
		selectedCoins, change, err := cs.CoinSelect(1)
		if !errors.Is(err, testCase.Expected.Err) {
			t.Errorf("Error: %s", err)
			return
		}

		if change != testCase.Expected.Change {
			t.Errorf("Error: change incorrect %d != %d", change, testCase.Expected.Change)
			return
		}

		if len(selectedCoins) != testCase.Expected.NumOfSelectedUTXOs {
			t.Errorf("Error: wrong number of coins selected %d != %d", len(selectedCoins), testCase.Expected.NumOfSelectedUTXOs)
			return
		}
	}
}
