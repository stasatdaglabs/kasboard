package logging

import (
	"fmt"
	"github.com/kaspanet/kaspad/infrastructure/logger"
	"github.com/kaspanet/kaspad/util"
	"os"
	"path/filepath"
)

const (
	appDataDirectory = "kashboard"
	logFileName      = "kashboard.log"
	errorLogFileName = "kashboard_errors.log"
)

var (
	backendLog = logger.NewBackend()
	log        = backendLog.Logger("KSBD")
)

func init() {
	homeDir := util.AppDataDir(appDataDirectory, false)
	logFile := filepath.Join(homeDir, logFileName)
	errorLogFile := filepath.Join(homeDir, errorLogFileName)

	err := backendLog.AddLogFile(logFile, logger.LevelTrace)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error adding log file %s as log rotator for level %s: %s",
			logFileName, logger.LevelTrace, err)
		os.Exit(1)
	}
	err = backendLog.AddLogFile(errorLogFile, logger.LevelWarn)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error adding log file %s as log rotator for level %s: %s",
			errorLogFileName, logger.LevelWarn, err)
		os.Exit(1)
	}
}

func Logger() *logger.Logger {
	return log
}
