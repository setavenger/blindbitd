package main

import (
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/daemon"
	"github.com/setavenger/blindbitd/src/ipc"
	"github.com/setavenger/blindbitd/src/networking"
	"github.com/setavenger/blindbitd/src/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {

	wallet := src.NewWallet()
	c := networking.Client{BaseUrl: "http://localhost:8000"}

	d := daemon.Daemon{
		Wallet: wallet,
		Client: &c,
	}

	d.Run()
	serverIpc := ipc.NewServer(&d)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	sGRpc := grpc.NewServer()
	reflection.Register(sGRpc)

	pb.RegisterIpcServiceServer(sGRpc, serverIpc)
	if err = sGRpc.Serve(listener); err != nil {
		panic(err)
	}
}
