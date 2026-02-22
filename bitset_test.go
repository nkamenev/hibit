package hibit

import (
	"reflect"
	"testing"
)

func TestNewBitset(t *testing.T) {
	tests := map[string]struct {
		size int
		bits []int
		want []uint64
	}{
		"empty size": {
			size: 0,
			bits: nil,
			want: nil,
		},
		"single bit": {
			size: 1,
			bits: []int{0},
			want: []uint64{
				0b1,
			},
		},
		"multiple bits same word": {
			size: 8,
			bits: []int{0, 2, 7},
			want: []uint64{
				0b10000101,
			},
		},
		"multiple bits different words": {
			size: 130,
			bits: []int{0, 64, 129},
			want: []uint64{
				0b1,  // bit 0
				0b1,  // bit 64
				0b10, // bit 129 (129 - 128 = 1)
			},
		},
		"out of range bits ignored": {
			size: 10,
			bits: []int{-1, 0, 9, 10, 100},
			want: nil,
		},
		"duplicate bits": {
			size: 5,
			bits: []int{1, 1, 1},
			want: []uint64{
				0b10,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			bs := NewBitset(tt.size, tt.bits...)
			if !reflect.DeepEqual(bs, tt.want) {
				t.Fatalf("got %064b\nwant %064b", bs, tt.want)
			}
		})
	}
}
