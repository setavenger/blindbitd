package lib

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbitd/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func NewClient(socketPath string) (pb.IpcServiceClient, *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to the server with a timeout context
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("unix://%s", socketPath), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewIpcServiceClient(conn)
	return client, conn
}
