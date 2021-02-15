package main

import (
	"github.com/kaspanet/kaspad/util/panics"
	configPackage "github.com/stasatdaglabs/kashboard/processing/infrastructure/config"
	interruptPackage "github.com/stasatdaglabs/kashboard/processing/infrastructure/interrupt"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	"github.com/stasatdaglabs/kashboard/processing/kaspad_sync"
	"os"
)

var (
	log   = logging.Logger()
	spawn = panics.GoroutineWrapperFunc(log)
)

func main() {
	interrupt := interruptPackage.InterruptListener()
	defer log.Info("Shutdown complete")

	config, err := configPackage.Parse()
	if err != nil {
		logErrorAndExit("Could not parse command line arguments: %s", err)
	}

	spawn("kaspad_sync", func() {
		err := kaspad_sync.Start(config.RPCServerAddress)
		if err != nil {
			logErrorAndExit("Received error from Kaspad sync: %s", err)
		}
	})

	<-interrupt
}

func logErrorAndExit(errorLog string, logParameters ...interface{}) {
	log.Errorf(errorLog, logParameters...)
	os.Exit(1)
}