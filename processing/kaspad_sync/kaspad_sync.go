package kaspad_sync

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/difficulty"
	"github.com/pkg/errors"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/database/model"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	"strconv"
	"time"
)

var log = logging.Logger()

func Start(rpcServerAddress string, database *database.Database) error {
	client, err := rpcclient.NewRPCClient(rpcServerAddress)
	if err != nil {
		return errors.Errorf("Could not connect to the Kaspad RPC server at %s: %s", rpcServerAddress, err)
	}

	return client.RegisterForBlockAddedNotifications(func(notification *appmessage.BlockAddedNotificationMessage) {
		handleBlockAddedNotifications(database, notification)
	})
}

func handleBlockAddedNotifications(database *database.Database, notification *appmessage.BlockAddedNotificationMessage) {
	bitsUint64, err := strconv.ParseUint(notification.BlockVerboseData.Bits, 16, 32)
	if err != nil {
		panic(err)
	}
	bitsUint32 := uint32(bitsUint64)
	bitsBigInt := difficulty.CompactToBig(bitsUint32)
	hashrateBigInt := hashrate(bitsBigInt, 1*time.Second)
	hashrateUint64 := hashrateBigInt.Uint64()

	block := &model.Block{
		BlockHash: notification.BlockVerboseData.Hash,
		BlueScore: notification.BlockVerboseData.BlueScore,
		Timestamp: notification.BlockVerboseData.Time,
		Bits:      bitsUint32,
		Hashrate:  hashrateUint64,
	}
	err = database.InsertBlock(block)
	if err != nil {
		panic(err)
	}
	log.Infof("Added block %s with blue score %d", block.BlockHash, block.BlueScore)
}
