package model

type Block struct {
	ID        uint64 `pg:",pk"`
	BlockHash string `pg:",use_zero"`
	BlueScore uint64 `pg:",use_zero"`
}
