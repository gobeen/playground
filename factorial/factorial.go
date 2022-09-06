package main

import (
	"fmt"
	"math/big"
)

func main() {
	var n int64
	fmt.Printf(">>> Factorial program using math/big package. Enter n: ")
	fmt.Scan(&n)
	bigFactorial := bigFactorial(n)
	fmt.Printf("%d! = %d\n", n, bigFactorial)
	fmt.Printf("number of bits: %d\n", bigFactorial.BitLen())

	fmt.Println(">>> Factorial program which returns int64.")
	fmt.Printf("%d! = %d\n", n, int64Factorial(n))
	fmt.Println("int64 overflow occures for n>20. 20! allocates 62 bits of memory.")

	fmt.Println(">>> Factorial via channel.")
	factorial := make(chan uint64)
	go factorialViaChannel(int(n), factorial)
	fmt.Printf("%d! = %d\n", n, <-factorial)
}

func bigFactorial(n int64) *big.Int {
	if n == 0 {
		return big.NewInt(1)
	}

	return big.NewInt(1).MulRange(1, n)
}

func int64Factorial(n int64) int64 {
	if n <= 0 {
		return 1
	}
	return n * int64Factorial(n-1)
}

func factorialViaChannel(n int, factorial chan uint64) {
	var computation uint64 = 1
	if n <= 0 {
		fmt.Printf("value cannot be less than zero: %d\n", n)
	} else {
		for i := 1; i <= n; i++ {
			computation *= uint64(i)
		}
	}
	factorial <- computation
}
