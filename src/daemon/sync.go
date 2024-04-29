package daemon

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil/gcs"
	"github.com/btcsuite/btcd/btcutil/gcs/builder"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/database"
	"github.com/setavenger/blindbitd/src/networking"
	"github.com/setavenger/gobip352"
	"time"
)

// syncBlock there are several possibilities how this returns no error and still an empty slice for FoundOutputs
func (d *Daemon) syncBlock(blockHeight uint64) ([]src.OwnedUTXO, error) {

	tweaks, err := d.Client.GetTweaks(blockHeight, d.Wallet.DustLimit)
	if err != nil {
		return nil, err
	}

	// otherwise change will not be found
	labelsToCheck := append([]gobip352.Label{d.ChangeLabel}, d.Labels...)

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
			return nil, err
		}

		var outputPubKey [32]byte
		outputPubKey, err = gobip352.CreateOutputPubKey(sharedSecret, d.Wallet.PubKeySpend, 0)
		if err != nil {
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
				panic(err)
			}
			potentialOutputs = append(potentialOutputs, append([]byte{0x51, 0x20}, labelPotentialOutputPrep[1:]...))

			// add label with uneven parity as well
			var negatedLabelPubKey [33]byte
			negatedLabelPubKey, err = gobip352.NegatePublicKey(label.PubKey)
			if err != nil {
				panic(err)
			}
			labelPotentialOutputPrep, err = gobip352.AddPublicKeys(outputPubKey33, negatedLabelPubKey)
			if err != nil {
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

	filterData, err := d.Client.GetFilter(blockHeight)
	if err != nil {
		return nil, err
	}

	c := chainhash.Hash{}

	err = c.SetBytes(gobip352.ReverseBytesCopy(filterData.BlockHash))
	if err != nil {
		return nil, err

	}

	filter, err := gcs.FromNBytes(builder.DefaultP, builder.DefaultM, filterData.Data)
	if err != nil {
		return nil, err
	}

	key := builder.DeriveKey(&c)

	isMatch, err := filter.HashMatchAny(key, potentialOutputs)
	if err != nil {
		return nil, err
	}

	if !isMatch {
		return nil, nil
	}

	utxos, err := d.Client.GetUTXOs(blockHeight)
	if err != nil {
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
			return nil, err
		}
		foundOutputs = append(foundOutputs, foundOutputsPerTweak...)
	}

	// use a map to not have to iterate for every found UTXOServed, map should be faster lookup
	matchUTXOMap := make(map[[32]byte]*networking.UTXOServed)
	for _, utxo := range utxos {
		matchUTXOMap[gobip352.ConvertToFixedLength32(utxo.ScriptPubKey[2:])] = utxo
	}

	var ownedUTXOs []src.OwnedUTXO
	for _, foundOutput := range foundOutputs {

		utxo, exists := matchUTXOMap[foundOutput.Output]
		if !exists {
			return nil, src.ErrNoMatchForUTXO
		}
		state := src.StateUnspent
		if utxo.Spent {
			state = src.StateSpent
		}
		ownedUTXOs = append(ownedUTXOs, src.OwnedUTXO{
			Txid:         utxo.Txid,
			Vout:         utxo.Vout,
			Amount:       utxo.Amount,
			PrivKeyTweak: foundOutput.SecKeyTweak,
			PubKey:       foundOutput.Output,
			Timestamp:    utxo.Timestamp,
			State:        state,             // should normally always be unspent here
			Label:        foundOutput.Label, // todo add m once gobip352 is updated
		})
	}

	return ownedUTXOs, err
}

func (d *Daemon) SyncToTip(chainTip uint64) error {
	var err error
	if chainTip == 0 {
		chainTip, err = d.Client.GetChainTip()
		if err != nil {
			return err
		}
	}

	fmt.Println("Tip:", chainTip)
	var startHeight = d.Wallet.BirthHeight
	if d.Wallet.LastScanHeight > startHeight {
		startHeight = d.Wallet.LastScanHeight
	}

	if startHeight >= chainTip {
		return nil
	}

	for i := startHeight; i < chainTip+1; i++ {
		// possible logging here to indicate to the user
		fmt.Println("syncing:", i)
		var ownedUTXOs []src.OwnedUTXO
		ownedUTXOs, err = d.syncBlock(i)
		if err != nil {
			return err
		}
		if ownedUTXOs == nil {
			continue
		}
		d.Wallet.UTXOs = append(d.Wallet.UTXOs, ownedUTXOs...)
		d.Wallet.LastScanHeight = i
		err = database.WriteToDB(src.PathDbWallet, d.Wallet, d.Password)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Daemon) ContinuousScan() error {

	for {
		<-time.NewTicker(5 * time.Second).C
		chainTip, err := d.Client.GetChainTip()
		if err != nil {
			return err
		}
		if chainTip <= d.LastScanHeight {
			continue
		}

		err = d.SyncToTip(chainTip)
		if err != nil {
			return err
		}
	}
}
