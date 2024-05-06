package ipc

import (
	"context"
	"errors"
	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/daemon"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strings"
	"time"
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

func (s *Server) SyncHeight(_ context.Context, _ *pb.Empty) (*pb.SyncHeightResponse, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	return &pb.SyncHeightResponse{Height: s.Daemon.Wallet.LastScanHeight}, nil
}

func (s *Server) Unlock(_ context.Context, in *pb.PasswordRequest) (*pb.BoolResponse, error) {
	// todo remove BoolResponse - just returning an error from server overrides this
	if !s.Daemon.Locked {
		return &pb.BoolResponse{Success: false, Error: "daemon already unlocked"}, errors.New("daemon already unlocked")
	}

	var response pb.BoolResponse

	s.Daemon.Status = pb.Status_STATUS_STARTING
	s.Daemon.Password = []byte(in.Password)
	if utils.CheckIfFileExists(src.PathToKeys) {
		err := s.Daemon.LoadDataFromDB()
		if err != nil {
			response.Success = false
			response.Error = err.Error()
			return &response, err
		}
	} else {
		response.Success = false
		response.Error = "keys not found"
		return &response, errors.New(response.Error)
	}

	response.Success = true
	response.Error = ""

	// send signal that wallet was unlocked successfully
	s.Daemon.ReadyChan <- struct{}{}
	s.Daemon.Status = pb.Status_STATUS_RUNNING
	return &response, nil
}

func (s *Server) SetPassword(_ context.Context, in *pb.PasswordRequest) (*pb.BoolResponse, error) {
	if in.Password == "" {
		return &pb.BoolResponse{Success: false, Error: "password can't be empty"}, errors.New("password can't be empty")
	}

	var response pb.BoolResponse
	if s.Daemon.Password != nil {
		response.Success = false
		response.Error = "already set"
		return &response, nil
	}
	s.Daemon.Password = []byte(in.Password)

	response.Success = true
	response.Error = ""

	// send signal that wallet was unlocked successfully
	s.Daemon.ReadyChan <- struct{}{}
	s.Daemon.Status = pb.Status_STATUS_RUNNING
	return &response, nil
}

func (s *Server) Shutdown(_ context.Context, _ *pb.Empty) (*pb.BoolResponse, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	var response pb.BoolResponse

	s.Daemon.Status = pb.Status_STATUS_SHUTTING_DOWN

	err := s.Daemon.Shutdown()
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return &response, err
	}

	response.Success = true
	s.Daemon.ShutdownChan <- struct{}{}

	return &response, nil
}

func (s *Server) ListUTXOs(_ context.Context, _ *pb.Empty) (*pb.UTXOCollection, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	return &pb.UTXOCollection{Utxos: convertWalletUTXOs(s.Daemon.Wallet.UTXOs, s.Daemon.Wallet.LabelsMapping)}, nil
}

func (s *Server) ListAddresses(_ context.Context, _ *pb.Empty) (*pb.AddressesCollection, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}

	// todo return addresses sorted by m and standard should be first
	var addressCollection pb.AddressesCollection
	walletAddresses, err := s.Daemon.Wallet.SortedAddresses()

	if err != nil {
		return nil, err
	}

	for _, address := range walletAddresses {
		addressCollection.Addresses = append(addressCollection.Addresses, &pb.Address{
			Address: address.Address,
			Comment: address.Comment,
		})
	}
	return &addressCollection, nil
}

func (s *Server) CreateNewLabel(_ context.Context, in *pb.NewLabelRequest) (*pb.Address, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	label, err := s.Daemon.Wallet.GenerateNewLabel(in.Comment)
	if err != nil {
		return nil, err
	}
	return &pb.Address{Address: label.Address, Comment: label.Comment}, nil
}

func (s *Server) CreateTransaction(_ context.Context, in *pb.CreateTransactionRequest) (*pb.RawTransaction, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	recipients := convertToRecipients(in.Recipients)
	// todo UTXOs have to be marked as spent after creating the transaction; broadcast and mark as spent
	signedTx, err := s.Daemon.SendToRecipients(recipients, in.FeeRate)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}
	return &pb.RawTransaction{RawTx: signedTx}, nil
}

func (s *Server) CreateTransactionAndBroadcast(_ context.Context, in *pb.CreateTransactionRequest) (*pb.NewTransaction, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	recipients := convertToRecipients(in.Recipients)
	// todo UTXOs have to be marked as spent after creating the transaction; broadcast and mark as spent
	signedTx, err := s.Daemon.SendToRecipients(recipients, in.FeeRate)
	if err != nil {
		return nil, err
	}
	txid, err := s.Daemon.BroadcastTx(signedTx)
	if err != nil {
		return nil, err
	}

	go func() {
		// give a delay such that the electrum server can update the state
		<-time.After(3 * time.Second)

		err = s.Daemon.CheckUnspentUTXOs()
		if err != nil {
			// we only log the error here as it is not relevant to the general execution of the ipc call
			logging.ErrorLogger.Println(err)
		}
	}()

	return &pb.NewTransaction{Txid: txid}, nil
}

func (s *Server) BroadcastRawTx(_ context.Context, in *pb.RawTransaction) (*pb.NewTransaction, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}

	txid, err := s.Daemon.BroadcastTx(in.RawTx)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	return &pb.NewTransaction{Txid: txid}, nil
}

func (s *Server) CreateNewWallet(_ context.Context, in *pb.NewWalletRequest) (*pb.Mnemonic, error) {
	var err error
	// todo add checks that existing wallet is not overridden
	// check if a keys file already exists // todo wrap in function for check
	if utils.CheckIfFileExists(src.PathToKeys) {
		return nil, errors.New("keys file already exists")
	}
	if utils.CheckIfFileExists(src.PathDbWallet) {
		return nil, errors.New("wallet file already exists")
	}

	if in.EncryptionPassword == "" {
		return nil, errors.New("encryption password can't be empty")
	}
	s.Daemon.Password = []byte(in.EncryptionPassword)

	s.Daemon.Locked = false // temporarily set locked to false in order to allow writing to files during process
	defer func() {
		if err == nil {
			return
		}
		s.Daemon.Locked = true
		err = os.Remove(src.PathToKeys)
		if err != nil {
			logging.ErrorLogger.Println(err)
			// don't kill here try to delete the other file as well
		}
		err = os.Remove(src.PathDbWallet)
		if err != nil {
			logging.ErrorLogger.Println(err)
			// don't kill here try to delete the other file as well
		}
	}()

	err = s.Daemon.CreateNewKeys(in.SeedPassphrase)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	s.Daemon.ReadyChan <- struct{}{}
	s.Daemon.Status = pb.Status_STATUS_RUNNING

	return &pb.Mnemonic{Mnemonic: s.Daemon.Mnemonic}, nil
}

func (s *Server) RecoverWallet(_ context.Context, in *pb.RecoverWalletRequest) (*pb.BoolResponse, error) {
	var err error
	var response pb.BoolResponse
	if utils.CheckIfFileExists(src.PathToKeys) {
		response.Success = false
		response.Error = "keys file already exists"
		return &response, errors.New(response.Error)
	}
	if utils.CheckIfFileExists(src.PathDbWallet) {
		return nil, errors.New("wallet file already exists")
	}
	if in.EncryptionPassword == "" {
		response.Success = false
		response.Error = "encryption password can't be empty"
		return &response, errors.New(response.Error)
	}
	s.Daemon.Password = []byte(in.EncryptionPassword)

	var seedPassphrase string
	if in.SeedPassphrase != nil {
		seedPassphrase = *in.SeedPassphrase
	}

	s.Daemon.Locked = false // temporarily set locked to false in order to allow writing to files during process
	defer func() {
		if err == nil {
			return
		}
		s.Daemon.Locked = true
		err = os.Remove(src.PathToKeys)
		if err != nil {
			logging.ErrorLogger.Println(err)
			// don't kill here try to delete the other file as well
		}
		err = os.Remove(src.PathDbWallet)
		if err != nil {
			logging.ErrorLogger.Println(err)
			// don't kill here try to delete the other file as well
		}
	}()

	err = s.Daemon.RecoverFromSeed(in.Mnemonic, seedPassphrase, in.BirthHeight)
	if err != nil {
		logging.ErrorLogger.Println(err)
		response.Success = false
		response.Error = err.Error()
		return &response, err
	}

	s.Daemon.ReadyChan <- struct{}{}
	s.Daemon.Status = pb.Status_STATUS_RUNNING

	response.Success = true
	return &response, err
}

func (s *Server) ForceRescanFromHeight(_ context.Context, in *pb.RescanRequest) (*pb.BoolResponse, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	var response pb.BoolResponse

	err := s.Daemon.ForceSyncFrom(in.GetHeight())
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		return &response, err
	}

	response.Success = true

	return &response, nil
}

func (s *Server) GetMnemonic(_ context.Context, _ *pb.Empty) (*pb.Mnemonic, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	return &pb.Mnemonic{Mnemonic: s.Daemon.Mnemonic}, nil
}

func (s *Server) GetChain(_ context.Context, _ *pb.Empty) (*pb.Chain, error) {
	if s.Daemon.Locked {
		return nil, src.ErrDaemonIsLocked
	}
	return convertChainParam(src.ChainParams), nil
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

	//sGRpc.GracefulStop() //todo this via channel that is fed from shutdown

	pb.RegisterIpcServiceServer(sGRpc, s)
	if err = sGRpc.Serve(listener); err != nil {
		panic(err)
	}
	return err
}
