package hibit

import "github.com/nkamenev/hibit/utils"

// BitTree is a compact binary tree representation of a bitset.
//
// Leaves store 64-bit words of the bitset, while internal nodes store
// the bitwise OR of their children. This allows fast pruning of entire
// subtrees during queries such as intersection.
//
// The tree is stored in a flat array using 1-based indexing semantics:
//   - index 1 is the root
//   - for node i, children are at 2*i and 2*i+1
//
// BitTree is immutable after construction.
type BitTree struct {
	tree []uint64
}

// NewBitTree builds a BitTree from a slice of 64-bit words.
//
// The input slice represents a bitset where each uint64 stores 64 bits.
// The number of leaves is rounded up to the next power of two; missing
// leaves are implicitly treated as zero.
func NewBitTree(src []uint64) *BitTree {
	if len(src) == 0 {
		return nil
	}

	leafCount := utils.NextPow2(len(src))
	tree := make([]uint64, leafCount<<1)

	// leaves
	copy(tree[leafCount:], src)

	// build bottom-up
	for i := leafCount - 1; i > 0; i-- {
		tree[i] = tree[i<<1] | tree[i<<1|1]
	}

	return &BitTree{tree: tree}
}

// NewBitTreeWithLeafCount builds a BitTree with a specified number of leaves.
func NewBitTreeWithLeafCount(leafCount int, src []uint64) *BitTree {
	if !utils.IsPow2(leafCount) ||
		len(src) == 0 ||
		leafCount < len(src) {
		return nil
	}

	tree := make([]uint64, leafCount<<1)

	// leaves
	copy(tree[leafCount:], src)

	// build bottom-up
	for i := leafCount - 1; i > 0; i-- {
		tree[i] = tree[i<<1] | tree[i<<1|1]
	}

	return &BitTree{tree: tree}
}

// Len returns the total number of nodes in the tree, including leaves and internal nodes.
func (bt *BitTree) Len() int {
	return len(bt.tree)
}
