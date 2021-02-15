package kaspad_sync

import (
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/pkg/errors"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
)

var log = logging.Logger()

func Start(rpcServerAddress string) error {
	client, err := rpcclient.NewRPCClient(rpcServerAddress)
	if err != nil {
		return errors.Errorf("Could not connect to the Kaspad RPC server at %s: %s", rpcServerAddress, err)
	}

	response, err := client.GetVirtualSelectedParentBlueScore()
	if err != nil {
		return errors.Errorf("Could not get response: %s", err)
	}
	log.Infof("blueScore: %d", response.BlueScore)
	return nil
}
