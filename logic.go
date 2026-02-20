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

// IntersectBitTrees computes the intersection of two BitTrees.
//
// It finds all bit indices that are set in both trees and writes them
// into the provided output slice `out`.
//
// The function:
//   - assumes both trees have the same logical length
//   - performs a depth-first traversal of the bit tree
//   - skips entire subtrees when their AND is zero
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
func (l *Logic) IntersectBitTrees(left, right *BitTree, out []int) int {
	if left.Len() != right.Len() {
		return 0
	}

	if l.stack == nil || l.stack.Cap() < 64 {
		l.stack = utils.NewStack(64)
	} else {
		l.stack.Reset()
	}

	leafStart := len(left.tree) >> 1

	l.stack.Push(1)

	write := 0

	for !l.stack.IsEmpty() {
		i := l.stack.Pop()

		and := left.tree[i] & right.tree[i]
		if and == 0 {
			continue
		}

		if i >= leafStart {
			if write >= len(out) {
				break
			}

			leafIdx := i - leafStart
			offset := leafIdx << 6

			n := utils.PopBitIndexes(and, offset, out[write:])
			write += n
			continue
		}

		l.stack.Push(i<<1 | 1)
		l.stack.Push(i << 1)
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
