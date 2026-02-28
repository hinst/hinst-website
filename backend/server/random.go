package server

import (
	"crypto/rand"
	"math/big"
)

func createRandomInt(limit int) int {
	var limitBigInt = big.NewInt(int64(limit))
	var randomBigInt = AssertResultError(rand.Int(rand.Reader, limitBigInt))
	return int(randomBigInt.Int64())
}
