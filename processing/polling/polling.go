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
	return pollHeaderAmount(database, client)
}

func pollHeaderAmount(database *database.Database, client *rpcclient.RPCClient) error {
	getBlockCountResponse, err := client.GetBlockCount()
	if err != nil {
		return err
	}

	headerAmount := &model.HeaderAmount{
		Timestamp: mstime.Now().UnixMilliseconds(),
		Amount:    getBlockCountResponse.HeaderCount,
	}
	return database.InsertHeaderAmount(headerAmount)
}
