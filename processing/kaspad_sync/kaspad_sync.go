package kaspad_sync

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/pkg/errors"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/database/model"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/config"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	hashratePackage "github.com/stasatdaglabs/kashboard/processing/kaspad_sync/hashrate"
)

var log = logging.Logger()
var spawn = panics.GoroutineWrapperFunc(log)

func Start(config *config.Config, database *database.Database) (chan *model.Block, error) {
	client, err := rpcclient.NewRPCClient(config.RPCServerAddress)
	if err != nil {
		return nil, errors.Errorf("Could not connect to the Kaspad RPC server at %s: %s",
			config.RPCServerAddress, err)
	}

	blockAddedNotifications := make(chan *appmessage.BlockAddedNotificationMessage, config.ActiveNetParams.PruningDepth())
	err = client.RegisterForBlockAddedNotifications(func(notification *appmessage.BlockAddedNotificationMessage) {
		blockAddedNotifications <- notification
	})
	if err != nil {
		return nil, err
	}

	blockChan := make(chan *model.Block, config.ActiveNetParams.PruningDepth())
	spawn("handleBlockAddedNotifications", func() {
		for notification := range blockAddedNotifications {
			err := handleBlockAddedNotifications(config, database, notification, blockChan)
			if err != nil {
				panic(err)
			}
		}
	})
	return blockChan, nil
}

func handleBlockAddedNotifications(config *config.Config, database *database.Database,
	notification *appmessage.BlockAddedNotificationMessage, blockChan chan *model.Block) error {

	hashrate, err := hashratePackage.Hashrate(notification.BlockVerboseData.Bits, config.ActiveNetParams.TargetTimePerBlock)
	if err != nil {
		return err
	}

	block := &model.Block{
		BlockHash:         notification.BlockVerboseData.Hash,
		BlueScore:         notification.BlockVerboseData.BlueScore,
		Timestamp:         notification.BlockVerboseData.Time,
		Hashrate:          hashrate,
		ParentAmount:      uint16(len(notification.BlockVerboseData.ParentHashes)),
		TransactionAmount: uint16(len(notification.Block.Transactions)),
	}
	err = database.InsertBlock(block)
	if err != nil {
		return err
	}
	log.Infof("Added block %s with blue score %d", block.BlockHash, block.BlueScore)

	blockChan <- block
	return nil
}
