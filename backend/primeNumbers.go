package main

func calculatePrimeNumbers(limit int) (primeNumbers []int) {
	for i := 2; len(primeNumbers) < limit; i++ {
		var isPrime = true
		for _, prime := range primeNumbers {
			if i%prime == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primeNumbers = append(primeNumbers, i)
		}
	}
	return
}

func createRiddle(primeNumbers []int, steps int) (riddle int) {
	riddle = 1
	for range steps {
		var index = createRandomInt(len(primeNumbers))
		var primeNumber = primeNumbers[index]
		riddle = multiplyLimited(riddle, primeNumber, 1000_000)
	}
	return int(riddle)
}
