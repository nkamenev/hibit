package hibit

import (
	"github.com/nkamenev/hibit/utils"
)

// Logic provides algorithms operating on BitTree structures.
//
// Logic is reusable and keeps internal scratch state (a stack) to avoid
// allocations during operations. It is not safe for concurrent use.
type Logic struct {
	stack *utils.Stack
}

// NewLogic creates a new Logic instance.
//
// Optional functional options may be provided to configure internal
// parameters (for example, stack size).
func NewLogic(opts ...func(*Logic)) *Logic {
	l := &Logic{}
	for _, o := range opts {
		o(l)
	}
	return l
}

// IntersectBitTrees computes the intersection of one or more BitTrees.
//
// It finds all bit indices that are set in **all input trees** and writes
// them into the provided output slice `out`.
//
// The function:
//   - assumes all trees have the same logical length
//   - performs a depth-first traversal of the bit trees
//   - skips entire subtrees when the bitwise AND of all trees is zero
//   - does not allocate
//
// Returned value is the number of indices written to `out`.
// If `out` is too small, writing stops early.
//
// Bit indices written to `out` are global indices
// (word index * 64 + bit index).
//
// IntersectBitTrees is not safe for concurrent calls on the same Logic
// instance.
func (l *Logic) IntersectBitTrees(out []int, trees ...*BitTree) int {
	if len(trees) == 0 {
		return 0
	}
	treeSz := trees[0].Len()
	for _, t := range trees {
		if t.Len() != treeSz {
			return 0
		}
	}

	if l.stack == nil || l.stack.Cap() < 64 {
		l.stack = utils.NewStack(64)
	} else {
		l.stack.Reset()
	}

	leafStart := treeSz >> 1

	l.stack.Push(1)

	write := 0

	for !l.stack.IsEmpty() {
		i := l.stack.Pop()
		and := trees[0].tree[i]
		for j := 1; j < len(trees); j++ {
			and &= trees[j].tree[i]
			if and == 0 {
				goto nextNode
			}
		}

		if i >= leafStart {
			if write >= len(out) {
				break
			}

			offset := (i - leafStart) << 6
			write += utils.PopBitIndexes(and, offset, out[write:])
		} else {
			l.stack.Push(i<<1 | 1)
			l.stack.Push(i << 1)
		}
	nextNode:
	}

	return write
}

// WithStackSize configures the initial stack capacity for Logic.
//
// A larger stack may reduce the need for reallocation when working
// with deep or wide trees.
func WithStackSize(size int) func(*Logic) {
	return func(l *Logic) {
		l.stack = utils.NewStack(size)
	}
}
