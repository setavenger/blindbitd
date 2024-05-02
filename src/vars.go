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
	MinChangeAmount int64 // MinChangeAmount

	// ChainParams defines on which chain the wallet runs
	ChainParams *chaincfg.Params
)
