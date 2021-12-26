package polling

import (
	concensusModel "github.com/kaspanet/kaspad/domain/consensus/model"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/mstime"
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kasboard/processing/database"
	"github.com/stasatdaglabs/kasboard/processing/database/model"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
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
	err = pollTipAmountAndVirtualParentAmountAndPruningPointMovement(database, client)
	if err != nil {
		return err
	}
	err = pollMempoolSize(database, client)
	if err != nil {
		return err
	}
	return pollBlueHashrate(database, client)
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

func pollTipAmountAndVirtualParentAmountAndPruningPointMovement(database *database.Database, client *rpcclient.RPCClient) error {
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
	err = database.InsertVirtualParentAmount(virtualParentAmount)
	if err != nil {
		return err
	}

	mostRecentPruningPointMovement, err := database.MostRecentPruningPointMovement()
	if err != nil {
		return err
	}
	if mostRecentPruningPointMovement == nil ||
		mostRecentPruningPointMovement.PruningPointBlockHash != getBlockDAGInfoResponse.PruningPointHash {
		pruningPointMovement := &model.PruningPointMovement{
			Timestamp:             timestamp,
			PruningPointBlockHash: getBlockDAGInfoResponse.PruningPointHash,
		}
		err := database.InsertPruningPointMovement(pruningPointMovement)
		if err != nil {
			return err
		}
	}

	return nil
}

func pollMempoolSize(database *database.Database, client *rpcclient.RPCClient) error {
	getInfoResponse, err := client.GetInfo()
	if err != nil {
		return err
	}
	timestamp := mstime.Now().UnixMilliseconds()
	mempoolSize := &model.MempoolSize{
		Timestamp: timestamp,
		Size:      getInfoResponse.MempoolSize,
	}
	return database.InsertMempoolSize(mempoolSize)
}

func pollBlueHashrate(database *database.Database, client *rpcclient.RPCClient) error {
	virtualBlockHash := concensusModel.VirtualBlockHash.String()
	estimateNetworkHashesPerSecondResponse, err := client.EstimateNetworkHashesPerSecond(virtualBlockHash, 1000)
	if err != nil {
		return err
	}
	timestamp := mstime.Now().UnixMilliseconds()
	estimatedBlueHashrate := &model.EstimatedBlueHashrate{
		Timestamp:    timestamp,
		BlueHashrate: estimateNetworkHashesPerSecondResponse.NetworkHashesPerSecond,
	}
	return database.InsertEstimatedBlueHashrate(estimatedBlueHashrate)
}
