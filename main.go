package main

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/daemon"
	"github.com/setavenger/blindbitd/src/ipc"
	"github.com/setavenger/blindbitd/src/networking"
	"github.com/setavenger/blindbitd/src/pb"
	"os"
	"os/signal"
)

func main() {
	defer func() {
		err := os.Remove(src.PathIpcSocket)
		if err != nil {
			panic(err)
		}
	}()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// load settings
	src.SetPaths("./.blindbit")

	// create the daemon but locked and without Wallet data
	c := networking.Client{BaseUrl: "http://localhost:8000"}

	d := daemon.NewDaemon(nil, &c, &chaincfg.RegressionNetParams)

	serverIpc := ipc.NewServer(d)

	go func() {
		// is blocking hence go routine
		err := serverIpc.Start()
		if err != nil {
			return
		}
		d.Status = pb.Status_STATUS_STARTING
	}()

	if _, err := os.Stat(src.PathDbWallet); err == nil {
		// exists and needs to be unlocked
		fmt.Println("Waiting to be unlocked...")
		select {
		// Wait here until daemon is unlocked
		case <-d.ReadyChan:
			fmt.Println("Daemon is ready...")
		case <-interrupt:
			return
		}
	} else if errors.Is(err, os.ErrNotExist) {
		//  does *not* exist
		fmt.Println("Please set password...")
		select {
		// Wait here until daemon is unlocked
		case <-d.ReadyChan:
			fmt.Println("Daemon is ready...")
		case <-interrupt:
			return
		}
		var chainTip uint64
		chainTip, err = c.GetChainTip()
		if err != nil {
			panic(err)
		}
		d.Wallet = src.NewWallet(chainTip)
	} else {
		// Schrödinger: file may or may not exist. See err for details.
		panic(err)
	}

	d.Run()

	for {
		select {
		case <-interrupt:
			return
		}
	}

}
