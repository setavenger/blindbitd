package src

import (
	"bytes"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/coinset"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/btcutil/txsort"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/setavenger/blindbitd/src/utils"
	"github.com/setavenger/gobip352"
)

const ExtraDataAmountKey = "amount"

// SendToRecipients
// creates a signed transaction that sends to the specified recipients
func (w *Wallet) SendToRecipients(recipients []*Recipient, fee int64) ([]byte, error) {

	var sumAllOutputs int64
	for _, recipient := range recipients {
		sumAllOutputs += recipient.Amount
	}

	var allPossibleVins []gobip352.Vin
	for _, utxo := range w.UTXOs {
		vin := ConvertOwnedUTXOIntoVin(utxo)
		allPossibleVins = append(allPossibleVins, vin)
	}

	allPossibleCoins := make([]coinset.Coin, len(allPossibleVins))
	for i, vin := range allPossibleVins {
		vinCopy := vin
		allPossibleCoins[i] = &vinCopy
	}

	coins, err := coinset.MinNumberCoinSelector{
		MaxInputs:       len(w.UTXOs),
		MinChangeAmount: btcutil.Amount(MinChangeAmount),
	}.CoinSelect(btcutil.Amount(sumAllOutputs+fee), allPossibleCoins)
	if err != nil {
		// ErrCoinsNoSelectionAvailable
		return nil, err
	}

	// vins is the final selection of coins, which can then be used to derive silentPayment Outputs
	var vins = make([]*gobip352.Vin, len(coins.Coins()))
	for i, coin := range coins.Coins() {
		if vin, ok := coin.(*gobip352.Vin); ok {
			fullVinSecretKey := gobip352.AddPrivateKeys(*vin.SecretKey, w.secretKeySpend)
			vin.SecretKey = &fullVinSecretKey
			vins[i] = vin
		} else {
			panic("coin was not a vin")
		}
	}

	// todo we only do this and the fee calculation until we have a CoinSelector
	//  which incorporates a fee rate and automatically returns a change output

	// now we need the difference between the inputs and outputs so that we can assign a value for change
	var sumAllInputs int64
	for _, vin := range vins {
		sumAllInputs += int64(vin.Amount)
	}

	difference := sumAllInputs - sumAllOutputs

	switch changeAmount := difference - fee; {
	case changeAmount == 0:
	// do nothing, no change output needed
	case changeAmount < MinChangeAmount:
		// here we fail because the changeAmount is not enough
		return nil, ErrMinChangeAmountNotReached
	default:
		// change exists, and it is greater than the MinChangeAmount
		recipients = append(recipients, &Recipient{
			Address: w.ChangeLabel.Address,
			Amount:  changeAmount,
		})
	}

	// extract the ScriptPubKeys of the SP recipients with the selected txInputs
	recipients, err = ParseSPRecipients(recipients, vins)
	if err != nil {
		return nil, err
	}

	err = sanityCheckRecipientsForSending(recipients)
	if err != nil {
		return nil, err
	}

	packet, err := CreateUnsignedPsbt(recipients, vins)
	if err != nil {
		return nil, err
	}

	err = SignPsbt(packet, vins)
	if err != nil {
		return nil, err
	}

	err = psbt.MaybeFinalizeAll(packet)
	if err != nil {
		panic(err)
	}

	finalTx, err := psbt.Extract(packet)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = finalTx.Serialize(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

// ParseSPRecipients
// Checks all recipients whether we are sending to a Silent Payment address.
// If so we derive the corresponding output.
// Has to be called after the final coinSelection is done.
// Otherwise, the SP outputs will not be found by the receiver.
// SP Recipients are always at the end.
// Hence, the tx must be sorted according to BIP 69 to avoid a specific signature of this wallet.
// todo keep original order in case that is relevant for any use case?
func ParseSPRecipients(recipients []*Recipient, vins []*gobip352.Vin) ([]*Recipient, error) {
	var spRecipients []*gobip352.Recipient

	// newRecipients tracks the modified group of recipients in order to avoid clashes
	var newRecipients []*Recipient
	for _, recipient := range recipients {
		isSP := utils.IsSilentPaymentAddress(recipient.Address)
		if !isSP {
			newRecipients = append(newRecipients, recipient)
			continue
		} else {
			recipient.PkScript = nil
		}

		extraData := map[string]any{}
		extraData[ExtraDataAmountKey] = recipient.Amount

		spRecipients = append(spRecipients, &gobip352.Recipient{
			SilentPaymentAddress: recipient.Address,
			Amount:               uint64(recipient.Amount),
			Data:                 extraData,
		})
	}

	err := gobip352.SenderCreateOutputs(spRecipients, vins, false)
	if err != nil {
		return nil, err
	}

	for _, spRecipient := range spRecipients {
		newRecipients = append(newRecipients, ConvertSPRecipient(spRecipient))
	}

	// This case might not be realistic so the check could potentially be removed safely
	if len(recipients) != len(newRecipients) {
		// for some reason we have a different number of recipients after parsing them.
		return nil, ErrWrongLengthRecipients
	}

	return newRecipients, nil
}

// sanityCheckRecipientsForSending
// checks whether any of the Recipients lacks the necessary information to construct the transaction.
// required for every recipient: Recipient.PkScript and Recipient.Amount
func sanityCheckRecipientsForSending(recipients []*Recipient) error {
	for _, recipient := range recipients {
		if recipient.PkScript == nil || recipient.Amount == 0 {
			// if we choose a lot of logging in this module/program we could log the incomplete recipient here
			return ErrRecipientIncomplete
		}
	}
	return nil
}

func CreateUnsignedPsbt(recipients []*Recipient, vins []*gobip352.Vin) (*psbt.Packet, error) {
	var txOutputs []*wire.TxOut
	for _, recipient := range recipients {
		txOutputs = append(txOutputs, wire.NewTxOut(recipient.Amount, recipient.PkScript))
	}

	//var psbtInputs []psbt.PInput
	var txInputs []*wire.TxIn
	for _, vin := range vins {
		hash, err := chainhash.NewHash(gobip352.ReverseBytesCopy(vin.Txid[:]))
		if err != nil {
			return nil, err
		}
		prevOut := wire.NewOutPoint(hash, vin.Vout)
		txInputs = append(txInputs, wire.NewTxIn(prevOut, nil, nil))

		//psbtInput := psbt.PInput{
		//	WitnessUtxo: wire.NewTxOut(int64(vin.Amount), vin.ScriptPubKey),
		//	SighashType: txscript.SigHashDefault,
		//}
		//psbtInputs = append(psbtInputs, psbtInput)
	}

	unsignedTx := &wire.MsgTx{
		Version: 2,
		TxIn:    txInputs,
		TxOut:   txOutputs,
	}

	packet := &psbt.Packet{
		UnsignedTx: txsort.Sort(unsignedTx),
		//Inputs:     psbtInputs,
		//Outputs:    []psbt.POutput{},
	}

	return packet, nil
}

// SignPsbt
// fails if inputs in packet have a different order than vins
func SignPsbt(packet *psbt.Packet, vins []*gobip352.Vin) error {
	if len(packet.UnsignedTx.TxIn) != len(vins) {
		return ErrTxInputAndVinLengthMismatch
	}

	prevOutsForFetcher := make(map[wire.OutPoint]*wire.TxOut, len(vins))
	for i, vin := range vins {
		prevOutsForFetcher[packet.UnsignedTx.TxIn[i].PreviousOutPoint] = wire.NewTxOut(int64(vin.Amount), vin.ScriptPubKey)
	}

	multiFetcher := txscript.NewMultiPrevOutFetcher(prevOutsForFetcher)

	sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, multiFetcher)

	var pInputs []psbt.PInput

	for iOuter, input := range packet.UnsignedTx.TxIn {

		signatureHash, err := txscript.CalcTaprootSignatureHash(sigHashes, txscript.SigHashDefault, packet.UnsignedTx, iOuter, multiFetcher)
		if err != nil {
			panic(err)
		}

		pInput, err := matchAndSign(input, signatureHash, vins)
		if err != nil {
			panic(err)
		}

		pInputs = append(pInputs, pInput)

	}

	packet.Inputs = pInputs

	return nil

}

func matchAndSign(input *wire.TxIn, signatureHash []byte, vins []*gobip352.Vin) (psbt.PInput, error) {
	var psbtInput psbt.PInput

	for _, vin := range vins {
		if bytes.Equal(input.PreviousOutPoint.Hash[:], gobip352.ReverseBytesCopy(vin.Txid[:])) &&
			input.PreviousOutPoint.Index == vin.Vout {
			privKey, pk := btcec.PrivKeyFromBytes(vin.SecretKey[:])

			if pk.Y().Bit(0) == 1 {
				newBytes := privKey.Key.Negate().Bytes()
				privKey, _ = btcec.PrivKeyFromBytes(newBytes[:])
			}
			signature, err := schnorr.Sign(privKey, signatureHash)
			if err != nil {
				panic(err)
			}

			var witnessBytes bytes.Buffer
			err = psbt.WriteTxWitness(&witnessBytes, [][]byte{signature.Serialize()})
			if err != nil {
				panic(err)
			}

			return psbt.PInput{
				WitnessUtxo:        wire.NewTxOut(int64(vin.Amount), vin.ScriptPubKey),
				SighashType:        txscript.SigHashDefault,
				FinalScriptWitness: witnessBytes.Bytes(),
			}, err

		}
	}

	return psbtInput, ErrNoMatchingVinFoundForTxInput

}

/*  util functions */

// ConvertSPRecipient converts a gobip352.Recipient to a Recipient native to this program
func ConvertSPRecipient(recipient *gobip352.Recipient) *Recipient {
	return &Recipient{
		Address:    recipient.SilentPaymentAddress,
		PkScript:   append([]byte{0x51, 0x20}, recipient.Output[:]...),
		Amount:     int64(recipient.Amount),
		Annotation: recipient.Data,
	}
}
