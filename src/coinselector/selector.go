package coinselector

/*
Simplified as we don't expect to produce transactions with more than 252 inputs/outputs.
Witness data is also very standardised.
*/

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/utils"
	"math"
)

// FeeRateCoinSelector
// Custom CoinSelector implementation. Selects according to a given fee rate. Focused on taproot-only inputs.
// Needs the OwnedUTXOs to contain at least the Amount of the src.OwnedUTXO.
// The function will fail if not enough value could be added together.
// Other data in the OwnedUTXOs is preserved.
// At the moment it is always assumed that we receive a taproot input.
type FeeRateCoinSelector struct {
	OwnedUTXOs      src.UtxoCollection
	MinChangeAmount uint64
	Recipients      []*src.Recipient
}

// Length in bytes without witness discount

// See here for explanation of vByte sizes https://bitcoinops.org/en/tools/calc-size/
const (
	NTxVersionLen                   float64 = 4
	SegWitMarkerLenAndSegWitFlagLen float64 = 0.5
	NumInputsLen                    float64 = 1 // todo is a varInt

	TrOutpointTxidLen      float64 = 32
	TrOutpointVoutLen      float64 = 4
	TrOutpointSequenceLen  float64 = 4
	TrEmptyRedeemScriptLen float64 = 1
	TrInputOutpointLen     float64 = TrOutpointTxidLen + TrOutpointVoutLen + TrEmptyRedeemScriptLen + TrOutpointSequenceLen

	// TrWitnessDataLen already discounted by 0.25 complete length (varInt + actual data)
	TrWitnessDataLen float64 = 16.25

	OutputValueLen float64 = 8

	WitnessCountLen float64 = 1 // todo is a varInt

	NLockTimeLen float64 = 4

	ScriptPubKeyTaprootLen = 34
)

func NewFeeRateCoinSelector(utxos src.UtxoCollection, minChangeAmount uint64, recipients []*src.Recipient) *FeeRateCoinSelector {
	return &FeeRateCoinSelector{
		OwnedUTXOs:      utxos,
		MinChangeAmount: minChangeAmount,
		Recipients:      recipients,
	}
}

// CoinSelect
// returns the utxos to select and the change amount in order to achieve the desired fee rate
func (s *FeeRateCoinSelector) CoinSelect(feeRate uint32) (src.UtxoCollection, uint64, error) {
	// todo should we somehow expose the resulting vBytes for later analysis?
	// todo reduce complexity in this function
	if feeRate < 1 {
		return nil, 0, src.ErrInvalidFeeRate
	}
	// track vBytes of the transaction
	var vByte float64 // todo make sure we don't face any decimal imprecision

	// OVERHEAD will always be there
	vByte += NTxVersionLen + SegWitMarkerLenAndSegWitFlagLen + NLockTimeLen
	vByte += NumInputsLen

	//
	outputLens, err := extractPkScriptsFromRecipients(s.Recipients)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, 0, err
	}

	vByte += float64(wire.VarIntSerializeSize(uint64(len(outputLens))))

	// END OVERHEAD should be 10.5 vByte here

	// add outputs to vByte
	for _, scriptPubKeyLen := range outputLens {
		vByte += OutputValueLen + float64(wire.VarIntSerializeSize(uint64(scriptPubKeyLen))) + float64(scriptPubKeyLen)
	}

	var sumTargetAmount uint64
	for _, recipient := range s.Recipients {
		if recipient.Amount > 0 {
			sumTargetAmount += uint64(recipient.Amount)
		} else {
			return nil, 0, src.ErrRecipientAmountIsZero
		}
	}

	var selectedInputs src.UtxoCollection
	var sumSelectedInputsAmounts uint64
	//var potentialVBytes = vByte // tracks a potential increase before actually adding to the main vByte tracking

	for i, utxo := range s.OwnedUTXOs {
		_ = i
		// we check that the sum of selected input amounts exceeds the (target Value + fees + (min. change))
		selectedInputs = append(selectedInputs, utxo)
		sumSelectedInputsAmounts += utxo.Amount

		//if i == 0 {
		vByte += WitnessCountLen / 4
		//}

		// outpoint size
		vByte += TrInputOutpointLen
		vByte += TrWitnessDataLen

		// todo also check that the fee rate is as we want it
		if sumSelectedInputsAmounts > sumTargetAmount+NeededFeeAbsolutSats(vByte, feeRate) {
			if sumSelectedInputsAmounts-(sumTargetAmount+NeededFeeAbsolutSats(vByte, feeRate)) < s.MinChangeAmount {
				continue
			}
			// todo account that change was considered in the vByte tx size
			return selectedInputs, sumSelectedInputsAmounts - (sumTargetAmount + NeededFeeAbsolutSats(vByte, feeRate)), err
		}
	}

	return nil, 0, src.ErrInsufficientFunds

}

func extractPkScriptsFromRecipients(recipients []*src.Recipient) ([]int, error) {

	var pkScriptLens []int

	for _, recipient := range recipients {
		if recipient.PkScript != nil {
			// skip if a pkScript is already present (for what ever reason)
			pkScriptLens = append(pkScriptLens, len(recipient.PkScript))
			continue
		}
		isSP := utils.IsSilentPaymentAddress(recipient.Address)
		if isSP {
			pkScriptLens = append(pkScriptLens, ScriptPubKeyTaprootLen)
			continue
		}

		// do this for all non SP addresses
		address, err := btcutil.DecodeAddress(recipient.Address, src.ChainParams)
		if err != nil {
			logging.ErrorLogger.Printf("Failed to decode address: %v", err)
			return nil, err
		}
		scriptPubKey, err := txscript.PayToAddrScript(address)
		if err != nil {
			logging.ErrorLogger.Printf("Failed to create scriptPubKey: %v", err)
			return nil, err
		}
		pkScriptLens = append(pkScriptLens, len(scriptPubKey))
	}

	return pkScriptLens, nil
}

func NeededFeeAbsolutSats(vByte float64, feeRate uint32) uint64 {
	return uint64(math.Ceil(vByte * float64(feeRate)))
}
