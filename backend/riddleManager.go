package main

type riddleManager struct {
	steps int
	limit int
}

func (me *riddleManager) create() (result int) {
	result = 1
	for range me.steps {
		var index = createRandomInt(len(globalPrimeNumbers))
		var primeNumber = globalPrimeNumbers[index]
		result = multiplyLimited(result, primeNumber, me.limit)
	}
	return
}
