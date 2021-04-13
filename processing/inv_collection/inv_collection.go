package inv_collection

import (
	"fmt"
	"github.com/kaspanet/kaspad/app/appmessage"
	kaspadConfigPackage "github.com/kaspanet/kaspad/infrastructure/config"
	"github.com/kaspanet/kaspad/infrastructure/network/netadapter/standalone"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kasboard/processing/database"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/config"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
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
	invCount, err := collectInvs(config, client)
	if err != nil {
		return err
	}
	log.Infof("invCount: %d", invCount)

	return nil
}

func collectInvs(config *config.Config, client *rpcclient.RPCClient) (uint64, error) {
	minimalNetAdapter, err := buildMinimalNetAdapter(config)
	if err != nil {
		return 0, err
	}
	connectedPeerInfo, err := client.GetConnectedPeerInfo()
	if err != nil {
		return 0, err
	}

	peersToRoutes := make(map[string]*standalone.Routes)
	for _, peerInfo := range connectedPeerInfo.Infos {
		routes, err := minimalNetAdapter.Connect(peerInfo.Address)
		if err != nil {
			return 0, err
		}
		peersToRoutes[peerInfo.Address] = routes
	}

	isRunning := true
	invCount := uint64(0)
	for peer, routes := range peersToRoutes {
		spawn(fmt.Sprintf("collectInvs-%s", peer), func() {
			for {
				if !isRunning {
					return
				}
				_, err := routes.WaitForMessageOfType(appmessage.CmdInvTransaction, time.Minute)
				if err != nil {
					return
				}
				invCount++
			}
		})
	}
	time.Sleep(time.Minute)
	isRunning = false

	return invCount, nil
}

func buildMinimalNetAdapter(config *config.Config) (*standalone.MinimalNetAdapter, error) {
	kaspadConfig := kaspadConfigPackage.DefaultConfig()
	kaspadConfig.NetworkFlags = config.NetworkFlags
	return standalone.NewMinimalNetAdapter(kaspadConfig)
}
