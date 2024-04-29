package ipc

import (
	"context"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/daemon"
	"github.com/setavenger/blindbitd/src/database"
	"github.com/setavenger/blindbitd/src/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strings"
)

type Server struct {
	pb.UnimplementedIpcServiceServer
	Daemon *daemon.Daemon
}

func NewServer(d *daemon.Daemon) *Server {
	return &Server{Daemon: d}
}

func (s *Server) Status(_ context.Context, _ *pb.Empty) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{Status: s.Daemon.Status}, nil
}

func (s *Server) Unlock(_ context.Context, in *pb.PasswordRequest) (*pb.BoolResponse, error) {

	var response pb.BoolResponse
	var wallet src.Wallet

	s.Daemon.Status = pb.Status_STATUS_STARTING
	s.Daemon.Password = []byte(in.Password)

	err := database.ReadFromDB(src.PathDbWallet, &wallet, s.Daemon.Password)
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return &response, err
	}
	s.Daemon.Wallet = &wallet

	response.Success = true
	response.Error = ""

	// send signal that wallet was unlocked successfully
	s.Daemon.ReadyChan <- struct{}{}
	s.Daemon.Status = pb.Status_STATUS_RUNNING
	return &response, nil
}

func (s *Server) SetPassword(_ context.Context, in *pb.PasswordRequest) (*pb.BoolResponse, error) {
	var response pb.BoolResponse
	s.Daemon.Password = []byte(in.Password)

	response.Success = true
	response.Error = ""

	// send signal that wallet was unlocked successfully
	s.Daemon.ReadyChan <- struct{}{}
	s.Daemon.Status = pb.Status_STATUS_RUNNING
	return &response, nil
}

func (s *Server) Shutdown(_ context.Context, _ *pb.Empty) (*pb.BoolResponse, error) {

	var response pb.BoolResponse

	s.Daemon.Status = pb.Status_STATUS_SHUTTING_DOWN

	err := s.Daemon.Shutdown()
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return &response, err
	}

	response.Success = true

	return &response, nil
}

func (s *Server) ListUTXOs(_ context.Context, _ *pb.Empty) (*pb.UTXOCollection, error) {
	return &pb.UTXOCollection{Utxos: convertWalletUTXOs(s.Daemon.UTXOs)}, nil
}

func (s *Server) CreateTransaction(_ context.Context, in *pb.CreateTransactionRequest) (*pb.RawTransaction, error) {
	recipients := convertToRecipients(in.Recipients)
	signedTx, err := s.Daemon.SendToRecipients(recipients, in.FeeRate)
	if err != nil {
		return nil, err
	}
	return &pb.RawTransaction{RawTx: signedTx}, nil
}

func (s *Server) Start() error {
	if s.Daemon == nil {
		return src.ErrDaemonNotSet
	}

	err := os.Remove(src.PathIpcSocket)
	if err != nil && !strings.Contains(err.Error(), "no such file") {
		panic(err)
	}
	listener, err := net.Listen("unix", src.PathIpcSocket)
	if err != nil {
		panic(err)
	}

	sGRpc := grpc.NewServer()
	reflection.Register(sGRpc)

	pb.RegisterIpcServiceServer(sGRpc, s)
	if err = sGRpc.Serve(listener); err != nil {
		panic(err)
	}
	return err
}
