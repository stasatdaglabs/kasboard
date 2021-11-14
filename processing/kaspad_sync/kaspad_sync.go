package kaspad_sync

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kasboard/processing/database"
	"github.com/stasatdaglabs/kasboard/processing/database/model"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/config"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
	hashratePackage "github.com/stasatdaglabs/kasboard/processing/kaspad_sync/hashrate"
)

var log = logging.Logger()
var spawn = panics.GoroutineWrapperFunc(log)

func Start(config *config.Config, database *database.Database, client *rpcclient.RPCClient) (chan *model.Block, error) {
	blockAddedNotifications := make(chan *appmessage.BlockAddedNotificationMessage, config.ActiveNetParams.PruningDepth())
	err := client.RegisterForBlockAddedNotifications(func(notification *appmessage.BlockAddedNotificationMessage) {
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

	hashrate, err := hashratePackage.Hashrate(notification.Block.Header.Bits, config.ActiveNetParams.TargetTimePerBlock)
	if err != nil {
		return err
	}
	block := &model.Block{
		BlockHash:         notification.Block.VerboseData.Hash,
		BlueScore:         notification.Block.VerboseData.BlueScore,
		Timestamp:         notification.Block.Header.Timestamp,
		Hashrate:          hashrate,
		ParentAmount:      uint16(len(notification.Block.Header.Parents[0].ParentHashes)),
		TransactionAmount: uint16(len(notification.Block.Transactions)),
	}
	err = database.InsertBlock(block)
	if err != nil {
		return err
	}
	log.Infof("Added block %s with timestamp %d", block.BlockHash, block.Timestamp)

	blockChan <- block
	return nil
}
