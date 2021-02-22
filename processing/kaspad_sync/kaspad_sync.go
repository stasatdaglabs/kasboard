package kaspad_sync

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/pkg/errors"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/database/model"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	hashratePackage "github.com/stasatdaglabs/kashboard/processing/kaspad_sync/hashrate"
)

var log = logging.Logger()
var spawn = panics.GoroutineWrapperFunc(log)

var blockAddedNotifications = make(chan *appmessage.BlockAddedNotificationMessage, 1_000_000)

func Start(rpcServerAddress string, database *database.Database) error {
	client, err := rpcclient.NewRPCClient(rpcServerAddress)
	if err != nil {
		return errors.Errorf("Could not connect to the Kaspad RPC server at %s: %s", rpcServerAddress, err)
	}

	err = client.RegisterForBlockAddedNotifications(func(notification *appmessage.BlockAddedNotificationMessage) {
		blockAddedNotifications <- notification
	})
	if err != nil {
		return err
	}

	spawn("handleBlockAddedNotifications", func() {
		for notification := range blockAddedNotifications {
			handleBlockAddedNotifications(database, notification)
		}
	})
	return nil
}

func handleBlockAddedNotifications(database *database.Database, notification *appmessage.BlockAddedNotificationMessage) {
	hashrate, err := hashratePackage.Hashrate(notification.BlockVerboseData.Bits)
	if err != nil {
		return
	}

	block := &model.Block{
		BlockHash: notification.BlockVerboseData.Hash,
		BlueScore: notification.BlockVerboseData.BlueScore,
		Timestamp: notification.BlockVerboseData.Time,
		Hashrate:  hashrate,
	}
	err = database.InsertBlock(block)
	if err != nil {
		panic(err)
	}
	log.Infof("Added block %s with blue score %d", block.BlockHash, block.BlueScore)
}
