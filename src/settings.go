package src

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/blindbitd/src/utils"
	"github.com/spf13/viper"
)

const DefaultDirectoryPath = "~/.blindbitd"

var (
	DirectoryPath = "~/.blindbitd"

	PathLogs string

	PathIpcSocketDir string
	PathIpcSocket    string

	PathConfig string

	PathDbWallet string

	PathToKeys string
)

const PathEndingSocketDirPath = "/run"
const PathEndingSocketPath = PathEndingSocketDirPath + "/blindbit.socket"

const logsPath = "/logs"
const dataPath = "/data"

const PathEndingConfig = "/blindbit.toml"

const PathEndingWallet = dataPath + "/wallet"

const PathEndingKeys = dataPath + "/keys"

func SetPaths(baseDirectory string) {
	if baseDirectory != "" {
		DirectoryPath = baseDirectory
	}

	PathLogs = DirectoryPath + logsPath
	PathIpcSocketDir = DirectoryPath + PathEndingSocketDirPath
	PathIpcSocket = DirectoryPath + PathEndingSocketPath

	PathConfig = DirectoryPath + PathEndingConfig
	PathDbWallet = DirectoryPath + PathEndingWallet

	PathToKeys = DirectoryPath + PathEndingKeys

	// create the directories
	utils.TryCreateDirectoryPanic(DirectoryPath)
	utils.TryCreateDirectoryPanic(PathIpcSocketDir)

	utils.TryCreateDirectoryPanic(DirectoryPath + dataPath)
	utils.TryCreateDirectoryPanic(PathLogs)

	return
}

func LoadConfigs(pathToConfig string) {
	// Set the file name of the configurations file
	viper.SetConfigFile(pathToConfig)

	// Handle errors reading the config file
	if err := viper.ReadInConfig(); err != nil {
		logging.ErrorLogger.Fatalf("Error reading config file, %s", err)
	}

	/* set defaults */
	viper.SetDefault("blindbit_server", "http://localhost:8000")
	viper.SetDefault("electrum_server", "localhost:50000")
	viper.SetDefault("minchange_amount", 1000)
	viper.SetDefault("chain", "signet")

	/* read and set config variables */
	BlindBitServerAddress = viper.GetString("network.blindbit_server")
	ElectrumServerAddress = viper.GetString("network.electrum_server")

	MinChangeAmount = viper.GetInt64("wallet.minchange_amount")

	// extract the chain data and set the params
	chain := viper.GetString("network.chain")
	switch chain {
	case "main":
		ChainParams = &chaincfg.MainNetParams
	case "test":
		ChainParams = &chaincfg.TestNet3Params
	case "signet":
		ChainParams = &chaincfg.SigNetParams
	case "regtest":
		ChainParams = &chaincfg.RegressionNetParams
	default:
		logging.ErrorLogger.Fatalf("Error reading config file, invalid chain: %s", chain)
	}
}
