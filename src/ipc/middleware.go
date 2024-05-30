package ipc

import (
	"context"

	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/logging"
	"google.golang.org/grpc"
)

type accessMode uint8

const (
	lockedDaemon accessMode = iota
	scanOnly
)

type AccessRule struct {
	Methods []string
	Allow   bool
}

var accessPolicies = map[accessMode][]AccessRule{
	// list of allowed endpoints
	lockedDaemon: {
		{Methods: []string{pb.IpcService_Status_FullMethodName, pb.IpcService_Unlock_FullMethodName, pb.IpcService_CreateNewLabel_FullMethodName, pb.IpcService_RecoverWallet_FullMethodName, pb.IpcService_SetupScanOnly_FullMethodName}, Allow: true},
		{Methods: []string{"*"}, Allow: false},
	},
	// list of forbidden endpoints
	scanOnly: {
		{Methods: []string{pb.IpcService_CreateTransaction_FullMethodName, pb.IpcService_CreateTransactionAndBroadcast_FullMethodName, pb.IpcService_GetMnemonic_FullMethodName, pb.IpcService_SetMnemonic_FullMethodName}, Allow: false},
		{Methods: []string{"*"}, Allow: true},
	},
}

// checkAccess checks whether a function can be accessed when the daemon is locked
func checkAccess(mode accessMode, method string) bool {
	if rules, exists := accessPolicies[mode]; exists {
	forCatchAll:
		for _, rule := range rules {
			for _, m := range rule.Methods {
				if m == method {
					return rule.Allow
				}
			}
		}

		// we change the method to the catch all and reiterate for the match in he rules. Could be simplified most likely.
		method = "*"
		goto forCatchAll
	}

	return false
}

func (s *Server) stateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if s.Daemon.Locked && !checkAccess(lockedDaemon, info.FullMethod) {
		logging.DebugLogger.Printf("Method %s not allowed when locked\n", info.FullMethod)
		return nil, src.ErrDaemonIsLocked
	}

	if src.ScanOnly && !checkAccess(scanOnly, info.FullMethod) {
		logging.DebugLogger.Printf("Method %s not allowed in scan only mode", info.FullMethod)
		return nil, src.ErrDaemonIsScanOnly
	}

	// Proceed with the handler if checks pass
	return handler(ctx, req)
}
