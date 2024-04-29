package src

import (
	"github.com/setavenger/blindbitd/src/utils"
)

var (
	DirectoryPath = "~/.blindbit"

	PathIpcSocketDir = ""
	PathIpcSocket    = ""

	PathConfig = ""

	PathDbUTXOs    = ""
	PathDbSettings = ""
	PathDbWallet   = ""

	PathToSecretScan  = ""
	PathToSecretSpend = ""
)

const PathEndingSocketDirPath = "/run"
const PathEndingSocketPath = PathEndingSocketDirPath + "/blindbit.socket"

const dataPath = "/data"

const PathEndingConfig = "/config"

const PathEndingUTXOs = dataPath + "/utxos"
const PathEndingSettings = dataPath + "/settings"
const PathEndingWallet = dataPath + "/wallet"

const PathEndingSecretScan = dataPath + "/secretScan"
const PathEndingSecretSpend = dataPath + "/secretSpend"

func SetPaths(path string) {
	if path != "" {
		DirectoryPath = path
	}

	PathIpcSocketDir = DirectoryPath + PathEndingSocketDirPath
	PathIpcSocket = DirectoryPath + PathEndingSocketPath

	PathConfig = DirectoryPath + PathEndingConfig
	PathDbUTXOs = DirectoryPath + PathEndingUTXOs
	PathDbSettings = DirectoryPath + PathEndingSettings
	PathDbWallet = DirectoryPath + PathEndingWallet

	PathToSecretScan = DirectoryPath + PathEndingSecretScan
	PathToSecretSpend = DirectoryPath + PathEndingSecretSpend

	// create the directories
	utils.TryCreateDirectoryPanic(DirectoryPath)
	utils.TryCreateDirectoryPanic(PathIpcSocketDir)

	utils.TryCreateDirectoryPanic(DirectoryPath + dataPath)

	return
}
