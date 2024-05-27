package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/daemon"
	"github.com/setavenger/blindbitd/src/ipc"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/networking"
	"github.com/setavenger/blindbitd/src/utils"
	"github.com/setavenger/go-electrum/electrum"
)

var testEnvironment bool
var dataDirectory string

func init() {
	// todo can this double reference work?
	flag.StringVar(&dataDirectory, "datadir", src.DefaultDirectoryPath, "Set the base directory for the blindbit daemon. Default directory is ~/.blindbitd")
	flag.BoolVar(&testEnvironment, "test", false, "NEVER USE IN PRODUCTION. If set to true the program will load predefined test keys")
	flag.Parse()
}

func main() {
	defer func() {
		if exists := utils.CheckIfFileExists(src.PathIpcSocket); exists {
			err := os.Remove(src.PathIpcSocket)
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}
		}
		fmt.Println("blindbitd shut down")
	}()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	d, err := daemon.NewDaemon(nil, nil, nil) // initialise with nothing to avoid deadlocks and nil pointers down below
	if err != nil {
		panic(err)
	}

	defer func() {
		err = d.Shutdown()
		if err != nil {
			logging.ErrorLogger.Println(err)
			panic(err)
		}
	}()

	go func() {
		// todo remove after development. Replace with command line arg
		src.SetPaths(dataDirectory)

		// initialise loggers
		logging.LoadLoggers(src.PathLogs)
		logging.DebugLogger.Println("paths loaded")

		// load config settings
		src.LoadConfigs(src.PathConfig)
		logging.DebugLogger.Println("config loaded")

		// create the daemon but locked and without Wallet data
		clientBlindBit := networking.ClientBlindBit{BaseUrl: src.BlindBitServerAddress}
		var clientElectrum *electrum.Client
		var err error

		if src.UseElectrum {
			logging.DebugLogger.Println("connecting to Electrum server")
			clientElectrum, err = networking.CreateElectrumClient()
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}
		}

		d, err = daemon.NewDaemon(nil, &clientBlindBit, clientElectrum)
		if err != nil {
			logging.ErrorLogger.Println(err)
			panic(err)
		}
		d.Status = pb.Status_STATUS_STARTING

		serverIpc := ipc.NewServer(d)

		go func() {
			logging.DebugLogger.Println("Starting IPC server")
			// is blocking hence go routine
			err = serverIpc.Start()
			if err != nil {
				return
			}
		}()

		// todo can this be more robust, especially considering the different unlocking/initialisation paths available
		if utils.CheckIfFileExists(src.PathToKeys) {
			d.Status = pb.Status_STATUS_LOCKED
			// exists and needs to be unlocked
			logging.InfoLogger.Println("Waiting to be unlocked...")
			select {
			// Wait here until daemon is unlocked
			case <-d.ReadyChan:
				logging.InfoLogger.Println("Daemon is ready...")
				d.Locked = false
			case <-interrupt:
				return
			}
		} else {
			// does *not* exist
			d.Status = pb.Status_STATUS_NO_WALLET
			logging.InfoLogger.Println("Please create new wallet...")
			select {
			// Wait here until wallet is set up
			case <-d.ReadyChan:
				logging.InfoLogger.Println("Daemon is ready...")
				d.Locked = false
			case <-interrupt:
				return
			}
			logging.InfoLogger.Println("New wallet created")
		}

		d.Status = pb.Status_STATUS_STARTING

		if testEnvironment {
			err = d.LoadTestData()
			if err != nil {
				logging.ErrorLogger.Println(err)
				return
			}
		}

		err = d.Wallet.CheckAndInitialiseFields()
		if err != nil {
			logging.ErrorLogger.Println(err)
			return
		}

		go func() {
			err = d.Run()
			if err != nil {
				logging.ErrorLogger.Println(err)
				panic(err)
			}
		}()

	}()

	for {
		select {
		case <-d.ShutdownChan:
			fmt.Println("Daemon is shutting down...")
			return
		case <-interrupt:
			return
		}
	}
}
