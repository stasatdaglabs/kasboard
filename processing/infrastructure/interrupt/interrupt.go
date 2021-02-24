package interrupt

import (
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
	"os"
	"os/signal"
)

var log = logging.Logger()

// InterruptListener listens for SIGINT (Ctrl+C).
// It returns a channel that is closed when either signal is received.
func InterruptListener() chan struct{} {
	interruptChan := make(chan struct{})
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, os.Interrupt)

		select {
		case <-interruptChannel:
			log.Infof("Shutting down...")
		}
		close(interruptChan)
	}()
	return interruptChan
}
