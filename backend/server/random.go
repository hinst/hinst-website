package server

import (
	"crypto/rand"
	"math/big"

	"github.com/hinst/go-common"
)

func createRandomInt(limit int) int {
	var limitBigInt = big.NewInt(int64(limit))
	var randomBigInt = common.AssertResultError(rand.Int(rand.Reader, limitBigInt))
	return int(randomBigInt.Int64())
}
