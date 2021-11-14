package inv_collection

import (
	"fmt"
	"github.com/kaspanet/kaspad/app/appmessage"
	kaspadConfigPackage "github.com/kaspanet/kaspad/infrastructure/config"
	"github.com/kaspanet/kaspad/infrastructure/network/netadapter/standalone"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/mstime"
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kasboard/processing/database"
	"github.com/stasatdaglabs/kasboard/processing/database/model"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/config"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
	"net"
	"time"
)

var log = logging.Logger()
var spawn = panics.GoroutineWrapperFunc(log)

func Start(config *config.Config, database *database.Database, client *rpcclient.RPCClient) {
	spawn("inv-collection", func() {
		for {
			err := poll(config, database, client)
			if err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}
	})
}

func poll(config *config.Config, database *database.Database, client *rpcclient.RPCClient) error {
	countOfBlockInvs, countOfTransactionInvs, err := collectInvs(config, client)
	if err != nil {
		return err
	}
	timestamp := mstime.Now().UnixMilliseconds()
	blockInvCount := &model.BlockInvCount{
		Timestamp: timestamp,
		Count:     countOfBlockInvs,
	}
	err = database.InsertBlockInvCount(blockInvCount)
	if err != nil {
		return err
	}
	transactionInvCount := &model.TransactionInvCount{
		Timestamp: timestamp,
		Count:     countOfTransactionInvs,
	}
	return database.InsertTransactionInvCount(transactionInvCount)
}

func collectInvs(config *config.Config, client *rpcclient.RPCClient) (uint32, uint32, error) {
	minimalNetAdapter, err := buildMinimalNetAdapter(config)
	if err != nil {
		return 0, 0, err
	}
	connectedPeerInfo, err := client.GetConnectedPeerInfo()
	if err != nil {
		return 0, 0, err
	}


	const maxPeersToConnectTo = 8
	peerInfos := connectedPeerInfo.Infos
	if len(peerInfos) > maxPeersToConnectTo {
		peerInfos = peerInfos[:maxPeersToConnectTo]
	}
	peersToRoutes := make(map[string]*standalone.Routes)
	for _, peerInfo := range peerInfos {
		host, _, err := net.SplitHostPort(peerInfo.Address)
		if err != nil {
			return 0, 0, err
		}
		connectAddress := net.JoinHostPort(host, config.NetParams().DefaultPort)
		routes, err := minimalNetAdapter.Connect(connectAddress)
		if err != nil {
			continue
		}
		peersToRoutes[peerInfo.Address] = routes
	}

	isRunning := true
	blockInvCount := uint32(0)
	transactionInvCount := uint32(0)
	for peer, routes := range peersToRoutes {
		routesCopy := routes
		spawn(fmt.Sprintf("collectInvs-%s", peer), func() {
			for {
				if !isRunning {
					routesCopy.Disconnect()
					return
				}
				message, err := routesCopy.IncomingRoute.DequeueWithTimeout(time.Minute)
				if err != nil {
					routesCopy.Disconnect()
					return
				}
				switch message.(type) {
				case *appmessage.MsgInvRelayBlock:
					blockInvCount++
				case *appmessage.MsgInvTransaction:
					transactionInvCount++
				}
			}
		})
	}
	time.Sleep(time.Minute)
	isRunning = false

	return blockInvCount, transactionInvCount, nil
}

func buildMinimalNetAdapter(config *config.Config) (*standalone.MinimalNetAdapter, error) {
	kaspadConfig := kaspadConfigPackage.DefaultConfig()
	kaspadConfig.NetworkFlags = config.NetworkFlags
	return standalone.NewMinimalNetAdapter(kaspadConfig)
}
