package main

import (
	databasePackage "github.com/stasatdaglabs/kashboard/processing/database"
	configPackage "github.com/stasatdaglabs/kashboard/processing/infrastructure/config"
	interruptPackage "github.com/stasatdaglabs/kashboard/processing/infrastructure/interrupt"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	"github.com/stasatdaglabs/kashboard/processing/kaspad_sync"
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

	err = kaspad_sync.Start(config, database)
	if err != nil {
		logErrorAndExit("Received error from Kaspad sync: %s", err)
	}

	<-interrupt
}

func logErrorAndExit(errorLog string, logParameters ...interface{}) {
	log.Errorf(errorLog, logParameters...)
	os.Exit(1)
}
