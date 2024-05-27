package daemon

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/btcutil/txsort"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/coinselector"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/utils"
	"github.com/setavenger/go-bip352"
)

const ExtraDataAmountKey = "amount"

// SendToRecipients
// creates a signed transaction that sends to the specified recipients
// todo should all these functions just be Daemon functions
// use markSpent to set the used UTXOs to spent_unconfirmed
// use useSpentUnconfirmed to also include spent_undconfirmed UTXOs in the coinSelection process
func (d *Daemon) SendToRecipients(recipients []*src.Recipient, feeRate int64, markSpent, useSpentUnconfirmed bool) ([]byte, error) {

	selector := coinselector.NewFeeRateCoinSelector(d.Wallet.GetFreeUTXOs(useSpentUnconfirmed), uint64(src.MinChangeAmount), recipients)

	selectedUTXOs, changeAmount, err := selector.CoinSelect(uint32(feeRate))
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	// vins is the final selection of coins, which can then be used to derive silentPayment Outputs
	var vins = make([]*bip352.Vin, len(selectedUTXOs))
	for i, utxo := range selectedUTXOs {
		vin := src.ConvertOwnedUTXOIntoVin(utxo)
		fullVinSecretKey := bip352.AddPrivateKeys(*vin.SecretKey, d.Wallet.SecretKeySpend())
		vin.SecretKey = &fullVinSecretKey
		vins[i] = &vin
	}

	// now we need the difference between the inputs and outputs so that we can assign a value for change
	var sumAllInputs int64
	for _, vin := range vins {
		sumAllInputs += int64(vin.Amount)
	}

	if changeAmount > 0 {
		// change exists, and it should be greater than the MinChangeAmount
		recipients = append(recipients, &src.Recipient{
			Address: d.Wallet.ChangeLabel.Address,
			Amount:  int64(changeAmount),
		})
	}

	// extract the ScriptPubKeys of the SP recipients with the selected txInputs
	recipients, err = ParseRecipients(recipients, vins, src.ChainParams)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	err = sanityCheckRecipientsForSending(recipients)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	packet, err := CreateUnsignedPsbt(recipients, vins)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	err = SignPsbt(packet, vins)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	err = psbt.MaybeFinalizeAll(packet)
	if err != nil {
		logging.ErrorLogger.Println(err)
		panic(err) // todo remove panic
	}

	finalTx, err := psbt.Extract(packet)
	if err != nil {
		logging.ErrorLogger.Println(err)
		panic(err) // todo remove panic
	}

	var sumAllOutputs int64
	for _, recipient := range recipients {
		sumAllOutputs += recipient.Amount
	}
	vSize := mempool.GetTxVirtualSize(btcutil.NewTx(finalTx))
	actualFee := sumAllInputs - sumAllOutputs
	actualFeeRate := float64(actualFee) / float64(vSize)

	errorTerm := 0.25 // todo make variable
	if actualFeeRate > float64(feeRate)+errorTerm {
		err = fmt.Errorf("actual fee rate deviates to strong from desired fee rate: %f > %d", actualFeeRate, feeRate)
		return nil, err
	}

	if actualFeeRate < float64(feeRate)-errorTerm {
		err = fmt.Errorf("actual fee rate deviates to strong from desired fee rate: %f < %d", actualFeeRate, feeRate)
		return nil, err
	}

	var buf bytes.Buffer
	err = finalTx.Serialize(&buf)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	if markSpent {
		var found int
		// now that everything worked mark as spent if desired
		for _, vin := range vins {
			vinOutpoint, err := utils.SerialiseVinToOutpoint(*vin)
			if err != nil {
				logging.ErrorLogger.Println(err)
				return nil, err
			}
			for _, utxo := range d.Wallet.UTXOs {
				utxoOutpoint, err := utxo.SerialiseToOutpoint()
				if err != nil {
					logging.ErrorLogger.Println(err)
					return nil, err
				}
				if bytes.Equal(vinOutpoint[:], utxoOutpoint[:]) {
					utxo.State = src.StateUnconfirmedSpent
					found++
					logging.DebugLogger.Printf("Marked %x as spent\n", utxoOutpoint)
				}
			}
		}
		if found != len(vins) {
			err = fmt.Errorf("we could not mark enough utxos as spent. marked %d, needed %d", found, len(vins))
			return nil, err
		}
	}

	return buf.Bytes(), err
}

// ParseRecipients
// Checks all recipients and adds the PkScript based on the given address.
// Silent Payment addresses are also parsed and the outputs will be computed based on the vins.
// For that reason this function has to be called after the final coinSelection is done.
// Otherwise, the SP outputs will NOT be found by the receiver.
// SP Recipients are always at the end.
// Hence, the tx must be sorted according to BIP 69 to avoid a specific signature of this wallet.
//
// NOTE: Existing PkScripts will NOT be overridden, those recipients will be skipped and returned as given
// todo keep original order in case that is relevant for any use case?
func ParseRecipients(recipients []*src.Recipient, vins []*bip352.Vin, chainParam *chaincfg.Params) ([]*src.Recipient, error) {
	var spRecipients []*bip352.Recipient

	// newRecipients tracks the modified group of recipients in order to avoid clashes
	var newRecipients []*src.Recipient
	for _, recipient := range recipients {
		if recipient.PkScript != nil {
			// skip if a pkScript is already present (for what ever reason)
			newRecipients = append(newRecipients, recipient)
			continue
		}
		isSP := utils.IsSilentPaymentAddress(recipient.Address)
		if !isSP {
			address, err := btcutil.DecodeAddress(recipient.Address, chainParam)
			if err != nil {
				logging.ErrorLogger.Printf("Failed to decode address: %v", err)
				return nil, err
			}
			scriptPubKey, err := txscript.PayToAddrScript(address)
			if err != nil {
				logging.ErrorLogger.Printf("Failed to create scriptPubKey: %v", err)
				return nil, err
			}
			recipient.PkScript = scriptPubKey

			newRecipients = append(newRecipients, recipient)
			continue
		}

		spRecipients = append(spRecipients, &bip352.Recipient{
			SilentPaymentAddress: recipient.Address,
			Amount:               uint64(recipient.Amount),
		})
	}

	var mainnet bool
	if src.ChainParams.Name == chaincfg.MainNetParams.Name {
		mainnet = true
	}

	if len(spRecipients) > 0 {
		err := bip352.SenderCreateOutputs(spRecipients, vins, mainnet, false)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return nil, err
		}
	}

	for _, spRecipient := range spRecipients {
		newRecipients = append(newRecipients, ConvertSPRecipient(spRecipient))
	}

	// This case might not be realistic so the check could potentially be removed safely
	if len(recipients) != len(newRecipients) {
		// for some reason we have a different number of recipients after parsing them.
		return nil, src.ErrWrongLengthRecipients
	}

	return newRecipients, nil
}

// sanityCheckRecipientsForSending
// checks whether any of the Recipients lacks the necessary information to construct the transaction.
// required for every recipient: Recipient.PkScript and Recipient.Amount
func sanityCheckRecipientsForSending(recipients []*src.Recipient) error {
	for _, recipient := range recipients {
		if recipient.PkScript == nil || recipient.Amount == 0 {
			// if we choose a lot of logging in this module/program we could log the incomplete recipient here
			return src.ErrRecipientIncomplete
		}
	}
	return nil
}

func CreateUnsignedPsbt(recipients []*src.Recipient, vins []*bip352.Vin) (*psbt.Packet, error) {
	var txOutputs []*wire.TxOut
	for _, recipient := range recipients {
		txOutputs = append(txOutputs, wire.NewTxOut(recipient.Amount, recipient.PkScript))
	}

	var txInputs []*wire.TxIn
	for _, vin := range vins {
		hash, err := chainhash.NewHash(bip352.ReverseBytesCopy(vin.Txid[:]))
		if err != nil {
			return nil, err
		}
		prevOut := wire.NewOutPoint(hash, vin.Vout)
		txInputs = append(txInputs, wire.NewTxIn(prevOut, nil, nil))
	}

	unsignedTx := &wire.MsgTx{
		Version: 2,
		TxIn:    txInputs,
		TxOut:   txOutputs,
	}

	packet := &psbt.Packet{
		UnsignedTx: txsort.Sort(unsignedTx),
	}

	return packet, nil
}

// SignPsbt
// fails if inputs in packet have a different order than vins
func SignPsbt(packet *psbt.Packet, vins []*bip352.Vin) error {
	if len(packet.UnsignedTx.TxIn) != len(vins) {
		return src.ErrTxInputAndVinLengthMismatch
	}

	prevOutsForFetcher := make(map[wire.OutPoint]*wire.TxOut, len(vins))

	// simple map to find correct vin for prevOutsForFetcher
	vinMap := make(map[string]bip352.Vin, len(vins))
	for _, v := range vins {
		vinMap[fmt.Sprintf("%x:%d", v.Txid, v.Vout)] = *v
	}

	for i := 0; i < len(vins); i++ {
		outpoint := packet.UnsignedTx.TxIn[i].PreviousOutPoint
		key := fmt.Sprintf("%x:%d", bip352.ReverseBytesCopy(outpoint.Hash[:]), outpoint.Index)
		vin, ok := vinMap[key]
		if !ok {
			err := fmt.Errorf("a vin was not found in the map, should not happen. upstream error in psbt and vin selection and or construction")
			logging.ErrorLogger.Println(err)
			return err
		}
		prevOutsForFetcher[outpoint] = wire.NewTxOut(int64(vin.Amount), vin.ScriptPubKey)
	}

	multiFetcher := txscript.NewMultiPrevOutFetcher(prevOutsForFetcher)

	sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, multiFetcher)

	var pInputs []psbt.PInput

	for iOuter, input := range packet.UnsignedTx.TxIn {
		signatureHash, err := txscript.CalcTaprootSignatureHash(sigHashes, txscript.SigHashDefault, packet.UnsignedTx, iOuter, multiFetcher)
		if err != nil {
			logging.ErrorLogger.Println(err)
			panic(err)
		}

		pInput, err := matchAndSign(input, signatureHash, vins)
		if err != nil {
			logging.ErrorLogger.Println(err)
			panic(err)
		}

		pInputs = append(pInputs, pInput)

	}

	packet.Inputs = pInputs

	return nil

}

func matchAndSign(input *wire.TxIn, signatureHash []byte, vins []*bip352.Vin) (psbt.PInput, error) {
	var psbtInput psbt.PInput

	for _, vin := range vins {
		if bytes.Equal(input.PreviousOutPoint.Hash[:], bip352.ReverseBytesCopy(vin.Txid[:])) &&
			input.PreviousOutPoint.Index == vin.Vout {
			privKey, pk := btcec.PrivKeyFromBytes(vin.SecretKey[:])

			if pk.Y().Bit(0) == 1 {
				newBytes := privKey.Key.Negate().Bytes()
				privKey, _ = btcec.PrivKeyFromBytes(newBytes[:])
			}
			signature, err := schnorr.Sign(privKey, signatureHash)
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}

			var witnessBytes bytes.Buffer
			err = psbt.WriteTxWitness(&witnessBytes, [][]byte{signature.Serialize()})
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}

			return psbt.PInput{
				WitnessUtxo:        wire.NewTxOut(int64(vin.Amount), vin.ScriptPubKey),
				SighashType:        txscript.SigHashDefault,
				FinalScriptWitness: witnessBytes.Bytes(),
			}, err
		}
	}

	return psbtInput, src.ErrNoMatchingVinFoundForTxInput

}

/*  util functions */

// ConvertSPRecipient converts a bip352.Recipient to a Recipient native to this program
func ConvertSPRecipient(recipient *bip352.Recipient) *src.Recipient {
	return &src.Recipient{
		Address:  recipient.SilentPaymentAddress,
		PkScript: append([]byte{0x51, 0x20}, recipient.Output[:]...),
		Amount:   int64(recipient.Amount),
		Data:     recipient.Data,
	}
}

// BroadcastTx
// broadcasts a transaction and returns the txid
func (d *Daemon) BroadcastTx(rawTx []byte) (string, error) {
	if !src.UseElectrum {
		return "", errors.New("currently can't broadcast without electrum; either activate electrum or publish on another channel")
	}
	// todo parse tx and check against the outpoints which where spent and mark UTXOs locally.
	txid, err := d.ClientElectrum.BroadcastTransaction(context.Background(), hex.EncodeToString(rawTx))
	if err != nil {
		return "", err
	}
	return txid, nil
}
