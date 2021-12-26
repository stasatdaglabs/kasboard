package model

type Block struct {
	ID                uint64  `pg:",pk"`
	BlockHash         string  `pg:",use_zero"`
	BlueScore         uint64  `pg:",use_zero"`
	Timestamp         int64   `pg:",use_zero"`
	Hashrate          uint64  `pg:",use_zero"`
	ParentAmount      uint16  `pg:",use_zero"`
	TransactionAmount uint16  `pg:",use_zero"`
	Difficulty        float64 `pg:",use_zero"`
}

type AnalyzedBlock struct {
	ID                  uint64  `pg:",pk"`
	Timestamp           int64   `pg:",use_zero"`
	AverageParentAmount float64 `pg:",use_zero"`
	BlockRate           float64 `pg:",use_zero"`
	TransactionRate     float64 `pg:",use_zero"`
}

type HeaderAmount struct {
	Timestamp int64  `pg:",use_zero"`
	Amount    uint64 `pg:",use_zero"`
}

type BlockAmount struct {
	Timestamp int64  `pg:",use_zero"`
	Amount    uint64 `pg:",use_zero"`
}

type TipAmount struct {
	Timestamp int64  `pg:",use_zero"`
	Amount    uint32 `pg:",use_zero"`
}

type VirtualParentAmount struct {
	Timestamp int64  `pg:",use_zero"`
	Amount    uint16 `pg:",use_zero"`
}

type MempoolSize struct {
	Timestamp int64  `pg:",use_zero"`
	Size      uint64 `pg:",use_zero"`
}

type PruningPointMovement struct {
	Timestamp             int64  `pg:",use_zero"`
	PruningPointBlockHash string `pg:",use_zero"`
}

type BlockInvCount struct {
	Timestamp int64  `pg:",use_zero"`
	Count     uint32 `pg:",use_zero"`
}

type TransactionInvCount struct {
	Timestamp int64  `pg:",use_zero"`
	Count     uint32 `pg:",use_zero"`
}

type EstimatedBlueHashrate struct {
	Timestamp    int64  `pg:",use_zero"`
	BlueHashrate uint64 `pg:",use_zero"`
}
