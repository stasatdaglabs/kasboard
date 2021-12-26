package analysis

import (
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kasboard/processing/database"
	"github.com/stasatdaglabs/kasboard/processing/database/model"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
	"time"
)

var log = logging.Logger()
var spawn = panics.GoroutineWrapperFunc(log)

func Start(database *database.Database, blockChan chan *model.Block) {
	spawn("analysis", func() {
		for block := range blockChan {
			err := handleBlock(database, block)
			if err != nil {
				panic(err)
			}
		}
	})
}

func handleBlock(database *database.Database, block *model.Block) error {
	const durationForAnalysis = 1 * time.Minute

	averageParentAmount, err := database.AverageParentAmount(block, durationForAnalysis)
	if err != nil {
		return err
	}

	blockCount, err := database.BlockCount(block, durationForAnalysis)
	if err != nil {
		return err
	}
	blockRate := float64(blockCount) / durationForAnalysis.Seconds()

	transactionCount, err := database.TransactionCountWithoutCoinbase(block, durationForAnalysis)
	if err != nil {
		return err
	}
	transactionRate := float64(transactionCount) / durationForAnalysis.Seconds()

	averagePropagationDelay, err := database.AveragePropagationDelay(block, durationForAnalysis)
	if err != nil {
		return err
	}

	analyzedBlock := &model.AnalyzedBlock{
		ID:                      block.ID,
		Timestamp:               block.Timestamp,
		AverageParentAmount:     averageParentAmount,
		BlockRate:               blockRate,
		TransactionRate:         transactionRate,
		AveragePropagationDelay: averagePropagationDelay,
	}
	err = database.InsertAnalyzedBlock(analyzedBlock)
	if err != nil {
		return err
	}

	log.Infof("Added analyzed block %s with timestamp %d", block.BlockHash, block.Timestamp)
	return nil
}
