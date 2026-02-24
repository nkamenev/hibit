package hibit

import (
	"github.com/nkamenev/hibit/utils"
)

// Relation represents a binary relation between entities.
//
// Each Relation has two BitTrees:
//   - left  — bitset of entities in the left part of the relation
//   - right — bitset of entities in the right part of the relation
//
// Rows correspond to entities, columns correspond to the relation itself.
// Both BitTrees must have the same length (number of entities).
type Relation struct {
	left, right *BitTree
}

// NewRelation constructs a Relation from slices of entity indices.
// - size   — total number of entities
// - left   — indices of entities in the left part of the relation
// - right  — indices of entities in the right part of the relation
func NewRelation(size int, left, right []int) Relation {
	if size == 0 {
		return Relation{}
	}
	if size < utils.Max(len(left), len(right)) {
		return Relation{}
	}
	return Relation{
		left:  NewBitTree(NewBitset(size, left...)),
		right: NewBitTree(NewBitset(size, right...)),
	}
}

// NewRelationFromBitsets constructs a Relation from pre-built bitsets.
func NewRelationFromBitsets(size int, left, right []uint64) Relation {
	if size == 0 {
		return Relation{}
	}
	leafCount := utils.NextPow2(utils.Max(len(left), len(right)))
	return Relation{
		left:  NewBitTreeWithLeafCount(leafCount, left),
		right: NewBitTreeWithLeafCount(leafCount, right),
	}
}

// Size returns the number of entities (rows) in the relation.
func (r Relation) Size() int {
	return r.left.Len()
}
