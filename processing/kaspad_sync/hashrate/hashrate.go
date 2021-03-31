package hashrate

import (
	"math/big"
	"time"

	"github.com/kaspanet/kaspad/util/difficulty"
)

var (
	// bigOne is 1 represented as a big.Int. It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// oneLsh256 is 1 shifted left 256 bits. It is defined here to avoid
	// the overhead of creating it multiple times.
	oneLsh256 = new(big.Int).Lsh(bigOne, 256)
)

func hashrate(target *big.Int, targetTimePerBlock time.Duration) *big.Int {
	// From: https://bitcoin.stackexchange.com/a/5557/40800
	// difficulty = hashrate / (2^256 / max_target / block_rate_in_seconds)
	// hashrate = difficulty * (2^256 / max_target / block_rate_in_seconds)
	// difficulty = max_target / target
	// hashrate = (max_target / target) * (2^256 / max_target / block_rate_in_seconds)
	// hashrate = 2^256 / (target * block_rate_in_seconds)

	tmp := new(big.Int)
	divisor := new(big.Int).Set(target)
	divisor.Mul(divisor, tmp.SetInt64(targetTimePerBlock.Milliseconds()))
	divisor.Div(divisor, tmp.SetInt64(int64(time.Second/time.Millisecond))) // Scale it up to seconds.
	divisor.Div(oneLsh256, divisor)
	return divisor
}

// Hashrate converts the given bits string to hashrate in uint64
func Hashrate(bits uint32, targetTimePerBlock time.Duration) (uint64, error) {
	bitsBigInt := difficulty.CompactToBig(bits)
	hashrateBigInt := hashrate(bitsBigInt, targetTimePerBlock)
	return hashrateBigInt.Uint64(), nil
}
