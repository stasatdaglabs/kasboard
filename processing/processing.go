package main

import (
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"
	"github.com/stasatdaglabs/kasboard/processing/analysis"
	databasePackage "github.com/stasatdaglabs/kasboard/processing/database"
	configPackage "github.com/stasatdaglabs/kasboard/processing/infrastructure/config"
	interruptPackage "github.com/stasatdaglabs/kasboard/processing/infrastructure/interrupt"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
	"github.com/stasatdaglabs/kasboard/processing/kaspad_sync"
	"github.com/stasatdaglabs/kasboard/processing/polling"
	"os"
)

var log = logging.Logger()

func main() {
	interrupt := interruptPackage.InterruptListener()
	defer log.Info("Shutdown complete")

	config, err := configPackage.Parse()
	if err != nil {
		logErrorAndExit("Could not parse command line arguments.\n%s", err)
	}

	database, err := databasePackage.Connect(config.DatabaseConnectionString)
	if err != nil {
		logErrorAndExit("Could not connect to database %s: %s", config.DatabaseConnectionString, err)
	}
	defer database.Close()

	client, err := rpcclient.NewRPCClient(config.RPCServerAddress)
	if err != nil {
		logErrorAndExit("Could not connect to the Kaspad RPC server at %s: %s",
			config.RPCServerAddress, err)
	}

	blockChan, err := kaspad_sync.Start(config, database, client)
	if err != nil {
		logErrorAndExit("Received error from Kaspad sync: %s", err)
	}
	analysis.Start(database, blockChan)
	polling.Start(database, client)

	<-interrupt
}

func logErrorAndExit(errorLog string, logParameters ...interface{}) {
	log.Errorf(errorLog, logParameters...)
	os.Exit(1)
}
