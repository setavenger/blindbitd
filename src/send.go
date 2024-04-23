package src

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/setavenger/gobip352"
)

const ExtraDataAmountKey = "amount"

type Recipient struct {
	Address    string
	Amount     uint64
	Annotation string
}

func (w *Wallet) SendToSPAddress(recipients []Recipient) {

	var silentPaymentRecipients []*gobip352.Recipient
	for _, recipient := range recipients {
		isSP := IsSilentPaymentAddress(recipient.Address)
		if isSP {
			extraData := map[string]any{}
			extraData[ExtraDataAmountKey] = recipient.Amount

			silentPaymentRecipients = append(silentPaymentRecipients, &gobip352.Recipient{
				SilentPaymentAddress: recipient.Address,
				ScanPubKey:           nil,
				SpendPubKey:          nil,
				Output:               [32]byte{},
				Data:                 extraData,
			})
		}
	}

	//gobip352.SenderCreateOutputs(silentPaymentRecipients, )

}

func (w *Wallet) SimpleSendToRecipient(recipient Recipient) ([]byte, error) {
	isSP := IsSilentPaymentAddress(recipient.Address)
	if !isSP {
		return nil, errors.New("not a sp address")
	}

	extraData := map[string]any{}
	extraData[ExtraDataAmountKey] = recipient.Amount

	recipientModified := &gobip352.Recipient{
		SilentPaymentAddress: recipient.Address,
		Data:                 extraData,
	}

	if len(w.UTXOs) < 1 {
		return nil, errors.New("no utxos in wallet")
	}

	utxo := w.UTXOs[0]
	vin := w.ConvertOwnedUTXOIntoVin(utxo)

	/// creates the necessary outputs
	err := gobip352.SenderCreateOutputs([]*gobip352.Recipient{recipientModified}, []*gobip352.Vin{&vin}, false)
	if err != nil {
		return nil, err
	}

	fmt.Printf("cOut: %x\n", recipientModified.Output)
	fmt.Printf("utxo_pubKey: %x\n", utxo.PubKey)
	fmt.Printf("utxo_tweak: %x\n", utxo.PrivKeyTweak)

	fmt.Printf("utxo_txid: %x\n", utxo.Txid)
	fmt.Printf("utxo_vout: %d\n", utxo.Vout)
	fmt.Printf("utxo_value: %d\n", utxo.Amount)

	fullPrivKey := gobip352.AddPrivateKeys(w.secretKeySpend, utxo.PrivKeyTweak)

	//// Transaction creation and signing below
	sendValue := utxo.Amount - 5000 // 200 is fee

	out1 := wire.NewTxOut(int64(sendValue), append([]byte{0x51, 0x20}, recipientModified.Output[:]...))

	hash1, err := chainhash.NewHash(gobip352.ReverseBytes(utxo.Txid[:]))
	if err != nil {
		return nil, err
	}
	prevOut1 := wire.NewOutPoint(hash1, utxo.Vout)

	outgoingTx := &wire.MsgTx{
		Version: 2,
		TxIn:    []*wire.TxIn{wire.NewTxIn(prevOut1, nil, nil)},
		TxOut:   []*wire.TxOut{out1},
	}
	fmt.Printf("%+v\n", out1)

	//============
	//============

	// Create the packet that we want to sign.
	packet := &psbt.Packet{
		UnsignedTx: outgoingTx,
		Inputs: []psbt.PInput{{
			WitnessUtxo:        wire.NewTxOut(int64(utxo.Amount), append([]byte{0x51, 0x20}, utxo.PubKey[:]...)),
			SighashType:        txscript.SigHashDefault,
			TaprootInternalKey: fullPrivKey[:],
		}},
		Outputs: []psbt.POutput{{TaprootInternalKey: append([]byte{0x51, 0x20}, recipientModified.Output[:]...)}},
	}

	//============
	//============

	// generalise the prepending and don't just do it adhoc
	fetcher := txscript.NewCannedPrevOutputFetcher(
		append([]byte{0x51, 0x20}, utxo.PubKey[:]...), int64(utxo.Amount),
	)

	sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, fetcher)
	// For example, calculating the sighash for the first input
	fmt.Printf("%x\n", sigHashes.HashPrevOutsV1[:])
	fmt.Printf("%x\n", sigHashes.HashInputAmountsV1[:])
	fmt.Printf("%x\n", sigHashes.HashInputScriptsV1[:])
	fmt.Printf("%x\n", sigHashes.HashSequenceV1[:])
	fmt.Printf("%x\n", sigHashes.HashOutputsV1[:])
	privKey, pk := btcec.PrivKeyFromBytes(fullPrivKey[:]) // private key is correct for that input pubKey

	if pk.Y().Bit(0) == 1 {
		newBytes := privKey.Key.Negate().Bytes()
		privKey, _ = btcec.PrivKeyFromBytes(newBytes[:])
	}

	fmt.Printf("\nsecretKey: %x\n", privKey.Serialize())
	witnessScript, err := txscript.TaprootWitnessSignature(
		packet.UnsignedTx, sigHashes, 0, int64(utxo.Amount),
		append([]byte{0x51, 0x20}, utxo.PubKey[:]...), txscript.SigHashDefault, privKey,
	)
	signatureHash, err := txscript.CalcTaprootSignatureHash(sigHashes, txscript.SigHashDefault, packet.UnsignedTx, 0, fetcher)
	if err != nil {
		panic(err)
	}
	signature, err := schnorr.Sign(privKey, signatureHash)
	if err != nil {
		panic(err)
	}

	_ = witnessScript

	var witnessBytes bytes.Buffer
	err = psbt.WriteTxWitness(&witnessBytes, [][]byte{signature.Serialize()})
	if err != nil {
		panic(err)
	}

	packet.Inputs[0].FinalScriptWitness = witnessBytes.Bytes()
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
