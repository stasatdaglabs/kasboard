package kaspad_sync

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/pkg/errors"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/database/model"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
)

var log = logging.Logger()

func Start(rpcServerAddress string, database *database.Database) error {
	client, err := rpcclient.NewRPCClient(rpcServerAddress)
	if err != nil {
		return errors.Errorf("Could not connect to the Kaspad RPC server at %s: %s", rpcServerAddress, err)
	}

	greatestBlueScore, err := database.GreatestBlueScore()
	if err != nil {
		return err
	}
	log.Infof("Greatest blue score: %d", greatestBlueScore)

	return client.RegisterForBlockAddedNotifications(func(notification *appmessage.BlockAddedNotificationMessage) {
		handleBlockAddedNotifications(database, notification)
	})
}

func handleBlockAddedNotifications(database *database.Database, notification *appmessage.BlockAddedNotificationMessage) {
	block := &model.Block{
		BlockHash: notification.BlockVerboseData.Hash,
		BlueScore: notification.BlockVerboseData.BlueScore,
	}
	err := database.InsertBlock(block)
	if err != nil {
		panic(err)
	}
	log.Infof("Added block %s with blue score %d", block.BlockHash, block.BlueScore)
}
