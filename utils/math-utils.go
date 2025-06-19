package utils

import "math/bits"

func Power(base int, exp int) int {
	result := 1
	x := base
	for exp > 0 {
		if exp%2 == 1 {
			result *= x
		}
		x *= x
		exp /= 2
	}
	return result
}

func Log2(x int) int {
	if x <= 0 {
		return -1 // Logarithm is undefined for non-positive numbers
	}
	return bits.Len(uint(x)) - 1
}

func IsPowerOfTwo(x int) bool {
	return x > 0 && (x&(x-1)) == 0
}
