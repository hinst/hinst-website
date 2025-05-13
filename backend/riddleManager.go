package main

import "time"

type riddles struct {
}

func (riddles) create(steps int) (result riddleItem) {
	result.createdAt = time.Now()
	result.Product = 1
	for range steps {
		var index = createRandomInt(len(globalPrimeNumbers))
		var primeNumber = globalPrimeNumbers[index]
		result.keys = append(result.keys, primeNumber)
		result.Product = multiplyLimited(result.Product, primeNumber, 1000_000)
	}
	return
}
