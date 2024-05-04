package ipc

import (
	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertWalletUTXOs(utxos []*src.OwnedUTXO, mapping src.LabelsMapping) []*pb.OwnedUTXO {
	var result []*pb.OwnedUTXO

	for _, utxo := range utxos {
		result = append(result, &pb.OwnedUTXO{
			Txid:               utils.CopyBytes(utxo.Txid[:]),
			Vout:               utxo.Vout,
			Amount:             utxo.Amount,
			PubKey:             utils.CopyBytes(utxo.PubKey[:]),
			TimestampConfirmed: &timestamppb.Timestamp{Seconds: int64(utxo.Timestamp)},
			UtxoState:          convertState(utxo.State),
			Label:              utxo.LabelComment(mapping),
		})
	}

	return result
}

func convertState(state src.UTXOState) pb.UTXOState {
	switch state {
	case src.StateUnconfirmed:
		return pb.UTXOState_UNCONFIRMED
	case src.StateUnspent:
		return pb.UTXOState_UNSPENT
	case src.StateSpent:
		return pb.UTXOState_SPENT
	case src.StateUnconfirmedSpent:
		return pb.UTXOState_SPENT_UNCONFIRMED
	default:
		return pb.UTXOState_UNKNOWN
	}
}

func convertToRecipients(recipients []*pb.TransactionRecipient) []*src.Recipient {
	var convertedRecipients []*src.Recipient

	for _, recipient := range recipients {
		convertedRecipients = append(convertedRecipients, &src.Recipient{
			Address:    recipient.Address,
			Amount:     int64(recipient.Amount),
			Annotation: recipient.Annotation,
		})
	}

	return convertedRecipients
}
