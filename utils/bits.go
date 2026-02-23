package utils

import "math/bits"

// PopBitIndexes extracts indexes of all set bits in x.
//
// For each set bit in x, its index (0..63) is added to the provided offset
// and written into out. The function returns the number of written indexes.
//
// Bits are processed from least significant to most significant.
// The caller must ensure that out has enough capacity to hold the result.
//
// This function performs no allocations and runs in O(k),
// where k is the number of set bits in x.
func PopBitIndexes(x uint64, offset int, out []int) int {
	i := 0
	for x != 0 {
		bit := bits.TrailingZeros64(x)
		out[i] = offset + bit
		i++
		x &= x - 1
	}
	return i
}

// NextPow2 returns the smallest power of two greater than or equal to n.
func NextPow2(n int) int {
	if n <= 1 {
		return 1
	}
	return 1 << bits.Len(uint(n-1))
}

// IsPow2 returns true if n is a power of two.
func IsPow2(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}
