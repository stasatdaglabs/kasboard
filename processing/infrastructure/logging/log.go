package logging

import (
	"github.com/kaspanet/kaspad/infrastructure/logger"
	"github.com/kaspanet/kaspad/util"
	"path/filepath"
)

var (
	homeDir      = util.AppDir("kasboard", false)
	logFile      = filepath.Join(homeDir, "kasboard.log")
	errorLogFile = filepath.Join(homeDir, "kasboard_errors.log")
	log          = logger.RegisterSubSystem("KSBD")
)

func InitLog() {
	logger.InitLog(logFile, errorLogFile)
	logger.SetLogLevels(logger.LevelInfo)
}

func Logger() *logger.Logger {
	return log
}

func Close() {
	log.Backend().Close()
}
