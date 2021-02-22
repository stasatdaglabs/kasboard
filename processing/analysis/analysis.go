package analysis

import (
	"github.com/kaspanet/kaspad/util/panics"
	"github.com/stasatdaglabs/kashboard/processing/database"
	"github.com/stasatdaglabs/kashboard/processing/database/model"
	"github.com/stasatdaglabs/kashboard/processing/infrastructure/logging"
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
	const durationForAverage = 5 * time.Minute
	averageParentAmount, err := database.AverageParentAmount(block, durationForAverage)
	if err != nil {
		return err
	}

	analyzedBlock := &model.AnalyzedBlock{
		ID:                  block.ID,
		Timestamp:           block.Timestamp,
		AverageParentAmount: averageParentAmount,
	}
	return database.InsertAnalyzedBlock(analyzedBlock)
}
