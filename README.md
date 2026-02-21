# hibit

`hibit` is a library of hierarchical zero-allocation bitsets based on a segment tree representation.

It is designed for fast intersection of large bitsets with sparse or uneven
distribution of set bits, where skipping empty subtrees gives a significant
speedup over linear scans.

## Key ideas

- **BitTree segment tree**:  
  Leaves store 64-bit words. Internal nodes store bitwise OR of children.  
  Allows fast pruning of empty subtrees.

- **Flat array layout**:  
  1-based indexing. Node i → children at 2*i, 2*i+1. Leaves hold actual bits.

- **DFS intersection**:  
  Uses zero-allocation stack. Skips branches where AND == 0. Traverses only nodes with possible common bits.

- **High performance**:  
  No heap allocations. Worst-case O(N), typical case much faster due to pruning.

## Intersection

Time complexity:

- **Best case:** `O(log N)` — early pruning of empty branches
- **Average case:** `O(K * log N)` — depends on sparsity and distribution of set bits
- **Worst case:** `O(N)` — full scan of all words (all bits set)

Where `N` is the number of 64-bit words and `K` is the number of visited nodes.

## Installation

```bash
go get github.com/nkamenev/hibit
```

## Example

```go
first := []uint64{
    0b1111, // bits 0..3
    0b0001, // bit 64
}

second := []uint64{
    0b0101, // bits 0 and 2
    0b0001, // bit 64
}

third := []uint64{
    0b0100, // bit 2
    0b0001, // bit 64
}

firstTree := hibit.NewBitTree(first)
secondTree := hibit.NewBitTree(second)
thirdTree := hibit.NewBitTree(third)

logic := hibit.NewLogic()

out := make([]int, 24)
n := logic.IntersectBitTrees(out, firstTree, secondTree, thirdTree)

result := out[:n]
// result == []int{2, 64}
```

## Testing

To run tests, use:

```bash
go test -v
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
