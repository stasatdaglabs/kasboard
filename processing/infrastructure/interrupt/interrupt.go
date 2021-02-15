package interrupt

import (
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
	"os"
	signalPackage "os/signal"
)

var log = logging.Logger()

// interruptSignals defines the signals to catch in order to do a proper shutdown
var interruptSignals = []os.Signal{os.Interrupt}

// InterruptListener listens for OS Signals such as SIGINT (Ctrl+C).
// It returns a channel that is closed when either signal is received.
func InterruptListener() chan struct{} {
	interruptChan := make(chan struct{})
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signalPackage.Notify(interruptChannel, interruptSignals...)

		select {
		case signal := <-interruptChannel:
			log.Infof("Received signal (%s). Shutting down...", signal)
		}
		close(interruptChan)
	}()
	return interruptChan
}
