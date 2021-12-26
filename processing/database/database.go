package database

import (
	"fmt"
	pg "github.com/go-pg/pg/v9"
	"github.com/kaspanet/kaspad/util/mstime"
	"github.com/pkg/errors"
	"github.com/stasatdaglabs/kasboard/processing/database/model"
	"github.com/stasatdaglabs/kasboard/processing/infrastructure/logging"
	"strings"
	"time"
)

var (
	log              = logging.Logger()
	allowedTimeZones = map[string]struct{}{
		"UTC":     {},
		"Etc/UTC": {},
	}
)

// Connect connects to the database mentioned in the config variable.
func Connect(connectionString string) (*Database, error) {
	migrator, driver, err := openMigrator(connectionString)
	if err != nil {
		return nil, err
	}
	isCurrent, version, err := isCurrent(migrator, driver)
	if err != nil {
		return nil, errors.Wrapf(err, "error checking whether the database is current")
	}
	if !isCurrent {
		log.Warnf("Database is not current (version %d). Migrating...", version)
		err := migrate(connectionString)
		if err != nil {
			return nil, errors.Wrapf(err, "could not migrate database")
		}
	}

	connectionOptions, err := pg.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	pgDB := pg.Connect(connectionOptions)

	err = validateTimeZone(pgDB)
	if err != nil {
		return nil, errors.Wrapf(err, "could not validate database timezone")
	}

	return &Database{
		database: pgDB,
	}, nil
}

func validateTimeZone(db *pg.DB) error {
	var timeZone string
	_, err := db.QueryOne(pg.Scan(&timeZone), `SELECT current_setting('TIMEZONE') as time_zone`)

	if err != nil {
		return errors.WithMessage(err, "some errors were encountered when "+
			"checking the database timezone:")
	}

	if _, ok := allowedTimeZones[timeZone]; !ok {
		return errors.Errorf("to prevent conversion errors - Kasparov should only run with "+
			"a database configured to use one of the allowed timezone. Currently configured timezone "+
			"is %s. Allowed time zones: %s", timeZone, allowedTimezonesString())
	}
	return nil
}

func allowedTimezonesString() string {
	keys := make([]string, 0, len(allowedTimeZones))
	for allowedTimeZone := range allowedTimeZones {
		keys = append(keys, fmt.Sprintf("'%s'", allowedTimeZone))
	}
	return strings.Join(keys, ", ")
}

type Database struct {
	database *pg.DB
}

func (db *Database) InsertBlock(block *model.Block) error {
	return db.database.Insert(block)
}

func (db *Database) AverageParentAmount(fromBlock *model.Block, duration time.Duration) (float64, error) {
	endTimestamp := mstime.UnixMilliseconds(fromBlock.Timestamp)
	startTimestamp := endTimestamp.Add(-duration)

	var result struct {
		Average float64
	}
	_, err := db.database.QueryOne(&result, "SELECT AVG(parent_amount) as average FROM blocks WHERE timestamp > ? AND timestamp < ?",
		startTimestamp.UnixMilliseconds(), endTimestamp.UnixMilliseconds())
	if err != nil {
		return 0, err
	}
	return result.Average, nil
}

func (db *Database) BlockCount(fromBlock *model.Block, duration time.Duration) (uint64, error) {
	endTimestamp := mstime.UnixMilliseconds(fromBlock.Timestamp)
	startTimestamp := endTimestamp.Add(-duration)

	var result struct {
		Count uint64
	}
	_, err := db.database.QueryOne(&result, "SELECT COUNT(*) as count FROM blocks WHERE timestamp > ? AND timestamp < ?",
		startTimestamp.UnixMilliseconds(), endTimestamp.UnixMilliseconds())
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func (db *Database) TransactionCountWithoutCoinbase(fromBlock *model.Block, duration time.Duration) (uint64, error) {
	endTimestamp := mstime.UnixMilliseconds(fromBlock.Timestamp)
	startTimestamp := endTimestamp.Add(-duration)

	var result struct {
		TransactionAmount uint64
		BlockAmount       uint64
	}
	_, err := db.database.QueryOne(&result, "SELECT SUM(transaction_amount) as transaction_amount, COUNT(transaction_amount) as block_amount FROM blocks WHERE timestamp > ? AND timestamp < ?",
		startTimestamp.UnixMilliseconds(), endTimestamp.UnixMilliseconds())
	if err != nil {
		return 0, err
	}
	return result.TransactionAmount - result.BlockAmount, nil
}

func (db *Database) InsertAnalyzedBlock(analyzedBlock *model.AnalyzedBlock) error {
	return db.database.Insert(analyzedBlock)
}

func (db *Database) InsertHeaderAmount(headerAmount *model.HeaderAmount) error {
	return db.database.Insert(headerAmount)
}

func (db *Database) InsertBlockAmount(blockAmount *model.BlockAmount) error {
	return db.database.Insert(blockAmount)
}

func (db *Database) InsertTipAmount(tipAmount *model.TipAmount) error {
	return db.database.Insert(tipAmount)
}

func (db *Database) InsertVirtualParentAmount(virtualParentAmount *model.VirtualParentAmount) error {
	return db.database.Insert(virtualParentAmount)
}

func (db *Database) InsertMempoolSize(mempoolSize *model.MempoolSize) error {
	return db.database.Insert(mempoolSize)
}

func (db *Database) InsertEstimatedBlueHashrate(estimatedBlueHashrate *model.EstimatedBlueHashrate) error {
	return db.database.Insert(estimatedBlueHashrate)
}

func (db *Database) MostRecentPruningPointMovement() (*model.PruningPointMovement, error) {
	var pruningPointMovements []*model.PruningPointMovement
	_, err := db.database.Query(&pruningPointMovements, "SELECT * FROM pruning_point_movements ORDER BY timestamp DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	if len(pruningPointMovements) == 0 {
		return nil, nil
	}
	return pruningPointMovements[0], nil
}

func (db *Database) InsertPruningPointMovement(pruningPointMovement *model.PruningPointMovement) error {
	return db.database.Insert(pruningPointMovement)
}

func (db *Database) InsertBlockInvCount(blockInvCount *model.BlockInvCount) error {
	return db.database.Insert(blockInvCount)
}

func (db *Database) InsertTransactionInvCount(transactionInvCount *model.TransactionInvCount) error {
	return db.database.Insert(transactionInvCount)
}

func (db *Database) Close() {
	_ = db.database.Close()
}
