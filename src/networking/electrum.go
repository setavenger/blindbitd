package networking

import (
	"context"
	"github.com/setavenger/blindbitd/src"
	"github.com/setavenger/blindbitd/src/logging"
	"github.com/setavenger/go-electrum/electrum"
	"time"
)

func CreateElectrumClient() (*electrum.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := electrum.NewClientTCP(ctx, src.ElectrumServerAddress, src.ElectrumTorProxyHost)
	if err != nil {
		logging.ErrorLogger.Println(err)
		return nil, err
	}

	return client, err
}
