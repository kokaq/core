package tests

import (
	"testing"

	"github.com/kokaq/core/utils"
)

func TestPower(t *testing.T) {
	tests := []struct {
		base, exp, want int
	}{
		{2, 0, 1},
		{2, 1, 2},
		{2, 3, 8},
		{5, 2, 25},
		{3, 4, 81},
		{10, 1, 10},
		{1, 100, 1},
		{0, 5, 0},
		{0, 0, 1}, // 0^0 is usually defined as 1 in programming
	}
	for _, tt := range tests {
		got := utils.Power(tt.base, tt.exp)
		if got != tt.want {
			t.Errorf("Power(%d, %d) = %d; want %d", tt.base, tt.exp, got, tt.want)
		}
	}
}

func TestLog2(t *testing.T) {
	tests := []struct {
		x, want int
	}{
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{1024, 10},
		{0, -1},
		{-5, -1},
	}
	for _, tt := range tests {
		got := utils.Log2(tt.x)
		if got != tt.want {
			t.Errorf("Log2(%d) = %d; want %d", tt.x, got, tt.want)
		}
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		x    int
		want bool
	}{
		{1, true},
		{2, true},
		{4, true},
		{8, true},
		{16, true},
		{3, false},
		{0, false},
		{-2, false},
		{1023, false},
		{1024, true},
	}
	for _, tt := range tests {
		got := utils.IsPowerOfTwo(tt.x)
		if got != tt.want {
			t.Errorf("IsPowerOfTwo(%d) = %v; want %v", tt.x, got, tt.want)
		}
	}
}
