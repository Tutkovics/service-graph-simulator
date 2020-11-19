package main

import (
	"fmt"
	"sync"
)

// func main() {
// 	algo()
// }

// Sieve of Eratosthenes
func algoTest(number int, waitgroup *sync.WaitGroup) {
	max := number
	numbers := make([]bool, max+1)
	// Set values to ture
	for i := range numbers {
		numbers[i] = true
	}

	// main algorithm
	for p := 2; p*p <= max; p++ {
		if numbers[p] {
			for i := p * p; i <= max; i += p {

				numbers[i] = false
			}
		}
	}

	for p := 2; p <= max; p++ {
		if numbers[p] {
			fmt.Printf("%d ", p)
		}
	}

	waitgroup.Done()

	// fmt.Println(sum)
}
