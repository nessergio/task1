package main

import (
	"errors"
	"fmt"
	"golang.org/x/exp/mmap"
	"math/big"
)

// f is the Fibonacci cache table, it values found by fib() func
var f = []*big.Int{big.NewInt(0), big.NewInt(1)}

// maxFib stores n of maximum fibonacci number found
var maxFib uint64 = 2

// Using big math because maximum Fibonacci number that can be returned for int64 is
// F(93) = 12,200,160,415,121,876,738

// fib calculates Fibonacci number using big math.
func fib(n uint64) *big.Int {
	// Updating fibonacci table up to the needed number, if it is not already cached
	for ; maxFib <= n; maxFib++ {
		c := big.Int{}
		f = append(f, c.Add(f[maxFib-1], f[maxFib-2]))
	}
	return f[n]
}

// getPossibleCombinations calculates # of possible ways to decode the message
func getPossibleCombinations(p []byte) (*big.Int, error) {
	clusterSize := uint64(0)
	x := big.NewInt(1)
	a := p[0]
	// If string starts with 0 or not digit, return 0
	if a == 0x30 {
		return big.NewInt(0), errors.New("string starts with 0")
	} else if a < 0x31 || a > 0x39 {
		return big.NewInt(0), errors.New("string starts with non-digit character")
	}
	for i, b := range p[1:] {
		// If there is non-digit in the string, return 0
		if b < 0x30 || b > 0x39 {
			return big.NewInt(0), fmt.Errorf("encountered non-digit character at pos. %d", i)
		}
		// if next digit is 0, and value is not 10 or 20, return 0
		if b == 0x30 && a != 0x31 && a != 0x32 {
			return big.NewInt(0), fmt.Errorf("encountered 0 which can not be attached to %c at pos. %d", a, i)
		}
		// checking if value is in range 11-19, 21-26 (can be treated both like separate or single)
		if (a == 0x31 && b > 0x30) || (a == 0x32 && b > 0x30 && b <= 0x36) {
			// we are in a cluster
			clusterSize++
		} else if clusterSize > 0 {
			// exited from a cluster
			x.Mul(x, fib(clusterSize+2))
			clusterSize = 0
		}
		a = b
	}
	if clusterSize > 0 {
		// if the string ended on a cluster
		x.Mul(x, fib(clusterSize+2))
	}
	return x, nil
}

func main() {
	r, err := mmap.Open("test2.txt")
	if err != nil {
		panic(err)
	}
	p := make([]byte, r.Len())
	if _, err = r.ReadAt(p, 0); err != nil {
		panic(err)
	}
	x, err := getPossibleCombinations(p)
	if err != nil {
		panic(err)
	}
	fmt.Print(x)
}
