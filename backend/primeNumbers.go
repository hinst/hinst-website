package main

const primeNumbersFileName = "primeNumbers.json"

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

var globalPrimeNumbers []int

func init() {
	if checkFileExists(primeNumbersFileName) {
		readJsonFile(primeNumbersFileName, &globalPrimeNumbers)
		assertCondition(
			len(globalPrimeNumbers) > 0,
			func() string { return "Need prime numbers" },
		)
	} else {
		panic("Need prime numbers")
	}
}

func createRiddle(steps int) (keys []int, riddle int) {
	riddle = 1
	for range steps {
		var index = createRandomInt(len(globalPrimeNumbers))
		var primeNumber = globalPrimeNumbers[index]
		keys = append(keys, primeNumber)
		riddle = multiplyLimited(riddle, primeNumber, 1000_000)
	}
	return
}
