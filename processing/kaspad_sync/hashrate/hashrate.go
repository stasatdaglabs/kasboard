package hashrate

import (
	"math"
	"math/big"
	"time"

	difficultyPackage "github.com/kaspanet/kaspad/util/difficulty"
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
	bitsBigInt := difficultyPackage.CompactToBig(bits)
	hashrateBigInt := hashrate(bitsBigInt, targetTimePerBlock)
	return hashrateBigInt.Uint64(), nil
}

// GetDifficultyRatio returns the proof-of-work difficulty as a multiple of the
// minimum difficulty using the passed bits field from the header of a block.
func GetDifficultyRatio(bits uint32, powMax *big.Int) float64 {
	// The minimum difficulty is the max possible proof-of-work limit bits
	// converted back to a number. Note this is not the same as the proof of
	// work limit directly because the block difficulty is encoded in a block
	// with the compact form which loses precision.
	target := difficultyPackage.CompactToBig(bits)

	difficulty := new(big.Rat).SetFrac(powMax, target)
	difficultyInFloat, _ := difficulty.Float64()

	roundingPrecision := float64(100)
	difficultyInFloat = math.Round(difficultyInFloat*roundingPrecision) / roundingPrecision

	return difficultyInFloat
}
