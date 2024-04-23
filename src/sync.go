package src

import (
	"errors"
	"github.com/btcsuite/btcd/btcutil/gcs/builder"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/gcs"
	"github.com/setavenger/gobip352"
)

// syncBlock there are several possibilities how this returns no error and still an empty slice for FoundOutputs
func (d *Daemon) syncBlock(blockHeight uint64) ([]OwnedUTXO, error) {

	tweaks, err := d.Client.GetTweaks(blockHeight, d.Wallet.DustLimit)
	if err != nil {
		return nil, err
	}

	var potentialOutputs = make([][]byte, len(tweaks)+len(d.Wallet.PubKeysToWatch))

	for i, tweak := range tweaks {
		sharedSecret, err := gobip352.CreateSharedSecret(tweak, d.Wallet.secretKeyScan, nil)
		if err != nil {
			return nil, err
		}
		outputPubKey, err := gobip352.CreateOutputPubKey(sharedSecret, d.Wallet.PubKeySpend, 0)
		if err != nil {
			return nil, err
		}

		// todo we do this for now until the filters are changed to the 32byte x-only taproot pub keys
		potentialOutputs[i] = append([]byte{0x51, 0x20}, outputPubKey[:]...)
	}

	for i, scriptsToWatch := range d.Wallet.PubKeysToWatch {
		potentialOutputs[len(tweaks)+i] = append([]byte{0x51, 0x20}, scriptsToWatch[:]...)
	}

	if len(potentialOutputs) == 0 {
		return nil, nil
	}

	filterData, err := d.Client.GetFilter(blockHeight)
	if err != nil {
		return nil, err
	}

	c := chainhash.Hash{}

	err = c.SetBytes(gobip352.ReverseBytes(filterData.BlockHash))
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
		foundOutputsPerTweak, err = gobip352.ReceiverScanTransaction(d.secretKeyScan, d.PubKeySpend, d.Labels, blockOutputs, tweak, nil)
		if err != nil {
			return nil, err
		}
		foundOutputs = append(foundOutputs, foundOutputsPerTweak...)
	}

	// use a map to not have to iterate for every found UTXO map is faster lookup
	matchUTXOMap := make(map[[32]byte]*UTXO)
	for _, utxo := range utxos {
		matchUTXOMap[gobip352.ConvertToFixedLength32(utxo.ScriptPubKey[2:])] = utxo
	}

	var ownedUTXOs []OwnedUTXO
	for _, foundOutput := range foundOutputs {
		utxo, exists := matchUTXOMap[foundOutput.Output]
		if !exists {
			return nil, errors.New("could not match UTXO to foundOutput, should not happen")
		}
		ownedUTXOs = append(ownedUTXOs, OwnedUTXO{
			Txid:               utxo.Txid,
			Vout:               utxo.Vout,
			Amount:             utxo.Value,
			PrivKeyTweak:       foundOutput.SecKeyTweak,
			PubKey:             foundOutput.Output,
			TimestampConfirmed: 0,
			State:              StateUnspent, // should normally always be unspent here
			Label:              0,            // todo add m once gobip352 is updated
		})
	}

	return ownedUTXOs, err
}
