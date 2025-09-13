package server

const primeNumbersFileName = "primeNumbers.json"
const primeNumbersLimit = 1000

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
		if len(globalPrimeNumbers) > primeNumbersLimit {
			globalPrimeNumbers = globalPrimeNumbers[0:primeNumbersLimit]
		}
		assertCondition(
			len(globalPrimeNumbers) > 0,
			func() string { return "Need prime numbers" },
		)
	} else {
		panic("Need prime numbers, but the file is missing: " + primeNumbersFileName)
	}
}
