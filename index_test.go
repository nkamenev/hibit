package hibit

import (
	"reflect"
	"testing"
)

func TestBitIndexJoin(t *testing.T) {
	tests := map[string]struct {
		left       [][]uint64
		right      [][]uint64
		relIndices []int
		want       []int
	}{
		"chain three": {
			left: [][]uint64{
				{0b1100},
				{0b1110},
				{0b1011},
			},
			right: [][]uint64{
				{0b1010},
				{0b1110},
				{0b1110},
			},
			relIndices: []int{0, 1, 2},
			want:       []int{1, 3},
		},
		"two relations": {
			left: [][]uint64{
				{0b1100},
				{0b1110},
			},
			right: [][]uint64{
				{0b1010},
				{0b1110},
			},
			relIndices: []int{0, 1},
			want:       []int{1, 3},
		},
		"single relation": {
			left: [][]uint64{
				{0b1100},
			},
			right: [][]uint64{
				{0b1010},
			},
			relIndices: []int{0},
			want:       nil,
		},
		"empty indices": {
			left:       nil,
			right:      nil,
			relIndices: []int{},
			want:       nil,
		},
		"no overlap": {
			left: [][]uint64{
				{0b0010},
				{0b0100},
			},
			right: [][]uint64{
				{0b1000},
				{0b0001},
			},
			relIndices: []int{0, 1},
			want:       []int{},
		},
		"single bit overlap": {
			left: [][]uint64{
				{0b0101},
				{0b0011},
			},
			right: [][]uint64{
				{0b0101},
				{0b0010},
			},
			relIndices: []int{0, 1},
			want:       []int{0},
		},
		"non adjacent": {
			left: [][]uint64{
				{0b1111}, // rel 0
				{0b0000}, // rel 1
				{0b0010}, // rel 2
			},
			right: [][]uint64{
				{0b1111}, // rel 0
				{0b0000},
				{0b0010}, // rel 2
			},
			relIndices: []int{0, 2},
			want:       []int{1}, // 1111 & 0010 = 0010 -> bit 1
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bi := NewBitIndexFromBitsets(tt.left, tt.right)
			var out []int
			if bi != nil && len(bi.left) > 0 {
				out = make([]int, bi.left[0].Len()*64)
			}
			got := bi.Join(out, tt.relIndices)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("%s: Join(%v) = %v, want %v", name, tt.relIndices, got, tt.want)
			}
		})
	}
}
