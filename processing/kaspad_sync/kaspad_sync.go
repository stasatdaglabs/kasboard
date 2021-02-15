package kaspad_sync

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/pkg/errors"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
)

var log = logging.Logger()

func Start(rpcServerAddress string, database *database.Database) error {
	client, err := rpcclient.NewRPCClient(rpcServerAddress)
	if err != nil {
		return errors.Errorf("Could not connect to the Kaspad RPC server at %s: %s", rpcServerAddress, err)
	}

	return client.RegisterForBlockAddedNotifications(handleBlockAddedNotifications)
}

func handleBlockAddedNotifications(notification *appmessage.BlockAddedNotificationMessage) {
	log.Infof("Received %d", notification.BlockVerboseData.BlueScore)
}
