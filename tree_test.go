package hibit

import (
	"fmt"
	"reflect"
	"testing"
)

var sizes = []int{
	1,
	8,
	64,
	256,
	1024,
	4096,
	8192,
	16384,
}

func TestNewBitTree(t *testing.T) {
	tests := map[string]struct {
		src  []uint64
		want []uint64
	}{
		"empty": {
			src:  nil,
			want: nil,
		},
		"single element": {
			src: []uint64{0b1010},
			// leafCount = 1
			// tree size = 2
			// index: 0 unused, 1=root/leaf
			want: []uint64{
				0,
				0b1010,
			},
		},
		"two elements": {
			src: []uint64{0b0101, 0b0010},
			// leafCount = 2
			// leaves at [2], [3]
			// root = OR(2,3)
			want: []uint64{
				0,
				0b0111,
				0b0101,
				0b0010,
			},
		},
		"four elements": {
			src: []uint64{0b0101, 0b0010, 0b0000, 0b1110},
			// leafCount = 4
			// internal nodes:
			// 2 = 4|5 = 0111
			// 3 = 6|7 = 1110
			// 1 = 2|3 = 1111
			want: []uint64{
				0,
				0b1111,
				0b0111,
				0b1110,
				0b0101,
				0b0010,
				0b0000,
				0b1110,
			},
		},
		"non power of two": {
			src: []uint64{1, 2, 4},
			// leafCount = 4
			// missing leaf is zero
			// leaves: [1,2,4,0]
			// 2 = 1|2 = 3
			// 3 = 4|0 = 4
			// 1 = 3|4 = 7
			want: []uint64{
				0,
				7,
				3,
				4,
				1,
				2,
				4,
				0,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bt := NewBitTree(tt.src)

			if tt.want == nil {
				if bt != nil {
					t.Fatalf("expected nil tree, got %v", bt.tree)
				}
				return
			}

			if !reflect.DeepEqual(bt.tree, tt.want) {
				t.Fatalf(
					"tree mismatch\nsrc:  %064b\nwant: %v\ngot:  %v",
					tt.src,
					tt.want,
					bt.tree,
				)
			}
		})
	}
}

func TestIntersectBitTrees(t *testing.T) {
	tests := map[string]struct {
		first  []uint64
		second []uint64
		third  []uint64
		want   []int
	}{
		"empty intersection": {
			first:  []uint64{0b0000},
			second: []uint64{0b0000},
			third:  []uint64{0b0000},
			want:   []int{},
		},
		"single bit": {
			first:  []uint64{0b0100},
			second: []uint64{0b0100},
			third:  []uint64{0b0100},
			want:   []int{2},
		},
		"no overlap": {
			first:  []uint64{0b1010},
			second: []uint64{0b0101},
			third:  []uint64{0b1010},
			want:   []int{},
		},
		"partial overlap": {
			first:  []uint64{0b1111},
			second: []uint64{0b0101},
			third:  []uint64{0b0100},
			want:   []int{2},
		},
		"multiple words": {
			first: []uint64{
				0b1111,
				0b0001,
			},
			second: []uint64{
				0b0101,
				0b0001,
			},
			third: []uint64{
				0b0100,
				0b0001,
			},
			want: []int{
				2,  // first word
				64, // second word
			},
		},
		"all bits": {
			first:  []uint64{^uint64(0)},
			second: []uint64{^uint64(0)},
			third:  []uint64{^uint64(0)},
			want: func() []int {
				out := make([]int, 64)
				for i := range 64 {
					out[i] = i
				}
				return out
			}(),
		},
	}

	logic := NewLogic(WithStackSize(64))

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			firstTree := NewBitTree(tt.first)
			secondTree := NewBitTree(tt.second)
			thirdTree := NewBitTree(tt.third)

			out := make([]int, 128)
			n := logic.IntersectBitTrees(out, firstTree, secondTree, thirdTree)
			got := out[:n]

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNewBitTree(b *testing.B) {
	for _, n := range sizes {
		b.Run(fmt.Sprintf("words=%d", n), func(b *testing.B) {
			src := make([]uint64, n)
			for i := range src {
				src[i] = uint64(i) * 0x9e3779b97f4a7c15
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = NewBitTree(src)
			}
		})
	}
}

func BenchmarkBuildLevels(b *testing.B) {
	const leafCount = 4096
	tree := make([]uint64, leafCount<<1)

	for i := range tree {
		tree[i] = uint64(i)
	}

	for b.Loop() {
		for j := leafCount - 1; j > 0; j-- {
			tree[j] = tree[j<<1] | tree[j<<1|1]
		}
	}
}

func genWords(n int, density float64) []uint64 {
	out := make([]uint64, n)

	if density <= 0 {
		return out
	}
	if density >= 1 {
		for i := range out {
			out[i] = ^uint64(0)
		}
		return out
	}

	bitsPerWord := int(64 * density)
	if bitsPerWord == 0 {
		bitsPerWord = 1
	}

	for i := range out {
		var v uint64
		for b := 0; b < bitsPerWord; b++ {
			v |= 1 << uint(b*64/bitsPerWord)
		}
		out[i] = v
	}
	return out
}

func BenchmarkIntersectBitTrees(b *testing.B) {
	var sizes = []int{
		// 1,
		// 8,
		// 64,
		256,
		1024,
		// 4096,
		// 8192,
		// 16384,
	}
	densities := []struct {
		name    string
		density float64
	}{
		// {"1bit", 1.0 / 64},
		{"10%", 0.10},
		// {"20%", 0.20},
		// {"30%", 0.30},
		// {"40%", 0.40},
		{"50%", 0.50},
		// {"60%", 0.60},
		// {"70%", 0.70},
		// {"80%", 0.80},
		// {"90%", 0.90},
		{"100%", 1.0},
	}

	logic := NewLogic(WithStackSize(64))

	treesCount := []int{256, 1024}

	for _, n := range sizes {
		for _, d := range densities {
			for _, tc := range treesCount {
				b.Run(
					fmt.Sprintf("words=%d/density=%s/trees=%d", n, d.name, tc),
					func(b *testing.B) {
						trees := make([]*BitTree, tc)
						for i := 0; i < tc; i++ {
							words := genWords(n, d.density)
							trees[i] = NewBitTree(words)
						}

						out := make([]int, n*64)

						b.ResetTimer()
						for i := 0; i < b.N; i++ {
							_ = logic.IntersectBitTrees(out, trees...)
						}
					},
				)
			}
		}
	}
}
