package daemon

import (
	"context"
	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/btcutil/gcs/builder"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/database"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/networking"
	"github.com/setavenger/blindbitd/src/utils"
	"github.com/setavenger/gobip352"
	"time"
)

// syncBlock there are several possibilities how this returns no error and still an empty slice for FoundOutputs
func (d *Daemon) syncBlock(blockHeight uint64) ([]*src.OwnedUTXO, error) {

	tweaks, err := d.ClientBlindBit.GetTweaks(blockHeight, src.DustLimit)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	// otherwise change will not be found
	labelsToCheck := append([]*gobip352.Label{d.Wallet.ChangeLabel}, d.Wallet.Labels...)

	// todo change back to assigning via index slice[i] once we are sure how long a slice will be; can we be sure how long it will always be?
	var potentialOutputs [][]byte
	// check for all tweaks normal outputs
	// + all tweaks * labels
	// + all tweaks * labels with opposite parity
	// + all prior found scripts in case someone resent to one of those scripts
	//var potentialOutputs = make([][]byte, len(tweaks)*(len(labelsToCheck)*2)+len(d.Wallet.PubKeysToWatch))

	for _, tweak := range tweaks {
		var sharedSecret [33]byte
		sharedSecret, err = gobip352.CreateSharedSecret(tweak, d.Wallet.SecretKeyScan(), nil)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return nil, err
		}

		var outputPubKey [32]byte
		outputPubKey, err = gobip352.CreateOutputPubKey(sharedSecret, d.Wallet.PubKeySpend, 0)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return nil, err
		}

		// todo we do this for now until the filters are changed to the 32byte x-only taproot pub keys
		potentialOutputs = append(potentialOutputs, append([]byte{0x51, 0x20}, outputPubKey[:]...))
		for _, label := range labelsToCheck {

			outputPubKey33 := gobip352.ConvertToFixedLength33(append([]byte{0x02}, outputPubKey[:]...))

			// even parity
			var labelPotentialOutputPrep [33]byte
			labelPotentialOutputPrep, err = gobip352.AddPublicKeys(outputPubKey33, label.PubKey)
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}
			potentialOutputs = append(potentialOutputs, append([]byte{0x51, 0x20}, labelPotentialOutputPrep[1:]...))

			// add label with uneven parity as well
			var negatedLabelPubKey [33]byte
			negatedLabelPubKey, err = gobip352.NegatePublicKey(label.PubKey)
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}
			labelPotentialOutputPrep, err = gobip352.AddPublicKeys(outputPubKey33, negatedLabelPubKey)
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}
			potentialOutputs = append(potentialOutputs, append([]byte{0x51, 0x20}, labelPotentialOutputPrep[1:]...))
		}
	}

	// todo change back to assigning via index slice[i] once we are sure how long a slice will be; can we be sure how long it will always be?
	for _, scriptsToWatch := range d.Wallet.PubKeysToWatch {
		potentialOutputs = append(potentialOutputs, append([]byte{0x51, 0x20}, scriptsToWatch[:]...))
	}

	if len(potentialOutputs) == 0 {
		return nil, nil
	}

	filterData, err := d.ClientBlindBit.GetFilter(blockHeight)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	c := chainhash.Hash{}

	err = c.SetBytes(gobip352.ReverseBytesCopy(filterData.BlockHash))
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err

	}

	filter, err := gcs.FromNBytes(builder.DefaultP, builder.DefaultM, filterData.Data)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	key := builder.DeriveKey(&c)

	isMatch, err := filter.HashMatchAny(key, potentialOutputs)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	if !isMatch {
		return nil, nil
	}

	utxos, err := d.ClientBlindBit.GetUTXOs(blockHeight)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	var foundOutputs []*gobip352.FoundOutput

	var blockOutputs = make([][32]byte, len(utxos)) // we use it as txOutputs we check against all outputs from the block
	for i, utxo := range utxos {
		blockOutputs[i] = gobip352.ConvertToFixedLength32(utxo.ScriptPubKey[2:])
	}

	for _, tweak := range tweaks {
		var foundOutputsPerTweak []*gobip352.FoundOutput
		foundOutputsPerTweak, err = gobip352.ReceiverScanTransaction(d.Wallet.SecretKeyScan(), d.Wallet.PubKeySpend, labelsToCheck, blockOutputs, tweak, nil)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return nil, err
		}
		foundOutputs = append(foundOutputs, foundOutputsPerTweak...)
	}

	// use a map to not have to iterate for every found UTXOServed, map should be faster lookup
	matchUTXOMap := make(map[[32]byte]*networking.UTXOServed)
	for _, utxo := range utxos {
		matchUTXOMap[gobip352.ConvertToFixedLength32(utxo.ScriptPubKey[2:])] = utxo
	}

	var ownedUTXOs []*src.OwnedUTXO
	for _, foundOutput := range foundOutputs {

		utxo, exists := matchUTXOMap[foundOutput.Output]
		if !exists {
			err = src.ErrNoMatchForUTXO
			logging.ErrorLogger.Println(err)
			return nil, err
		}
		state := src.StateUnspent
		if utxo.Spent {
			state = src.StateSpent
		}
		ownedUTXOs = append(ownedUTXOs, &src.OwnedUTXO{
			Txid:         utxo.Txid,
			Vout:         utxo.Vout,
			Amount:       utxo.Amount,
			PrivKeyTweak: foundOutput.SecKeyTweak,
			PubKey:       foundOutput.Output,
			Timestamp:    utxo.Timestamp,
			State:        state,
			Label:        foundOutput.Label,
		})
	}

	return ownedUTXOs, err
}

func (d *Daemon) SyncToTip(chainTip uint64) error {
	var err error
	if chainTip == 0 {
		chainTip, err = d.ClientBlindBit.GetChainTip()
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
	}

	logging.DebugLogger.Println("Tip:", chainTip)
	// todo find fixed points for mainnet/signet/testnet where startHeight can start from. Avoid scanning through non SP merged blocks
	var startHeight = d.Wallet.BirthHeight
	if d.Wallet.LastScanHeight > startHeight {
		startHeight = d.Wallet.LastScanHeight + 1
	}

	if startHeight >= chainTip {
		return nil
	}

	// don't check genesis block
	if startHeight == 0 {
		startHeight = 1
	}

	for i := startHeight; i < chainTip+1; i++ {
		// possible logging here to indicate to the user
		logging.DebugLogger.Println("syncing:", i)
		var ownedUTXOs []*src.OwnedUTXO
		ownedUTXOs, err = d.syncBlock(i)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
		if ownedUTXOs == nil {
			d.Wallet.LastScanHeight = i
			continue
		}
		err = d.Wallet.AddUTXOs(ownedUTXOs)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
		d.Wallet.LastScanHeight = i
		err = database.WriteToDB(src.PathDbWallet, d.Wallet, d.Password)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
	}

	err = d.CheckUnspentUTXOs()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func (d *Daemon) ForceSyncFrom(fromHeight uint64) error {
	chainTip, err := d.ClientBlindBit.GetChainTip()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	logging.InfoLogger.Printf("ForceSyncFrom: %d to %d\n", fromHeight, chainTip)

	// don't check genesis block
	if fromHeight == 0 {
		fromHeight = 1
	}

	for i := fromHeight; i < chainTip+1; i++ {
		// possible logging here to indicate to the user
		logging.DebugLogger.Println("syncing:", i)
		var ownedUTXOs []*src.OwnedUTXO
		ownedUTXOs, err = d.syncBlock(i)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
		if ownedUTXOs == nil {
			d.Wallet.LastScanHeight = i
			continue
		}
		err = d.Wallet.AddUTXOs(ownedUTXOs)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
		d.Wallet.LastScanHeight = i
		err = database.WriteToDB(src.PathDbWallet, d.Wallet, d.Password)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
	}

	err = d.CheckUnspentUTXOs()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func (d *Daemon) ContinuousScan() error {
	d.Status = pb.Status_STATUS_SCANNING
	for {
		select {
		case newBlock := <-d.NewBlockChan:
			<-time.After(5 * time.Second) // delay, indexing server does not index immediately after a block is found

			err := d.SyncToTip(uint64(newBlock.Height))
			if err != nil {
				logging.ErrorLogger.Println(err)
				return err
			}
			logging.InfoLogger.Printf("New balance: %d\n", d.Wallet.FreeBalance())
		case <-time.NewTicker(5 * time.Minute).C:
			// todo is this needed if NewBlockChan is very robust?
			// check every 5 minutes anyway
			chainTip, err := d.ClientBlindBit.GetChainTip()
			if err != nil {
				logging.ErrorLogger.Println(err)
				return err
			}

			if chainTip <= d.Wallet.LastScanHeight {
				continue
			}

			err = d.SyncToTip(chainTip)
			if err != nil {
				logging.ErrorLogger.Println(err)
				return err
			}
		case <-time.NewTicker(1 * time.Minute).C:
			// exclusively to check for spent UTXOs
			err := d.CheckUnspentUTXOs()
			if err != nil {
				logging.ErrorLogger.Println(err)
				return err
			}
		}
	}
}

// CheckUnspentUTXOs
// checks against electrum whether unspent owned UTXOs are now unspent
func (d *Daemon) CheckUnspentUTXOs() error {
	// todo this probably breaks if more than one UTXO are locked to a script
	//  this should never happen if the protocol is followed but still might occur
	for _, utxo := range d.Wallet.GetUTXOsByStates(src.StateUnspent, src.StateUnconfirmedSpent) {
		balance, err := d.ClientElectrum.GetBalance(context.Background(), utils.ConvertPubKeyToScriptHash(utxo.PubKey))
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
		if balance.Confirmed == 0.0 && balance.Unconfirmed == 0.0 {
			utxo.State = src.StateSpent
			continue
		}
		if balance.Unconfirmed < 0 {
			utxo.State = src.StateUnconfirmedSpent
			continue
		}
	}
	return nil
}
