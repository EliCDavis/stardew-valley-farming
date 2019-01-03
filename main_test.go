package main

import (
	"testing"
)

func TestValidBoards(t *testing.T) {

	var flagtests = []struct {
		board  []byte
		width  int
		height int
		out    bool
	}{
		{[]byte{'.'}, 1, 1, true},
		{[]byte{'x'}, 1, 1, true},
		{[]byte{'c'}, 1, 1, true},
		{[]byte{'c', 'c', 'c', 'c'}, 2, 2, true},
		{[]byte{'c', 'x', 'c', 'c'}, 2, 2, true},
		{[]byte{'c', 'c', 'x', '.'}, 2, 2, true},
		{[]byte{'c', 'x', 'c', 'x'}, 2, 2, true},
		{[]byte{
			'c', 'c', 'c',
			'c', 'c', 'c',
			'c', 'c', 'c',
		}, 3, 3, false},
		{[]byte{
			'c', 'c', 'c',
			'c', '.', 'c',
			'c', 'c', 'c',
		}, 3, 3, true},
		{[]byte{
			'c', 'c', 'c',
			'c', 'x', 'c',
			'c', 'c', 'c',
		}, 3, 3, true},
	}

	for _, tt := range flagtests {
		s := validLayout(tt.width, tt.height, tt.board)
		if s != tt.out {
			t.Errorf("%s got %t, want %t", tt.board, s, tt.out)
		}
	}

}

func BenchmarkValidBoards(b *testing.B) {

	var flagtests = []struct {
		board  []byte
		width  int
		height int
		out    bool
	}{
		{[]byte{'.'}, 1, 1, true},
		{[]byte{'x'}, 1, 1, true},
		{[]byte{'c'}, 1, 1, true},
		{[]byte{'c', 'c', 'c', 'c'}, 2, 2, true},
		{[]byte{'c', 'x', 'c', 'c'}, 2, 2, true},
		{[]byte{'c', 'c', 'x', '.'}, 2, 2, true},
		{[]byte{'c', 'x', 'c', 'x'}, 2, 2, true},
		{[]byte{
			'c', 'c', 'c',
			'c', 'c', 'c',
			'c', 'c', 'c',
		}, 3, 3, false},
		{[]byte{
			'c', 'c', 'c',
			'c', '.', 'c',
			'c', 'c', 'c',
		}, 3, 3, true},
		{[]byte{
			'c', 'c', 'c',
			'c', 'x', 'c',
			'c', 'c', 'c',
		}, 3, 3, true},
	}

	for n := 0; n < b.N; n++ {
		for _, tt := range flagtests {
			validLayout(tt.width, tt.height, tt.board)
		}
	}

}
