package hibit

import "github.com/nkamenev/hibit/utils"

// BitIndex represents a collection of binary relations as BitTrees.
//
// Conceptually, it is a relation matrix:
//   - rows correspond to entities (each row = one entity)
//   - columns correspond to relations (each column = one binary relation)
//
// Each Relation has a left and right BitTree:
//   - left  — bitset of entities in the left part of the relation
//   - right — bitset of entities in the right part of the relation
//
// BitIndex provides efficient intersection (join) operations across relations
// using a reusable Logic instance. The Logic stack is lazily initialized
// if not provided via WithLogic.
type BitIndex struct {
	rels  []Relation
	logic *Logic
}

// NewBitIndex creates a BitIndex from a slice of Relations.
func NewBitIndex(relations []Relation, opts ...func(*BitIndex)) *BitIndex {
	bi := &BitIndex{}
	for _, o := range opts {
		o(bi)
	}

	if cap(bi.rels) < len(relations) {
		bi.rels = make([]Relation, 0, len(relations))
	}
	bi.rels = bi.rels[:0]
	bi.rels = append(bi.rels, relations...)

	if bi.logic == nil {
		bi.logic = NewDefaultLogic()
	}
	return bi
}

// NewBitIndexFromBitsets constructs a BitIndex from slices of bitsets.
//
// Each entry in left and right represents the indices of set bits
// for a single binary relation. A BitTree is built for each set.
//
// left and right must have the same length. Returns nil if not.
//
// Optional functional options can be provided (e.g., WithLogic).
func NewBitIndexFromBitsets(left, right [][]uint64, opts ...func(*BitIndex)) *BitIndex {
	if len(left) != len(right) {
		return nil
	}
	rels := make([]Relation, len(left))
	for i, lb := range left {
		rels[i] = NewRelationFromBitsets(utils.Max(len(lb), len(right[i])), lb, right[i])
	}
	bi := &BitIndex{
		rels: rels,
	}
	for _, o := range opts {
		o(bi)
	}
	if bi.logic == nil {
		bi.logic = NewDefaultLogic()
	}
	return bi
}

// AddRelation appends a new relation to the BitIndex.
func (bi *BitIndex) AddRelation(rel Relation) {
	bi.rels = append(bi.rels, rel)
}

// RowsNum returns the number of entities (rows) in the BitIndex.
func (bi *BitIndex) RowsNum() int {
	if len(bi.rels) == 0 {
		return 0
	}
	return bi.rels[0].Size()
}

// RelationsNum returns the number of relations (columns) in the BitIndex.
func (bi *BitIndex) RelationsNum() int {
	return len(bi.rels)
}

// Join performs a chain join across multiple relations specified by relIndices.
//
// For each consecutive pair of relations in relIndices:
//   - intersect the right BitTree of the first with the left BitTree of the next
//   - propagate the intersection
//
// The resulting set of entity indices is returned in `out[:n]`.
//
// Returns nil if relIndices has fewer than 2 elements, any index is out of range,
// or any left/right BitTree is nil.
func (bi *BitIndex) Join(out []int, relIndices []int) []int {
	if len(relIndices) < 2 || len(bi.rels) == 0 {
		return nil
	}

	trees := make([]*BitTree, 0, len(relIndices)*2)
	for i := 0; i < len(relIndices)-1; i++ {
		curr := relIndices[i]
		next := relIndices[i+1]
		if curr < 0 ||
			curr >= len(bi.rels) ||
			next < 0 ||
			next >= len(bi.rels) {
			return nil
		}
		if bi.rels[curr].right == nil || bi.rels[next].left == nil {
			return nil
		}
		trees = append(trees, bi.rels[curr].right, bi.rels[next].left)
	}

	n := bi.logic.IntersectBitTrees(out, trees...)
	return out[:n]
}

// WithLogic sets a custom Logic instance to be used by BitIndex.
func WithLogic(logic *Logic) func(*BitIndex) {
	return func(bi *BitIndex) {
		bi.logic = logic
	}
}

// WithRelationsCap preallocates capacity for relations in the BitIndex.
func WithRelationsCap(capacity int) func(*BitIndex) {
	return func(bi *BitIndex) {
		bi.rels = make([]Relation, 0, capacity)
	}
}
