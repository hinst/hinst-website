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
			func() string { return "Need prime numbers, but the array is empty" },
		)
	} else {
		panic("Need prime numbers, but the file is missing: " + primeNumbersFileName)
	}
}
