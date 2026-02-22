package hibit

// NewBitset returns a bitset capable of holding exactly size bits.
//
// The bitset is represented as a slice of uint64 words, where each word
// stores 64 consecutive bits.
//
// Active bits can be specified via activeBits. Each value represents a
// zero-based bit index to set. If any index is negative or greater than
// or equal to size, NewBitset returns nil.
func NewBitset(size int, activeBits ...int) []uint64 {
	if size <= 0 {
		return nil
	}
	bs := make([]uint64, (size+63)>>6)
	for _, bitIdx := range activeBits {
		if bitIdx < 0 || bitIdx >= size {
			return nil
		}
		wordPos := bitIdx >> 6
		bitPos := bitIdx & 63
		bs[wordPos] |= 1 << bitPos
	}

	return bs
}
