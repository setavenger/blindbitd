package daemon

import (
	"context"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/pb"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/database"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/networking"
	"github.com/setavenger/blindbitd/src/utils"
	"github.com/setavenger/go-electrum/electrum"
)

type Daemon struct {
	Status         pb.Status
	Password       []byte
	Locked         bool
	ReadyChan      chan struct{} // for the startup signal; either unlocking or setting password on initial startup
	ShutdownChan   chan struct{}
	Mnemonic       string
	ClientElectrum *electrum.Client
	ClientBlindBit *networking.ClientBlindBit
	Wallet         *src.Wallet
	NewBlockChan   <-chan *electrum.SubscribeHeadersResult
}

func NewDaemon(wallet *src.Wallet, clientBlindBit *networking.ClientBlindBit, clientElectrum *electrum.Client, network *chaincfg.Params) *Daemon {
	channel, err := clientElectrum.SubscribeHeaders(context.Background())
	if err != nil {
		panic(err)
	}
	return &Daemon{
		Status:         pb.Status_STATUS_UNSPECIFIED,
		Wallet:         wallet,
		ClientBlindBit: clientBlindBit,
		ClientElectrum: clientElectrum,
		Locked:         true,
		ReadyChan:      make(chan struct{}),
		ShutdownChan:   make(chan struct{}),
		NewBlockChan:   channel,
	}
}

func (d *Daemon) Run() error {
	d.Status = pb.Status_STATUS_RUNNING

	// first we sync up and then we scan continuously
	err := d.SyncToTip(0)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	logging.InfoLogger.Println("Balance:", d.Wallet.FreeBalance())

	// todo add a recovery mechanism
	err = d.ContinuousScan() // blocking function if it returns, it returns an error and Run is closed as well
	return err
}

var exampleLabelComments = [5]string{"Hello", "Donations for project", "Family and Friends", "Deal 1", "Deal 2"}

// LoadDataFromDB
// Load keys and wallet data from disk
func (d *Daemon) LoadDataFromDB() error {
	var keys src.Keys
	d.Status = pb.Status_STATUS_STARTING
	err := database.ReadFromDB(src.PathToKeys, &keys, d.Password)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	var wallet src.Wallet

	d.Mnemonic = keys.Mnemonic
	// load keys in any case other data will be read in next step if available
	wallet.LoadKeys(keys.ScanSecretKey, keys.SpendSecretKey)
	err = wallet.CheckAndInitialiseFields()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	if utils.CheckIfFileExists(src.PathDbWallet) {
		err = database.ReadFromDB(src.PathDbWallet, &wallet, d.Password)
		if err != nil {
			logging.ErrorLogger.Println(err)
			return err
		}
	}

	d.Wallet = &wallet

	return nil
}

func (d *Daemon) Shutdown() error {
	// todo save all data to a files
	fmt.Println("Process shutting down")

	if d.ClientElectrum != nil {
		d.ClientElectrum.Shutdown()
	}
	if d.Status == pb.Status_STATUS_NO_WALLET {
		// we don't store anything if the wallet was not initialised yet
		return nil
	}
	if d.Locked || d.Password == nil {
		return nil
	}
	err := database.WriteToDB(src.PathDbWallet, d.Wallet, d.Password)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}
	return nil
}

// CreateNewKeys
// WARNING: Must only be called if no other wallet is present. Will overwrite the old keys.
func (d *Daemon) CreateNewKeys(seedPassphrase string) error {

	var chainTip uint64
	chainTip, err := d.ClientBlindBit.GetChainTip()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	d.Wallet = src.NewWallet(chainTip)
	var newKeys *src.Keys
	newKeys, err = src.CreateNewKeys(seedPassphrase)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}
	d.Wallet.LoadKeys(newKeys.ScanSecretKey, newKeys.SpendSecretKey)
	if newKeys.Mnemonic == "" {
		return errors.New("mnemonic is empty")
	}
	d.Mnemonic = newKeys.Mnemonic
	if d.Locked || d.Password == nil {
		return errors.New("daemon is locked or has no encryption password")
	}
	err = database.WriteToDB(src.PathToKeys, newKeys, d.Password)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	// setup up the other important stuff needed
	err = d.Wallet.CheckAndInitialiseFields()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	_, err = d.Wallet.GenerateAddress()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func (d *Daemon) RecoverFromSeed(mnemonic, seedPassphrase string, birthHeight uint64) error {

	d.Wallet = src.NewWallet(birthHeight)
	newKeys, err := src.KeysFromMnemonic(mnemonic, seedPassphrase)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}
	d.Wallet.LoadKeys(newKeys.ScanSecretKey, newKeys.SpendSecretKey)
	if newKeys.Mnemonic == "" {
		return errors.New("mnemonic is empty")
	}
	d.Mnemonic = newKeys.Mnemonic
	if d.Locked || d.Password == nil {
		return errors.New("daemon is locked or has no encryption password")
	}
	err = database.WriteToDB(src.PathToKeys, newKeys, d.Password)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	// setup up the other important stuff needed
	err = d.Wallet.CheckAndInitialiseFields()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	_, err = d.Wallet.GenerateAddress()
	if err != nil {
		logging.ErrorLogger.Println(err)
		return err
	}

	return err
}
