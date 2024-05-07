// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: ipc.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	IpcService_Status_FullMethodName                        = "/ipc.IpcService/Status"
	IpcService_SyncHeight_FullMethodName                    = "/ipc.IpcService/SyncHeight"
	IpcService_Unlock_FullMethodName                        = "/ipc.IpcService/Unlock"
	IpcService_SetPassword_FullMethodName                   = "/ipc.IpcService/SetPassword"
	IpcService_Shutdown_FullMethodName                      = "/ipc.IpcService/Shutdown"
	IpcService_ListUTXOs_FullMethodName                     = "/ipc.IpcService/ListUTXOs"
	IpcService_ListAddresses_FullMethodName                 = "/ipc.IpcService/ListAddresses"
	IpcService_CreateNewLabel_FullMethodName                = "/ipc.IpcService/CreateNewLabel"
	IpcService_CreateTransaction_FullMethodName             = "/ipc.IpcService/CreateTransaction"
	IpcService_CreateTransactionAndBroadcast_FullMethodName = "/ipc.IpcService/CreateTransactionAndBroadcast"
	IpcService_BroadcastRawTx_FullMethodName                = "/ipc.IpcService/BroadcastRawTx"
	IpcService_GetMnemonic_FullMethodName                   = "/ipc.IpcService/GetMnemonic"
	IpcService_SetMnemonic_FullMethodName                   = "/ipc.IpcService/SetMnemonic"
	IpcService_CreateNewWallet_FullMethodName               = "/ipc.IpcService/CreateNewWallet"
	IpcService_RecoverWallet_FullMethodName                 = "/ipc.IpcService/RecoverWallet"
	IpcService_ForceRescanFromHeight_FullMethodName         = "/ipc.IpcService/ForceRescanFromHeight"
	IpcService_GetChain_FullMethodName                      = "/ipc.IpcService/GetChain"
)

// IpcServiceClient is the client API for IpcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IpcServiceClient interface {
	// Alive pings the daemon and returns true if the daemon is alive
	Status(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatusResponse, error)
	SyncHeight(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SyncHeightResponse, error)
	Unlock(ctx context.Context, in *PasswordRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	SetPassword(ctx context.Context, in *PasswordRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BoolResponse, error)
	ListUTXOs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UTXOCollection, error)
	ListAddresses(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AddressesCollection, error)
	CreateNewLabel(ctx context.Context, in *NewLabelRequest, opts ...grpc.CallOption) (*Address, error)
	CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*RawTransaction, error)
	CreateTransactionAndBroadcast(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*NewTransaction, error)
	BroadcastRawTx(ctx context.Context, in *RawTransaction, opts ...grpc.CallOption) (*NewTransaction, error)
	GetMnemonic(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Mnemonic, error)
	SetMnemonic(ctx context.Context, in *Mnemonic, opts ...grpc.CallOption) (*BoolResponse, error)
	CreateNewWallet(ctx context.Context, in *NewWalletRequest, opts ...grpc.CallOption) (*Mnemonic, error)
	RecoverWallet(ctx context.Context, in *RecoverWalletRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	ForceRescanFromHeight(ctx context.Context, in *RescanRequest, opts ...grpc.CallOption) (*BoolResponse, error)
	GetChain(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Chain, error)
}

type ipcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIpcServiceClient(cc grpc.ClientConnInterface) IpcServiceClient {
	return &ipcServiceClient{cc}
}

func (c *ipcServiceClient) Status(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, IpcService_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) SyncHeight(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SyncHeightResponse, error) {
	out := new(SyncHeightResponse)
	err := c.cc.Invoke(ctx, IpcService_SyncHeight_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) Unlock(ctx context.Context, in *PasswordRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, IpcService_Unlock_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) SetPassword(ctx context.Context, in *PasswordRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, IpcService_SetPassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, IpcService_Shutdown_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) ListUTXOs(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UTXOCollection, error) {
	out := new(UTXOCollection)
	err := c.cc.Invoke(ctx, IpcService_ListUTXOs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) ListAddresses(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AddressesCollection, error) {
	out := new(AddressesCollection)
	err := c.cc.Invoke(ctx, IpcService_ListAddresses_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) CreateNewLabel(ctx context.Context, in *NewLabelRequest, opts ...grpc.CallOption) (*Address, error) {
	out := new(Address)
	err := c.cc.Invoke(ctx, IpcService_CreateNewLabel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) CreateTransaction(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*RawTransaction, error) {
	out := new(RawTransaction)
	err := c.cc.Invoke(ctx, IpcService_CreateTransaction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) CreateTransactionAndBroadcast(ctx context.Context, in *CreateTransactionRequest, opts ...grpc.CallOption) (*NewTransaction, error) {
	out := new(NewTransaction)
	err := c.cc.Invoke(ctx, IpcService_CreateTransactionAndBroadcast_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) BroadcastRawTx(ctx context.Context, in *RawTransaction, opts ...grpc.CallOption) (*NewTransaction, error) {
	out := new(NewTransaction)
	err := c.cc.Invoke(ctx, IpcService_BroadcastRawTx_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) GetMnemonic(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Mnemonic, error) {
	out := new(Mnemonic)
	err := c.cc.Invoke(ctx, IpcService_GetMnemonic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) SetMnemonic(ctx context.Context, in *Mnemonic, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, IpcService_SetMnemonic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) CreateNewWallet(ctx context.Context, in *NewWalletRequest, opts ...grpc.CallOption) (*Mnemonic, error) {
	out := new(Mnemonic)
	err := c.cc.Invoke(ctx, IpcService_CreateNewWallet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) RecoverWallet(ctx context.Context, in *RecoverWalletRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, IpcService_RecoverWallet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) ForceRescanFromHeight(ctx context.Context, in *RescanRequest, opts ...grpc.CallOption) (*BoolResponse, error) {
	out := new(BoolResponse)
	err := c.cc.Invoke(ctx, IpcService_ForceRescanFromHeight_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ipcServiceClient) GetChain(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Chain, error) {
	out := new(Chain)
	err := c.cc.Invoke(ctx, IpcService_GetChain_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IpcServiceServer is the server API for IpcService service.
// All implementations must embed UnimplementedIpcServiceServer
// for forward compatibility
type IpcServiceServer interface {
	// Alive pings the daemon and returns true if the daemon is alive
	Status(context.Context, *Empty) (*StatusResponse, error)
	SyncHeight(context.Context, *Empty) (*SyncHeightResponse, error)
	Unlock(context.Context, *PasswordRequest) (*BoolResponse, error)
	SetPassword(context.Context, *PasswordRequest) (*BoolResponse, error)
	Shutdown(context.Context, *Empty) (*BoolResponse, error)
	ListUTXOs(context.Context, *Empty) (*UTXOCollection, error)
	ListAddresses(context.Context, *Empty) (*AddressesCollection, error)
	CreateNewLabel(context.Context, *NewLabelRequest) (*Address, error)
	CreateTransaction(context.Context, *CreateTransactionRequest) (*RawTransaction, error)
	CreateTransactionAndBroadcast(context.Context, *CreateTransactionRequest) (*NewTransaction, error)
	BroadcastRawTx(context.Context, *RawTransaction) (*NewTransaction, error)
	GetMnemonic(context.Context, *Empty) (*Mnemonic, error)
	SetMnemonic(context.Context, *Mnemonic) (*BoolResponse, error)
	CreateNewWallet(context.Context, *NewWalletRequest) (*Mnemonic, error)
	RecoverWallet(context.Context, *RecoverWalletRequest) (*BoolResponse, error)
	ForceRescanFromHeight(context.Context, *RescanRequest) (*BoolResponse, error)
	GetChain(context.Context, *Empty) (*Chain, error)
	mustEmbedUnimplementedIpcServiceServer()
}

// UnimplementedIpcServiceServer must be embedded to have forward compatible implementations.
type UnimplementedIpcServiceServer struct {
}

func (UnimplementedIpcServiceServer) Status(context.Context, *Empty) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedIpcServiceServer) SyncHeight(context.Context, *Empty) (*SyncHeightResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncHeight not implemented")
}
func (UnimplementedIpcServiceServer) Unlock(context.Context, *PasswordRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlock not implemented")
}
func (UnimplementedIpcServiceServer) SetPassword(context.Context, *PasswordRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPassword not implemented")
}
func (UnimplementedIpcServiceServer) Shutdown(context.Context, *Empty) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}
func (UnimplementedIpcServiceServer) ListUTXOs(context.Context, *Empty) (*UTXOCollection, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUTXOs not implemented")
}
func (UnimplementedIpcServiceServer) ListAddresses(context.Context, *Empty) (*AddressesCollection, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAddresses not implemented")
}
func (UnimplementedIpcServiceServer) CreateNewLabel(context.Context, *NewLabelRequest) (*Address, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewLabel not implemented")
}
func (UnimplementedIpcServiceServer) CreateTransaction(context.Context, *CreateTransactionRequest) (*RawTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransaction not implemented")
}
func (UnimplementedIpcServiceServer) CreateTransactionAndBroadcast(context.Context, *CreateTransactionRequest) (*NewTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransactionAndBroadcast not implemented")
}
func (UnimplementedIpcServiceServer) BroadcastRawTx(context.Context, *RawTransaction) (*NewTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastRawTx not implemented")
}
func (UnimplementedIpcServiceServer) GetMnemonic(context.Context, *Empty) (*Mnemonic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMnemonic not implemented")
}
func (UnimplementedIpcServiceServer) SetMnemonic(context.Context, *Mnemonic) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMnemonic not implemented")
}
func (UnimplementedIpcServiceServer) CreateNewWallet(context.Context, *NewWalletRequest) (*Mnemonic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewWallet not implemented")
}
func (UnimplementedIpcServiceServer) RecoverWallet(context.Context, *RecoverWalletRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverWallet not implemented")
}
func (UnimplementedIpcServiceServer) ForceRescanFromHeight(context.Context, *RescanRequest) (*BoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForceRescanFromHeight not implemented")
}
func (UnimplementedIpcServiceServer) GetChain(context.Context, *Empty) (*Chain, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChain not implemented")
}
func (UnimplementedIpcServiceServer) mustEmbedUnimplementedIpcServiceServer() {}

// UnsafeIpcServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IpcServiceServer will
// result in compilation errors.
type UnsafeIpcServiceServer interface {
	mustEmbedUnimplementedIpcServiceServer()
}

func RegisterIpcServiceServer(s grpc.ServiceRegistrar, srv IpcServiceServer) {
	s.RegisterService(&IpcService_ServiceDesc, srv)
}

func _IpcService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).Status(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_SyncHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).SyncHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_SyncHeight_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).SyncHeight(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_Unlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).Unlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_Unlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).Unlock(ctx, req.(*PasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_SetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).SetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_SetPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).SetPassword(ctx, req.(*PasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_Shutdown_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).Shutdown(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_ListUTXOs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).ListUTXOs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_ListUTXOs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).ListUTXOs(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_ListAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).ListAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_ListAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).ListAddresses(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_CreateNewLabel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewLabelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).CreateNewLabel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_CreateNewLabel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).CreateNewLabel(ctx, req.(*NewLabelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_CreateTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).CreateTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_CreateTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).CreateTransaction(ctx, req.(*CreateTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_CreateTransactionAndBroadcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).CreateTransactionAndBroadcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_CreateTransactionAndBroadcast_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).CreateTransactionAndBroadcast(ctx, req.(*CreateTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_BroadcastRawTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawTransaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).BroadcastRawTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_BroadcastRawTx_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).BroadcastRawTx(ctx, req.(*RawTransaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_GetMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).GetMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_GetMnemonic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).GetMnemonic(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_SetMnemonic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Mnemonic)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).SetMnemonic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_SetMnemonic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).SetMnemonic(ctx, req.(*Mnemonic))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_CreateNewWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).CreateNewWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_CreateNewWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).CreateNewWallet(ctx, req.(*NewWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_RecoverWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecoverWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).RecoverWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_RecoverWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).RecoverWallet(ctx, req.(*RecoverWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_ForceRescanFromHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RescanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).ForceRescanFromHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_ForceRescanFromHeight_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).ForceRescanFromHeight(ctx, req.(*RescanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IpcService_GetChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IpcServiceServer).GetChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IpcService_GetChain_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IpcServiceServer).GetChain(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// IpcService_ServiceDesc is the grpc.ServiceDesc for IpcService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IpcService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ipc.IpcService",
	HandlerType: (*IpcServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _IpcService_Status_Handler,
		},
		{
			MethodName: "SyncHeight",
			Handler:    _IpcService_SyncHeight_Handler,
		},
		{
			MethodName: "Unlock",
			Handler:    _IpcService_Unlock_Handler,
		},
		{
			MethodName: "SetPassword",
			Handler:    _IpcService_SetPassword_Handler,
		},
		{
			MethodName: "Shutdown",
			Handler:    _IpcService_Shutdown_Handler,
		},
		{
			MethodName: "ListUTXOs",
			Handler:    _IpcService_ListUTXOs_Handler,
		},
		{
			MethodName: "ListAddresses",
			Handler:    _IpcService_ListAddresses_Handler,
		},
		{
			MethodName: "CreateNewLabel",
			Handler:    _IpcService_CreateNewLabel_Handler,
		},
		{
			MethodName: "CreateTransaction",
			Handler:    _IpcService_CreateTransaction_Handler,
		},
		{
			MethodName: "CreateTransactionAndBroadcast",
			Handler:    _IpcService_CreateTransactionAndBroadcast_Handler,
		},
		{
			MethodName: "BroadcastRawTx",
			Handler:    _IpcService_BroadcastRawTx_Handler,
		},
		{
			MethodName: "GetMnemonic",
			Handler:    _IpcService_GetMnemonic_Handler,
		},
		{
			MethodName: "SetMnemonic",
			Handler:    _IpcService_SetMnemonic_Handler,
		},
		{
			MethodName: "CreateNewWallet",
			Handler:    _IpcService_CreateNewWallet_Handler,
		},
		{
			MethodName: "RecoverWallet",
			Handler:    _IpcService_RecoverWallet_Handler,
		},
		{
			MethodName: "ForceRescanFromHeight",
			Handler:    _IpcService_ForceRescanFromHeight_Handler,
		},
		{
			MethodName: "GetChain",
			Handler:    _IpcService_GetChain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ipc.proto",
}
