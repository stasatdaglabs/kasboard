package polling

import (
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/mstime"
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/database/model"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	"time"
)

var log = logging.Logger()
var spawn = panics.GoroutineWrapperFunc(log)

func Start(database *database.Database, client *rpcclient.RPCClient) {
	spawn("polling", func() {
		for {
			err := poll(database, client)
			if err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}
	})
}

func poll(database *database.Database, client *rpcclient.RPCClient) error {
	err := pollHeaderAmountAndBlockAmount(database, client)
	if err != nil {
		return err
	}
	return pollTipAmountAndVirtualParentAmount(database, client)
}

func pollHeaderAmountAndBlockAmount(database *database.Database, client *rpcclient.RPCClient) error {
	getBlockCountResponse, err := client.GetBlockCount()
	if err != nil {
		return err
	}

	headerAmount := &model.HeaderAmount{
		Timestamp: mstime.Now().UnixMilliseconds(),
		Amount:    getBlockCountResponse.HeaderCount,
	}
	err = database.InsertHeaderAmount(headerAmount)
	if err != nil {
		return err
	}

	blockAmount := &model.BlockAmount{
		Timestamp: mstime.Now().UnixMilliseconds(),
		Amount:    getBlockCountResponse.BlockCount,
	}
	return database.InsertBlockAmount(blockAmount)
}

func pollTipAmountAndVirtualParentAmount(database *database.Database, client *rpcclient.RPCClient) error {
	getBlockDAGInfoResponse, err := client.GetBlockDAGInfo()
	if err != nil {
		return err
	}
	timestamp := mstime.Now().UnixMilliseconds()

	tipAmount := &model.TipAmount{
		Timestamp: timestamp,
		Amount:    uint32(len(getBlockDAGInfoResponse.TipHashes)),
	}
	err = database.InsertTipAmount(tipAmount)
	if err != nil {
		return err
	}

	virtualParentAmount := &model.VirtualParentAmount{
		Timestamp: timestamp,
		Amount:    uint16(len(getBlockDAGInfoResponse.VirtualParentHashes)),
	}
	return database.InsertVirtualParentAmount(virtualParentAmount)
}
