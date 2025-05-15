package main

type riddles struct {
}

func (riddles) create(steps int) (result int) {
	result = 1
	for range steps {
		var index = createRandomInt(len(globalPrimeNumbers))
		var primeNumber = globalPrimeNumbers[index]
		result = multiplyLimited(result, primeNumber, 1000_000)
	}
	return
}
