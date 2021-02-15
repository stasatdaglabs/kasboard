package main

import (
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	configPackage "github.com/stasatdaglabs/kashboard/processing/config"
	"os"
)

func main() {
	config, err := configPackage.Parse()
	if err != nil {
		logErrorAndExit("Could not parse command line arguments: %s", err)
	}

	client, err := rpcclient.NewRPCClient(config.RPCServerAddress)
	if err != nil {
		logErrorAndExit("Could not connect to the Kaspad RPC server at %s: %s", config.RPCServerAddress, err)
	}

	response, err := client.GetVirtualSelectedParentBlueScore()
	if err != nil {
		logErrorAndExit("Could not get response: %s", err)
	}
	log.Infof("blueScore: %d", response.BlueScore)
}

func logErrorAndExit(errorLog string, logParameters ...interface{}) {
	log.Errorf(errorLog, logParameters...)
	os.Exit(1)
}
