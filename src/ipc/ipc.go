package ipc

import (
	"context"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/daemon"
	"github.com/setavenger/blindbitd/src/pb"
	"github.com/setavenger/blindbitd/src/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedIpcServiceServer
	Daemon *daemon.Daemon
}

func NewServer(d *daemon.Daemon) *Server {
	return &Server{Daemon: d}
}

func (s *Server) Status(ctx context.Context, in *pb.Empty) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{Status: s.Daemon.Status}, nil
}

func (s Server) ListUTXOs(ctx context.Context, in *pb.Empty) (*pb.UTXOCollection, error) {
	return &pb.UTXOCollection{Utxos: convertWalletUTXOs(s.Daemon.UTXOs)}, nil
}

func convertWalletUTXOs(utxos []src.OwnedUTXO) []*pb.OwnedUTXO {
	var result []*pb.OwnedUTXO

	for _, utxo := range utxos {
		result = append(result, &pb.OwnedUTXO{
			Txid:               utils.CopyBytes(utxo.Txid[:]),
			Vout:               utxo.Vout,
			Amount:             utxo.Amount,
			PubKey:             utils.CopyBytes(utxo.PubKey[:]),
			TimestampConfirmed: &timestamppb.Timestamp{Seconds: int64(utxo.TimestampConfirmed)},
			UtxoState:          convertState(utxo.State),
			Label:              utxo.LabelPubKey(),
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
	default:
		return pb.UTXOState_UNKNOWN
	}
}
