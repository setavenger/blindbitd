package src

import "github.com/btcsuite/btcd/chaincfg"

var (

	/* [Network] */

	// BlindBitServerAddress Indexing server for silent payments that follows the blindbit standard
	BlindBitServerAddress string
	// ElectrumServerAddress Electrum server
	ElectrumServerAddress string

	/* [Wallet] */

	// MinChangeAmount The wallet will never create change that is smaller than this value. Value has to be in sats.
	MinChangeAmount int64
	// DustLimit only receives amounts and checks tweaks where the biggest utxo exceeds the dust limit.
	// Note: that if you receive funds below this threshold you might not find them.
	// Rescan with DustLimit = 0 to find those.
	DustLimit uint64

	// ChainParams defines on which chain the wallet runs
	ChainParams *chaincfg.Params
)
