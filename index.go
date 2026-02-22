package hibit

// BitIndex represents a collection of BitTrees for multiple binary relations.
//
// Each binary relation has a left and a right BitTree:
//   - left[i] — bitset of entities in the left part of relation i
//   - right[i] — bitset of entities in the right part of relation i
//
// BitIndex provides methods to perform intersections (joins) between these
// relations efficiently using a reusable Logic instance.
//
// The Logic stack is lazily initialized if not provided via WithLogic.
type BitIndex struct {
	left, right []*BitTree
	logic       *Logic
}

// NewBitIndexFromBitsets constructs a BitIndex from slices of bitsets representing binary relations.
//
// Each entry in left and right represents the entities of a single binary relation.
// left and right must have the same length. A BitTree is built for each bitset.
//
// Optional functional options can be provided, e.g. WithLogic to set
// a custom Logic instance.
//
// Returns nil if left and right have different lengths.
func NewBitIndexFromBitsets(left, right [][]uint64, opts ...func(*BitIndex)) *BitIndex {
	if len(left) != len(right) {
		return nil
	}
	leftTrees := make([]*BitTree, len(left))
	rightTrees := make([]*BitTree, len(left))
	for i, lb := range left {
		leftTrees[i] = NewBitTree(lb)
		rightTrees[i] = NewBitTree(right[i])

	}
	bi := &BitIndex{
		left:  leftTrees,
		right: rightTrees,
	}
	for _, o := range opts {
		o(bi)
	}
	if bi.logic == nil {
		bi.logic = NewDefaultLogic()
	}
	return bi
}

// Join returns the intersection of a chain of relations specified
// by relIndices.
//
// It performs a chain join: for each consecutive pair of relations
// in relIndices, it intersects the right BitTree of the first with
// the left BitTree of the next. The resulting set of global bit
// indices is returned.
//
// If len(relIndices) < 2 or any index is out of range, Join returns nil.
//
// The returned slice is newly allocated and contains exactly the
// bit indices of the intersection.
func (bi *BitIndex) Join(out []int, relIndices []int) []int {
	if len(relIndices) < 2 {
		return nil
	}

	trees := make([]*BitTree, 0, len(relIndices)*2)
	for i := 0; i < len(relIndices)-1; i++ {
		curr := relIndices[i]
		next := relIndices[i+1]
		if curr < 0 ||
			curr >= len(bi.left) ||
			next < 0 ||
			next >= len(bi.left) {
			return nil
		}
		trees = append(trees, bi.right[curr], bi.left[next])
	}

	n := bi.logic.IntersectBitTrees(out, trees...)
	return out[:n]
}

// WithLogic sets a custom Logic instance to be used by BitIndex.
//
// This allows reusing a preconfigured Logic or providing a Logic
// with a different stack size.
func WithLogic(logic *Logic) func(*BitIndex) {
	return func(bi *BitIndex) {
		bi.logic = logic
	}
}
